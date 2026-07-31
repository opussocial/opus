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
  "path/filepath"
  "html/template"
  
  "github.com/opussocialcontent/opus-go/actions"
  "github.com/opussocialcontent/opus-go/quality"
  "github.com/julienschmidt/httprouter"
  "golang.org/x/time/rate"
)

// AssetPaths holds the simplified asset paths
var (
  // Base directories - these are relative to the theme loader
  // The theme loader already looks in resources/ directory
  ResourcesDir = "./resources"
  StaticAssetsDir = "static"
  ConsumersDir = "consumers"
  ThemeDir = "theme"
  
  // Subdirectories relative to static
  CSSDir = StaticAssetsDir + "/css"
  JSDir = StaticAssetsDir + "/js"
  FontsDir = StaticAssetsDir + "/fonts"
  ImagesDir = StaticAssetsDir + "/images"
  
  // Template directories
  TemplatesDir = "templates"  // New location
  LegacyTemplatesDir = "theme" // Old location
)

// AssetConfig defines which assets to load for different contexts
type AssetConfig struct {
  CSSFiles []string
  JSFiles  []string
}

// GetDefaultAssets returns the default asset configuration
// These paths are relative to the resources/ directory
func GetDefaultAssets() AssetConfig {
  return AssetConfig{
    CSSFiles: []string{
      "static/css/normalize.min.css",
      "static/css/semantic.min.css",
      "static/css/fontawesome.min.css",
      "static/css/fonts.css",
      "static/css/opus.css",
    },
    JSFiles: []string{
      "static/js/polyfill.min.js",
      "static/js/jquery-3.6.0.min.js",
      "static/js/d3.min.js",
      "static/js/semantic.min.js",
      "static/js/rxjs.umd.min.js",
      "static/js/petite-vue.iife.js",
      "static/js/director.min.js",
      "static/js/services.js",
    },
  }
}

// GetLegacyAssets returns the legacy asset paths for backward compatibility
// These paths are relative to the resources/ directory
func GetLegacyAssets() AssetConfig {
  return AssetConfig{
    CSSFiles: []string{
      "theme/assets/lib/normalize.min.css",
      "theme/assets/lib/semantic.min.css",
      "theme/assets/lib/fontawesome.min.css",
      "theme/assets/lib/fonts.css",
      "theme/assets/opus.css",
    },
    JSFiles: []string{
      "theme/assets/lib/polyfill.min.js",
      "theme/assets/lib/jquery-3.6.0.min.js",
      "theme/assets/lib/d3.min.js",
      "theme/assets/lib/semantic.min.js",
      "theme/assets/lib/rxjs.umd.min.js",
      "theme/assets/lib/petite-vue.iife.js",
      "theme/assets/lib/director.min.js",
      "theme/assets/services.js",
    },
  }
}

// GetDefaultAssetPaths returns the asset paths for the theme
func GetDefaultAssetPaths() ([]string, []string) {
  config := GetDefaultAssets()
  return config.CSSFiles, config.JSFiles
}

type HTTPService struct {
  container   actions.Container
  hub         *PubSubHub
  server      *http.Server
  handler     *RouteHandler
  router      *httprouter.Router
}

func NewHTTPService(
  container actions.Container,
  hub *PubSubHub,
) *HTTPService {

  router := httprouter.New()

  router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "Not found", http.StatusNotFound)
  })

  router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
  })

  cm := quality.NewConnectionManager(quality.ConnectionConfig{
    MaxConns:     1000,
    IdleTimeout:  time.Minute * 30,
    WriteTimeout: time.Second * 10,
  })

  rl := quality.NewRateLimiter(
    rate.Limit(100),
    50,
    time.Minute*10,
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
    container: container,
    hub:      hub,
    router:   router,
  }

  service.registerRoutes(router)

  return service
}

