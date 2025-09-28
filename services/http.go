package services

import (
  "context"
  "encoding/json"
  "errors"
  "log"
  "fmt"
  "strings"
  "net/http"
  "os"
  "path"
  "time"
  // "database/sql"
  "gitlab.com/pedrokoblitz/opus-go/actions"
  "gitlab.com/pedrokoblitz/opus-go/quality"
  "github.com/julienschmidt/httprouter"
  "golang.org/x/time/rate"

  "gitlab.com/pedrokoblitz/opus-go/modules/auth"
)

type HTTPService struct {
  container   actions.Container
  hub      *PubSubHub

  server  *http.Server
  handler *RouteHandler
  router  *httprouter.Router
}

func NewHTTPService(
  container actions.Container,
  hub *PubSubHub,
) *HTTPService {

  router := httprouter.New()

  router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "Not found hurr durr", http.StatusNotFound)
    // h.handleHTMLError(w, http.StatusNotFound, err.Error())
    
  })

  router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    // h.handleHTMLError(w, http.StatusNotFound, err.Error())
  })

  cm := quality.NewConnectionManager(quality.ConnectionConfig{
    MaxConns:     1000,
    IdleTimeout:  time.Minute * 30,
    WriteTimeout: time.Second * 10,
  })

  rl := quality.NewRateLimiter(
    rate.Limit(100), // 100 req/sec
    50,              // Burst capacity
    time.Minute*10,  // Cleanup interval
  )

  handler := NewRouteHandler(
    hub,
    container,
    cm,
    rl,
  )
  cfg := container.Config()
  port := cfg.HTTP.Port
  service := &HTTPService{
    server: &http.Server{
      Addr:    fmt.Sprintf(":%d", port),
      Handler: router,
    },
    handler: handler,
    container:  container,
    hub:     hub,
    router:  router,
  }

  service.registerRoutes(router)
  service.registerTopics()

  return service
}

func (s *HTTPService) registerTopics() {
  s.hub.CreateTopic("before:action")
  s.hub.CreateTopic("error:action")
  s.hub.CreateTopic("success:action")
}

func (s *HTTPService) registerActionTopics(rc actions.HttpRoute) {
  beforeAction := "before:" + rc.Action
  s.hub.CreateTopic(beforeAction)
  errorAction := "error:" + rc.Action
  s.hub.CreateTopic(errorAction)
  successAction := "success:" + rc.Action
  s.hub.CreateTopic(successAction)
  log.Println(beforeAction, errorAction, successAction)
}

func (s *HTTPService) registerRoutes(router *httprouter.Router) {
  // Serve static files from multiple directories
  staticDirs := []string{
    "./resources/assets",
  }

  for _, dir := range staticDirs {
    if _, err := os.Stat(dir); err == nil {
      // Convert httprouter to handle static files
      fs := http.FileServer(http.Dir(dir))
      prefix := "/assets/"
      router.Handler(http.MethodGet, prefix+"*filepath", http.StripPrefix(prefix, fs))
      log.Printf("Serving static files from %s at %s", dir, prefix)
    }
  }

  // TODO: webhook route

  // Register modules
  cfg := s.container.Config()
  modules := cfg.Modules
  for _, module := range modules {
    modulePath := "./resources/modules/" + module + "/module.yml"
    moduleConfig, err := actions.LoadModuleConfig(modulePath)
    if err != nil {
      log.Fatal(err)
    }
    for _, sub := range moduleConfig.Module.Subscribers {

      opts := SubscriberOptions{
        Concurrency: 1,
        RetryPolicy: RetryPolicy{
          MaxAttempts: 2,
          Backoff: 10,
        },
      }
      fmt.Println("SUBSCRIBED!!!")
      s.hub.Subscribe(sub.Topic, sub.Action, opts)
    }

    var fullPath string
    for _, rc := range moduleConfig.Module.Http.Api.Routes {
      rc.Module = moduleConfig.Module.Name

      fullPath = path.Join("/", moduleConfig.Module.Http.Api.Prefix, rc.Pattern)
      // Convert HTTP method to httprouter format
      method := "GET"
      if rc.Method != "" {
        method = rc.Method
      }
      
      router.Handle(method, fullPath, s.handler.HandleHTTPRouter(rc))
      log.Println("api route:", fullPath, rc.Method)
      
      // pubsub topics
      s.registerActionTopics(rc)
    }

    for _, rc := range moduleConfig.Module.Http.Views.Routes {
      rc.Module = moduleConfig.Module.Name

      if rc.Document == "" {
        rc.Document = moduleConfig.Module.Http.Views.Document
      }
      
      fullPath = path.Join("/", moduleConfig.Module.Http.Views.Prefix, rc.Pattern)
      // Views typically use GET method
      router.Handle(http.MethodGet, fullPath, s.handler.HandleHTTPRouter(rc))
      log.Println("view route:", fullPath, rc.Method)
    }
  }
}

