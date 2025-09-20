package services

import (
	"log"
	"gitlab.com/pedrokoblitz/opus-go/internal/adapters"
	"gitlab.com/pedrokoblitz/opus-go/internal/actions"
)

type Container struct {
	Config *ServiceConfig
	DB adapters.DatabaseAdapter
	Registry *actions.ActionRegistry
	Theme *adapters.ThemeAdapter
	Email *adapters.EmailAdapter
}

func NewContainer(cfg *ServiceConfig) Container {
    registry := NewActionRegistry()
	db, err := adapters.NewSQLAdapter(cfg.Service.DSN)
	if err != nil {
		log.Fatal("invalid db conf")
	}
	email := adapters.NewEmailAdapter(cfg.Service.Smtp)
	theme := adapters.NewThemeAdapter("./resources/", "page")
	// theme := adapters.NewThemeAdapter(cfg.Service.ResourcesDir, cfg.Service.DefaultTemplate)

	return Container{
		Config: cfg,
		Registry: registry,
		DB: db,
		Email: email,
		Theme: theme,
	}
}

func (c Container) GetDB() adapters.DatabaseAdapter {
	return c.DB
}

func (c Container) GetRegistry() *actions.ActionRegistry {
	return c.Registry
}

func (c Container) GetEmail() *adapters.EmailAdapter {
	return c.Email
}

func (c Container) GetTheme() *adapters.ThemeAdapter {
	return c.Theme
}