func (s *HTTPService) registerRoutes(router *httprouter.Router) {
  // Serve static files from the resources/static directory
  staticDir := filepath.Join(ResourcesDir, "static")
  if _, err := os.Stat(staticDir); err == nil {
    fs := http.FileServer(http.Dir(staticDir))
    // Mount at /static/
    router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static/", fs))
    log.Printf("Serving static files from %s at /static/", staticDir)
    
    // Also support old /assets/ path for backward compatibility
    router.Handler(http.MethodGet, "/assets/*filepath", http.StripPrefix("/assets/", fs))
    log.Printf("Serving static files from %s at /assets/ (backward compatibility)", staticDir)
  }

  // Register providers
  cfg := s.container.Config()
  providers := cfg.Providers
  for _, provider := range providers {
    // Look for provider.yml in the resources/consumers/ directory
    providerPath := filepath.Join(ResourcesDir, ConsumersDir, provider, "provider.yml")
    
    // Check if provider exists
    if _, err := os.Stat(providerPath); os.IsNotExist(err) {
      log.Printf("Provider config not found at %s", providerPath)
      continue
    }

    defaultProviderConfig := &actions.ProviderConfig{
        Provider: actions.ProviderDefinition{
            Name:    provider,
            Depends: []string{},
            Http: actions.HttpConfig{
                Api: actions.HttpRouteGroup{
                    Protected: false,
                    Prefix:    "/api",
                },
                Views: actions.HttpRouteGroup{
                    Protected: false,
                    Prefix:    "/",
                    Document:  "page",
                },
            },
        },
    }
    
    providerConfig, err := actions.LoadProviderConfig(providerPath, defaultProviderConfig)
    if err != nil {
      log.Printf("Failed to load provider config for %s, using defaults: %v", provider, err)
      providerConfig = defaultProviderConfig
    }

    for _, sub := range providerConfig.Provider.Subscribers {
      opts := SubscriberOptions{
        Concurrency: 1,
        RetryPolicy: RetryPolicy{
          MaxAttempts: 2,
          Backoff: 10,
        },
      }
      log.Println("SUBSCRIBED!!!")
      if s.hub != nil {
        s.hub.Subscribe(sub.Topic, sub.Action, opts)
      }
    }

    for _, rc := range providerConfig.Provider.Http.Api.Routes {
      rc.Provider = providerConfig.Provider.Name
      fullPath := path.Join("/", providerConfig.Provider.Http.Api.Prefix, rc.Pattern)
      
      method := "GET"
      if rc.Method != "" {
        method = rc.Method
      }
      
      router.Handle(method, fullPath, s.handler.HandleHTTPRouter(rc))
      log.Println("api route:", fullPath, rc.Method)
    }

    for _, rc := range providerConfig.Provider.Http.Views.Routes {
      rc.Provider = providerConfig.Provider.Name

      if rc.Document == "" {
        rc.Document = providerConfig.Provider.Http.Views.Document
      }
      
      fullPath := path.Join("/", providerConfig.Provider.Http.Views.Prefix, rc.Pattern)
      router.Handle(http.MethodGet, fullPath, s.handler.HandleHTTPRouter(rc))
      log.Println("view route:", fullPath, rc.Method)
    }
  }
}

func (s *HTTPService) Start() error {
  return s.server.ListenAndServe()
}

func (s *HTTPService) Shutdown(ctx context.Context) error {
  if s.handler.rl != nil {
    s.handler.rl.Stop()
  }
  
  if s.handler.cm != nil {
    s.handler.cm.Close()
  }
  return s.server.Shutdown(ctx)
}

type RouteHandler struct {
  hub       *PubSubHub
  container actions.Container
  cm        *quality.ConnectionManager
  rl        *quality.RateLimiter
  // Cache for template functions to avoid re-parsing on each request
  templateFuncs template.FuncMap
}

func NewRouteHandler(
  hub *PubSubHub,
  container actions.Container,
  cm *quality.ConnectionManager,
  rl *quality.RateLimiter,
) *RouteHandler {
  return &RouteHandler{
    hub:       hub,
    container: container,
    cm:        cm,
    rl:        rl,
    templateFuncs: template.FuncMap{
      // Add any custom template functions here
      "safeHTML": func(s string) template.HTML {
        return template.HTML(s)
      },
    },
  }
}

func (h *RouteHandler) HandleHTTPRouter(rc actions.HttpRoute) httprouter.Handle {
  return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    if !h.rl.GetVisitor(r.RemoteAddr).Allow() {
      http.Error(w, "too many requests", http.StatusTooManyRequests)
      return
    }

    h.cm.Acquire()
    defer h.cm.Release()

    err := h.handleRoute(w, r, rc, ps)
    if err != nil {
      h.handleError(w, r, err, rc)
    }
  }
}

type SimpleClaims struct {
  UserID string
}

