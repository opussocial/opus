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
    "github.com/opussocialcontent/opus-go/quality"
)

// ThemeAdapter handles template loading and rendering
type ThemeAdapter struct {
    basePath    string
    document    string  // Default document name (e.g., "page")
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
        document:    document,  // Store the default document name
        templates:   make(map[string]*template.Template),
        templateExt: ".html",
        assets:      make(map[string]string),
    }
    return ta
}

// LoadAssets loads and concatenates CSS/JS assets for a consumer
func (ta *ThemeAdapter) LoadAssets(css, js []string) error {
    ta.mu.Lock()
    defer ta.mu.Unlock()

    if len(css) > 0 {
        cssTxt, err := ta.concatenateTextFiles(css)
        if err != nil {
            return quality.ErrExecution.WithDetail(fmt.Sprintf("error concatenating css files: %v", err))
        }
        ta.assets["css"] = "<style>" + cssTxt + "</style>"
    } else {
        ta.assets["css"] = ""
    }

    if len(js) > 0 {
        jsTxt, err := ta.concatenateTextFiles(js)
        if err != nil {
            return quality.ErrExecution.WithDetail(fmt.Sprintf("error concatenating js files: %v", err))
        }
        ta.assets["js"] = "<script>" + jsTxt + "</script>"
    } else {
        ta.assets["js"] = ""
    }

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
func (ta *ThemeAdapter) Render(w http.ResponseWriter, document, consumer, view string, data map[string]interface{}) error {
    // Get or create the template set
    if document == "" {
        document = ta.document
    }
    tpl, err := ta.loadTemplateSet(document, consumer, view)
    if err != nil {
        return quality.ErrFsIO.WithDetail(fmt.Sprintf("template load failed: %v", err))
    }

    // Prepare template data
    if data == nil {
        data = make(map[string]interface{})
    }
    
    // Inject assets
    ta.mu.RLock()
    data["RenderedCss"] = template.HTML(ta.assets["css"])
    data["RenderedJs"] = template.HTML(ta.assets["js"])
    ta.mu.RUnlock()

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    
    // Always execute the "page" template
    // This works because all document templates define "page"
    return tpl.ExecuteTemplate(w, "page", data)
}

func (ta *ThemeAdapter) loadTemplateSet(document, consumer, view string) (*template.Template, error) {
    key := fmt.Sprintf("%s:%s:%s", document, consumer, view)

    // Singleflight protected load
    tpl, err, _ := ta.loader.Do(key, func() (interface{}, error) {
        // 1. Base document (dynamic filename) - looks for whatever document name is passed
        // But always defines it as "page" template
        docPaths := []string{
            filepath.Join(ta.basePath, "theme", document+ta.templateExt),
            filepath.Join(ta.basePath, "consumers", consumer, "theme", document+ta.templateExt),
            filepath.Join(ta.basePath, "consumers", consumer, "templates", document+ta.templateExt),
        }
        
        var docPath string
        for _, path := range docPaths {
            if _, err := os.Stat(path); err == nil {
                docPath = path
                break
            }
        }
        
        if docPath == "" {
            return nil, quality.ErrFsIO.WithDetail(fmt.Sprintf("document template file not found: %s (tried: %v)", document, docPaths))
        }
        
        // Parse the base document template
        // IMPORTANT: Use "page" as the template name, not the document variable
        tpl, err := template.New("page").Delims("{{", "}}").ParseFiles(docPath)
        if err != nil {
            return nil, quality.ErrFsIO.WithDetail(fmt.Sprintf("document template error: %v", err))
        }

        // 2. THEME OVERRIDE (REQUIRED) - theme.html
        themePaths := []string{
            filepath.Join(ta.basePath, "consumers", consumer, "theme", "theme"+ta.templateExt),
            filepath.Join(ta.basePath, "consumers", consumer, "templates", "theme"+ta.templateExt),
        }
        
        var themePath string
        for _, path := range themePaths {
            if _, err := os.Stat(path); err == nil {
                themePath = path
                break
            }
        }
        
        if themePath == "" {
            return nil, quality.ErrFsIO.WithDetail(fmt.Sprintf("theme template not found for consumer: %s", consumer))
        }

        // Parse the theme template
        tpl, err = tpl.ParseFiles(themePath)
        if err != nil {
            return nil, quality.ErrFsIO.WithDetail(fmt.Sprintf("theme template error: %v", err))
        }

        // 3. View template (required) - home.html, etc.
        viewPaths := []string{
            filepath.Join(ta.basePath, "consumers", consumer, "templates", view+ta.templateExt),
            filepath.Join(ta.basePath, "consumers", consumer, "theme", view+ta.templateExt),
        }
        
        var viewPath string
        for _, path := range viewPaths {
            if _, err := os.Stat(path); err == nil {
                viewPath = path
                break
            }
        }
        
        if viewPath == "" {
            return nil, quality.ErrFsIO.WithDetail(fmt.Sprintf("view template not found: %s", view))
        }

        // Parse the view template
        tpl, err = tpl.ParseFiles(viewPath)
        if err != nil {
            return nil, quality.ErrFsIO.WithDetail(fmt.Sprintf("view template error: %v", err))
        }

        return tpl, nil
    })

    if err != nil {
        return nil, err
    }
    return tpl.(*template.Template), nil
}