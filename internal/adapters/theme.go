package adapters

import (
    "fmt"
    "os"
    "sync"
    "bytes"
    "strings"
    "net/http"
    "io/ioutil"
    "html/template"
    "path/filepath"
    "golang.org/x/sync/singleflight"
    "gitlab.com/pedrokoblitz/opus-go/internal/quality"
)

// ThemeAdapter handles template loading and rendering
type ThemeAdapter struct {
    basePath    string
    document    string
    components  *template.Template
    templates   map[string]*template.Template
    templateExt string
    mu          sync.RWMutex
    assets      map[string]string // Stores concatenated assets
    loader      singleflight.Group
}

// NewThemeAdapter creates a new template adapter
func NewThemeAdapter(basePath, document string) *ThemeAdapter {
    ta := &ThemeAdapter{
        basePath:    basePath,
        document: document,
        templates:   make(map[string]*template.Template),
        templateExt: ".html", // default extension
        assets:      make(map[string]string),
    }
    return ta
}

// LoadAssets loads and concatenates CSS/JS assets for a module
func (ta *ThemeAdapter) LoadAssets(css, js []string) error {
    ta.mu.Lock()
    defer ta.mu.Unlock()

    cssTxt, err := ta.concatenateTextFiles(css)
    if err != nil {
        return quality.ErrExecution.WithDetail(fmt.Sprintf("error concatenating css files: %v", err))
    }
    ta.assets["css"] = "<style>" + cssTxt + "</style>"

    jsTxt, err := ta.concatenateTextFiles(js)
    if err != nil {
        return quality.ErrExecution.WithDetail(fmt.Sprintf("error concatenating js files: %v", err))
    }
    ta.assets["js"] = "<script>" + jsTxt + "</script>"

    return nil
}

// concatenateCSS combines CSS files with proper formatting
func (ta *ThemeAdapter) concatenateTextFiles(filePaths []string) (string, error) {
    var buffer bytes.Buffer

    for _, filePath := range filePaths {
        // Add file source comment
        filePath := filepath.Join(ta.basePath, filePath)
        buffer.WriteString(fmt.Sprintf("/* Source: %s */\n", filepath.Base(filePath)))

        content, err := ioutil.ReadFile(filePath)
        if err != nil {
            return "", quality.ErrFsIO.WithDetail(fmt.Sprintf("error reading %s: %v", filePath, err))
        }

        txt := strings.TrimSpace(string(content))
        buffer.WriteString(txt)
        buffer.WriteString("\n\n")
    }

    return buffer.String(), nil
}

// Render executes the template hierarchy
func (ts *ThemeAdapter) Render(w http.ResponseWriter, document, module, view string, data map[string]interface{}) error {
  // Get or create the template set
  if document == "" {
    document = ts.document
  }
  tpl, err := ts.loadTemplateSet(document, module, view)
  if err != nil {
    // TODO: quality err
    return quality.ErrFsIO.WithDetail(fmt.Sprintf("template load failed: %w", err))
  }

  // Prepare template data
  if data == nil {
    data = make(map[string]interface{})
  }
  
  // Inject assets
  ts.mu.Lock()
  data["RenderedCss"] = template.HTML(ts.assets["css"])
  data["RenderedJs"] = template.HTML(ts.assets["js"])
  ts.mu.Unlock()

  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  return tpl.ExecuteTemplate(w, ts.document, data)
}

func (ts *ThemeAdapter) loadTemplateSet(document, module, view string) (*template.Template, error) {
  key := fmt.Sprintf("%s:%s", module, view)

  // Singleflight protected load
  tpl, err, _ := ts.loader.Do(key, func() (interface{}, error) {

    // 1. Base document (required)
    basePath := filepath.Join(ts.basePath, "theme", document+ts.templateExt)
    tpl, err := template.New(document).Delims("<%", "%>").ParseFiles(basePath)
    if err != nil {
      // TODO: quality err
      return nil, quality.ErrFsIO.WithDetail(fmt.Sprintf("base template error: %w", err))
    }

    // 2. Module theme overrides (optional)
    themePath := filepath.Join(ts.basePath, "modules", module, "theme", "theme"+ts.templateExt)
    if _, err := os.Stat(themePath); err == nil {
      tpl, _ = tpl.ParseFiles(themePath)
    }

    // 3. Specific view template (required)
    viewPath := filepath.Join(ts.basePath, "modules", module, "theme", view+ts.templateExt)
    tpl, err = tpl.ParseFiles(viewPath)
    if err != nil {
      // TODO: quality err
      return nil, quality.ErrFsIO.WithDetail(fmt.Sprintf("view template error: %w", err))
    }

    return tpl, nil
  })

  if err != nil {
    return nil, err
  }
  return tpl.(*template.Template), nil
}