func (h *RouteHandler) handleRoute(w http.ResponseWriter, r *http.Request, rc actions.HttpRoute, ps httprouter.Params) error {
  ctx := r.Context()
  ctx = context.WithValue(ctx, "story", ps.ByName("story"))
  ctx = context.WithValue(ctx, "definition", ps.ByName("definition"))

  if rc.Protected == true {
    tokenHeader := r.Header.Get("Authorization")
    token := strings.ReplaceAll(tokenHeader, "Bearer: ", "")
    
    claims := &SimpleClaims{UserID: "anonymous"}
    if token != "" {
      claims = &SimpleClaims{UserID: "user_" + token}
    }
    
    ctx = context.WithValue(ctx, "authToken", token)
    ctx = context.WithValue(ctx, "userID", claims.UserID)
  }
  r = r.WithContext(ctx)

  switch {
  case rc.Template != "" && rc.Action == "" && rc.Payload == "":
    err := h.handleTemplateRoute(w, r, rc)
    if err != nil {
      return err
    }
  case rc.Method != "" && r.Method != rc.Method:
    return h.handleMethodNotAllowed(w, r)
  default:
    err := h.handleActionRoute(w, r, rc, ps)
    if err != nil {
      return err
    }
  }

  return nil
}

// findTemplatePath tries to find the template in various locations
func (h *RouteHandler) findTemplatePath(provider, template string) (string, error) {
  // Try new location: resources/consumers/{provider}/templates/{template}.html
  newPath := filepath.Join(ResourcesDir, ConsumersDir, provider, TemplatesDir, template+".html")
  if _, err := os.Stat(newPath); err == nil {
    log.Printf("Found template at: %s", newPath)
    return newPath, nil
  }
  
  // Try new location without .html extension (if template already includes it)
  newPathNoExt := filepath.Join(ResourcesDir, ConsumersDir, provider, TemplatesDir, template)
  if _, err := os.Stat(newPathNoExt); err == nil {
    log.Printf("Found template at: %s", newPathNoExt)
    return newPathNoExt, nil
  }
  
  // Try legacy location: resources/consumers/{provider}/theme/{template}.html
  legacyPath := filepath.Join(ResourcesDir, ConsumersDir, provider, LegacyTemplatesDir, template+".html")
  if _, err := os.Stat(legacyPath); err == nil {
    log.Printf("Found template at (legacy): %s", legacyPath)
    return legacyPath, nil
  }
  
  // Try legacy location without .html extension
  legacyPathNoExt := filepath.Join(ResourcesDir, ConsumersDir, provider, LegacyTemplatesDir, template)
  if _, err := os.Stat(legacyPathNoExt); err == nil {
    log.Printf("Found template at (legacy): %s", legacyPathNoExt)
    return legacyPathNoExt, nil
  }
  
  return "", fmt.Errorf("template not found: %s (tried: %s, %s, %s, %s)", 
    template, newPath, newPathNoExt, legacyPath, legacyPathNoExt)
}

func (h *RouteHandler) handleTemplateRoute(w http.ResponseWriter, r *http.Request, rc actions.HttpRoute) error {
  log.Println("render provider tpl", rc.Provider, rc.Template)

  // Load assets using paths relative to resources/
  cssFiles, jsFiles := GetDefaultAssetPaths()
  
  theme := h.container.Theme()
  
  // Try to load from simplified structure first
  err := theme.LoadAssets(cssFiles, jsFiles)
  if err != nil {
    // Fallback to legacy structure
    log.Printf("Failed to load assets from simplified structure, trying legacy: %v", err)
    legacyAssets := GetLegacyAssets()
    err = theme.LoadAssets(legacyAssets.CSSFiles, legacyAssets.JSFiles)
    if err != nil {
      return fmt.Errorf("failed to load assets: %w", err)
    }
  }

  // Prepare template data
  data := map[string]interface{}{
    "Title":  "Opus",
    "Route":  rc,
    "Lang":   "en",
  }

  // Determine the document template from the route config
  // Default to "page" if not specified
  document := "page"
  if rc.Document != "" {
    document = rc.Document
  }

  // Use the ThemeAdapter's Render method with the dynamic document
  err = theme.Render(w, document, rc.Provider, rc.Template, data)
  if err != nil {
    log.Printf("Failed to render template: %v", err)
    return fmt.Errorf("failed to render template: %w", err)
  }
  
  log.Printf("Successfully rendered template: %s with document: %s", rc.Template, document)
  return nil
}