func (s *HTTPService) Start() error {
  return s.server.ListenAndServe()
}

func (s *HTTPService) Shutdown(ctx context.Context) error {
  // Stop rate limiter cleanup goroutine
  if s.handler.rl != nil {
    s.handler.rl.Stop()
  }
  
  // Close connection manager
  if s.handler.cm != nil {
    s.handler.cm.Close()
  }
  return s.server.Shutdown(ctx)
}

type RouteHandler struct {
  hub      *PubSubHub
  container actions.Container
  cm       *quality.ConnectionManager
  rl       *quality.RateLimiter
}

// Update NewRouteHandler
func NewRouteHandler(
  hub *PubSubHub,
  container actions.Container,
  cm *quality.ConnectionManager,
  rl *quality.RateLimiter,
) *RouteHandler {

  handler := &RouteHandler{
    hub:   hub,
    container: container,
    cm:       cm,
    rl:       rl,
  }
  
  return handler
}

// HandleHTTPRouter returns an httprouter.Handle function
func (h *RouteHandler) HandleHTTPRouter(rc actions.HttpRoute) httprouter.Handle {
  return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    // Use shared rate limiter
    if !h.rl.GetVisitor(r.RemoteAddr).Allow() {
      http.Error(w, "too many requests", http.StatusTooManyRequests)
      return
    }

    // Use shared connection manager
    h.cm.Acquire()
    defer h.cm.Release()

    err := h.handleRoute(w, r, rc, ps)
    if err != nil {
      h.handleError(w, r, err, rc)
    }
  }
}

func (h *RouteHandler) handleRoute(w http.ResponseWriter, r *http.Request, rc actions.HttpRoute, ps httprouter.Params) error {
  var err error

  ctx := r.Context()
  ctx = context.WithValue(ctx, "story", ps.ByName("story"))
  ctx = context.WithValue(ctx, "definition", ps.ByName("definition"))

  if rc.Protected == true {
    tokenHeader := r.Header.Get("Authorization")
    token := strings.ReplaceAll(tokenHeader, "Bearer: ", "")
    claims, err := auth.ParseBearerToken(token)
    if err != nil {
      return err
    }
    ctx = context.WithValue(ctx, "authToken", token)
    ctx = context.WithValue(ctx, "userID", claims.UserID)
  }
  r = r.WithContext(ctx)

  switch {
  case rc.Template != "" && rc.Action == "" && rc.Payload == "":
    err = h.handleTemplateRoute(w, r, rc)
  case rc.Method != "" && r.Method != rc.Method:
    err = h.handleMethodNotAllowed(w, r)
  default:
    err = h.handleActionRoute(w, r, rc, ps)
  }

  if err != nil {
    return err
  }
  return nil
}

func (h *RouteHandler) handleTemplateRoute(w http.ResponseWriter, r *http.Request, rc actions.HttpRoute) error {
  var err error
  log.Println("render module tpl", rc.Module, rc.Template)

  cssFiles := []string{
    "theme/assets/lib/normalize.min.css",
    "theme/assets/lib/semantic.min.css",
    "theme/assets/lib/fontawesome.min.css",
    "theme/assets/lib/fonts.css",
    "theme/assets/opus.css",
  }
  
  jsFiles := []string{
    "theme/assets/lib/polyfill.min.js",
    "theme/assets/lib/jquery-3.6.0.min.js",
    "theme/assets/lib/d3.min.js",
    "theme/assets/lib/semantic.min.js",
    "theme/assets/lib/rxjs.umd.min.js",
    "theme/assets/lib/petite-vue.iife.js",
    "theme/assets/lib/director.min.js",
    "theme/assets/services.js",
  }
  theme := h.container.Theme()
  err = theme.LoadAssets(cssFiles, jsFiles)
  if err != nil {
    return err
  }

  data := map[string]interface{}{
    "Title":  "Opus",
    "Route":  rc,
  }

  err = theme.Render(w, rc.Document, rc.Module, rc.Template, data)
  if err != nil {
    return err
  }
  return nil
}

func (h *RouteHandler) handleActionRoute(w http.ResponseWriter, r *http.Request, rc actions.HttpRoute, ps httprouter.Params) error {
  payload, err := h.bindPayload(w, r, rc, ps)
  if err != nil {
    return err
  }

  err = h.hub.Publish("before:action", payload)
  if err != nil {
    return err
  }
  err = h.hub.Publish("before:" + rc.Action, payload)
  if err != nil {
    return err
  }
  registry := h.container.Registry()
  actionFunc, ok := registry.ResolveAction(rc.Action)
  if !ok {
    return quality.ErrInvalidAction
  }

  action := actions.NewAction(rc.Action, payload, h.container, actionFunc)
  if err := action.Execute(); err != nil {
    h.hub.Publish("error:" + rc.Action, payload)
    h.hub.Publish("error:action", payload)
    return err
  }

  err = h.hub.Publish("success:"+rc.Action, payload)
  if err != nil {
    return err
  }
  err = h.hub.Publish("success:action", payload)
  if err != nil {
    return err
  }

  err = h.handleResponse(w, r, rc, payload)
  if err != nil {
    return err
  }
  return nil
}

