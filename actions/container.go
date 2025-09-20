package actions

import (
	"log"
	"database/sql"
	"gitlab.com/pedrokoblitz/opus-go/actions/adapters"
)

type Container struct {
	cfg *ServiceConfig
	db *sql.DB
	registry *ActionRegistry
	theme *adapters.ThemeAdapter
	email *adapters.EmailAdapter
	mu           sync.RWMutex
}

func NewContainer(cfg *ServiceConfig) Container {
    registry := NewActionRegistry()
	db, err := adapters.NewSQLAdapter(cfg.Dsn())
	if err != nil {
		log.Fatal("invalid db conf")
	}
	email := adapters.NewEmailAdapter(cfg.Service.Smtp)
	theme := adapters.NewThemeAdapter("./resources/", "page")
	// theme := adapters.NewThemeAdapter(cfg.Service.ResourcesDir, cfg.Service.DefaultTemplate)

	return Container{
		cfg: cfg,
		registry: registry,
		db: db,
		email: email,
		theme: theme,
	}
}

func (c Container) Get(name string) interface{} {
	r.regMu.Lock()
	defer r.regMu.Unlock()
	switch name {
	case "db":
		return c.db
	case "email":
		return c.email
	case "tpl":
		return c.theme
	case "registry":
		return c.registry
	case "cfg":
		return c.cfg
	default:
		log.Fatal("invalid dep")
	}
	return nil
}