func (h *RouteHandler) handleTemplateResponse(w http.ResponseWriter, r *http.Request, rc actions.HttpRoute, p actions.Payload) error {
  // Load assets using paths relative to resources/
  cssFiles, jsFiles := GetDefaultAssetPaths()
  
  theme := h.container.Theme()
  
  // Try to load from simplified structure first
  err := theme.LoadAssets(cssFiles, jsFiles)
  if err != nil {
    // Fallback to legacy structure
    log.Printf("Failed to load assets from simplified structure, trying legacy: %v", err)
    legacyAssets := GetLegacyAssets()
    err = theme.LoadAssets(legacyAssets.CSSFiles, legacyAssets.JSFiles)
    if err != nil {
      return fmt.Errorf("failed to load assets: %w", err)
    }
  }

  // Prepare template data with payload
  data := map[string]interface{}{
    "Title":  "Opus",
    "Route":  rc,
    "Data":   p,
    "Lang":   "en",
  }

  // Determine the document template from the route config
  document := "page"
  if rc.Document != "" {
    document = rc.Document
  }

  // Use the ThemeAdapter's Render method with the dynamic document
  err = theme.Render(w, document, rc.Provider, rc.Template, data)
  if err != nil {
    log.Printf("Failed to render template: %v", err)
    return fmt.Errorf("failed to render template: %w", err)
  }

  return nil
}

func (h *RouteHandler) handleActionRoute(w http.ResponseWriter, r *http.Request, rc actions.HttpRoute, ps httprouter.Params) error {
  payload, err := h.bindPayload(w, r, rc, ps)
  if err != nil {
    return err
  }

  if h.hub != nil {
    h.hub.Publish("before:action", payload)
    h.hub.Publish("before:"+rc.Action, payload)
  }
  
  registry := h.container.Registry()
  actionFunc, ok := registry.ResolveAction(rc.Action)
  if !ok {
    return quality.ErrInvalidAction
  }

  action := actions.NewAction(rc.Action, payload, h.container, actionFunc)
  if err := action.Execute(); err != nil {
    if h.hub != nil {
      h.hub.Publish("error:"+rc.Action, payload)
      h.hub.Publish("error:action", payload)
    }
    return err
  }

  if h.hub != nil {
    h.hub.Publish("success:"+rc.Action, payload)
    h.hub.Publish("success:action", payload)
  }

  err = h.handleResponse(w, r, rc, payload)
  if err != nil {
    return err
  }
  return nil
}

func (h *RouteHandler) bindPayload(w http.ResponseWriter, r *http.Request, rc actions.HttpRoute, ps httprouter.Params) (actions.Payload, error) {
  if r.Body == nil {
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

func (h *RouteHandler) handleResponse(w http.ResponseWriter, r *http.Request, rc actions.HttpRoute, p actions.Payload) error {
  if rc.Template != "" {
    return h.handleTemplateResponse(w, r, rc, p)
  }
  return writeJSONResponse(w, http.StatusOK, p)
}

func (h *RouteHandler) handleError(w http.ResponseWriter, r *http.Request, err error, rc actions.HttpRoute) {
  var status int
  switch {
  case errors.Is(err, quality.ErrPayloadBinding), errors.Is(err, quality.ErrValidation):
    status = http.StatusBadRequest
  case errors.Is(err, quality.ErrInvalidAction), errors.Is(err, quality.ErrInvalidPayload):
    status = http.StatusInternalServerError
  default:
    status = http.StatusInternalServerError
  }

  // Only write header if not already written
  if r.Method == http.MethodGet && rc.Template != "" {
    h.handleHTMLError(w, status, err.Error())
  } else {
    writeJSONError(w, status, "Request failed", err.Error())
  }
}

func (h *RouteHandler) handleHTMLError(w http.ResponseWriter, status int, message string) {
  // Check if headers are already written
  if w.Header().Get("Content-Type") == "" {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
  }
  
  data := map[string]interface{}{
    "Title":  "Opus",
    "Error":  message,
  }
  
  // Try to find an error template
  errorTemplatePath, err := h.findTemplatePath("site", "error")
  if err == nil {
    tmpl, err := template.ParseFiles(errorTemplatePath)
    if err == nil {
      tmpl.Execute(w, data)
      w.WriteHeader(status)
      return
    }
  }
  
  // Fallback to simple HTML error
  w.WriteHeader(status)
  html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <title>Error %d</title>
  <style>
    body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
    h1 { color: #d32f2f; }
    .error { background: #ffebee; padding: 20px; border-radius: 4px; border-left: 4px solid #d32f2f; }
  </style>
</head>
<body>
  <h1>Error %d</h1>
  <div class="error">%s</div>
</body>
</html>`, status, status, message)
  w.Write([]byte(html))
}

func (h *RouteHandler) handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) error {
  w.Header().Set("Allow", "GET, POST, PUT, DELETE")
  w.WriteHeader(http.StatusMethodNotAllowed)
  return nil
}

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