func (h *RouteHandler) bindPayload(w http.ResponseWriter, r *http.Request, rc actions.HttpRoute, ps httprouter.Params) (actions.Payload, error) {
  if r.Body == nil {
    //TODO: move to quality
    return nil, errors.New("empty request body")
  }
  defer r.Body.Close()

  registry := h.container.Registry()
  payloadCtor, ok := registry.ResolvePayload(rc.Payload)
  if !ok {
    http.Error(w, "Invalid payload", http.StatusInternalServerError)
    return nil, quality.ErrInvalidPayload
  }

  payload := payloadCtor()  
  if err := payload.FromRequest(r); err != nil {
    writeJSONError(w, http.StatusBadRequest, "Error parsing request", err.Error())
    return nil, quality.ErrValidation
  }

  if rc.Method != "GET" && rc.Method != "DELETE" {
    if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
      http.Error(w, "Invalid payload format", http.StatusBadRequest)
      return nil, quality.ErrPayloadBinding
    }

    if err := payload.Validate(); err != nil {
      writeJSONError(w, http.StatusBadRequest, "Validation failed", err.Error())
      return nil, quality.ErrValidation
    }
  }
  return payload, nil
}

func (h *RouteHandler) handleTemplateResponse(w http.ResponseWriter, r *http.Request, rc actions.HttpRoute, p actions.Payload) error {
  var err error
  cssFiles := []string{
    "theme/assets/lib/normalize.min.css",
    "theme/assets/lib/semantic.min.css",
    "theme/assets/lib/fontawesome.min.css",
    "theme/assets/lib/fonts.css",
    "theme/assets/opus.css",
  }
  
  jsFiles := []string{
    "theme/assets/lib/polyfill.min.js",
    "theme/assets/lib/jquery-3.6.0.min.js",
    "theme/assets/lib/d3.min.js",
    "theme/assets/lib/semantic.min.js",
    "theme/assets/lib/rxjs.umd.min.js",
    "theme/assets/lib/petite-vue.iife.js",
    "theme/assets/lib/director.min.js",
    "theme/assets/services.js",
  }

  theme := h.container.Theme()
  err = theme.LoadAssets(cssFiles, jsFiles)
  if err != nil {
    return err
  }

  data := map[string]interface{}{
    "Title":  "Opus",
    "Route":  rc,
  }

  err = theme.Render(w, rc.Document, rc.Module, rc.Template, data)
  if err != nil {
    // TODO: move to quality
    return fmt.Errorf("failed to render template: %w", err)
  }

  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  w.WriteHeader(http.StatusOK)
  return nil
}

func (h *RouteHandler) handleResponse(w http.ResponseWriter, r *http.Request, rc actions.HttpRoute, p actions.Payload) error {
  if rc.Template != "" {
    return h.handleTemplateResponse(w, r, rc, p)
  }
  return writeJSONResponse(w, http.StatusOK, p)
}

func (h *RouteHandler) handleError(w http.ResponseWriter, r *http.Request, err error, rc actions.HttpRoute) {
  var status int
  // TODO: review the shit out of this
  switch {
  case errors.Is(err, quality.ErrPayloadBinding), errors.Is(err, quality.ErrValidation), errors.Is(err, quality.ErrValidation):
    status = http.StatusBadRequest
  case errors.Is(err, quality.ErrInvalidAction), errors.Is(err, quality.ErrInvalidPayload):
    status = http.StatusInternalServerError
  default:
    status = http.StatusInternalServerError
  }

  if r.Method == http.MethodGet && rc.Template != "" {
    h.handleHTMLError(w, status, err.Error())
  } else {
    writeJSONError(w, status, "Request failed", err.Error())
  }
}

func (h *RouteHandler) handleHTMLError(w http.ResponseWriter, status int, message string) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  data := map[string]interface{}{
    "Title":  "Opus",
  }
  theme := h.container.Theme()
  err := theme.Render(w, "", "", "", data)
  if err != nil {
    writeJSONError(w, status, "Request failed", err.Error())
  }
  w.WriteHeader(status)
}

func (h *RouteHandler) handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) error {
  w.Header().Set("Allow", "GET, POST, PUT, DELETE")
  w.WriteHeader(http.StatusMethodNotAllowed)
  return nil
}

// Helper functions
func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) error {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(status)
  return json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message, detail string) {
  writeJSONResponse(w, status, map[string]interface{}{
    "error":   message,
    "details": detail,
  })
}
