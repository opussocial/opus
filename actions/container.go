package actions

import (
	"errors"
	"fmt"
	"sync"
	"database/sql"

	"github.com/opussocialcontent/opus-go/actions/adapters"
)

type ServiceType string

type DBConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Minimal service config for actions.Container
type ServiceConfig struct {
	Providers []string          `yaml:"providers"`
	HTTP      HTTPConfig        `yaml:"http"`
	DB        DBConfig          `yaml:"db"`
	SMTP      adapters.EmailConfig `yaml:"smtp"`
}

func (s ServiceConfig) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		s.DB.User,
		s.DB.Password,
		s.DB.Host,
		s.DB.Port,
		s.DB.Database,
	)
}

const (
	DB       ServiceType = "db"
	Email    ServiceType = "email"
	Theme    ServiceType = "theme"
	Registry ServiceType = "registry"
	Config   ServiceType = "config"
)

var (
	ErrServiceNotFound = errors.New("service not found")
	ErrInvalidConfig   = errors.New("invalid configuration")
)

// Container manages dependency injection for the application
type Container struct {
	mu       sync.RWMutex
	services map[ServiceType]interface{}
	cfg      *ServiceConfig
	registry *ActionRegistry
	db       *sql.DB
	theme    *adapters.ThemeAdapter
	email    *adapters.EmailAdapter
}

// NewContainer creates and initializes a new Container with all dependencies
func NewContainer(cfg *ServiceConfig) (Container, error) {
	var container Container
	if cfg == nil {
		return container, fmt.Errorf("%w: config cannot be nil", ErrInvalidConfig)
	}

	// Initialize database adapter
	db, err := adapters.NewSQLAdapter(cfg.Dsn())
	if err != nil {
		return container, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize email adapter
	email := adapters.NewEmailAdapter(cfg.SMTP)

	// Initialize theme adapter
	theme := adapters.NewThemeAdapter("./resources/", "page")
	// theme := adapters.NewThemeAdapter(cfg.Service.ResourcesDir, cfg.Service.DefaultTemplate)

	registry := NewActionRegistry()

	container = Container{
		cfg:      cfg,
		db:       db,
		registry: registry,
		theme:    theme,
		email:    email,
		services: make(map[ServiceType]interface{}),
	}

	// Register all services
	container.registerServices()

	return container, nil
}

// registerServices populates the services map with all available services
func (c *Container) registerServices() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.services[DB] = c.db
	c.services[Email] = c.email
	c.services[Theme] = c.theme
	c.services[Registry] = c.registry
	c.services[Config] = c.cfg
}

// Get retrieves a service by type with type safety
func (c *Container) Get(serviceType ServiceType) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	service, exists := c.services[serviceType]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrServiceNotFound, serviceType)
	}

	return service, nil
}

// Type-safe getter methods for better developer experience

// DB returns the database adapter
func (c *Container) DB() *sql.DB {
	return c.db
}

// Email returns the email adapter
func (c *Container) Email() *adapters.EmailAdapter {
	return c.email
}

// Template returns the theme/template adapter
func (c *Container) Theme() *adapters.ThemeAdapter {
	return c.theme
}

// Registry returns the action registry
func (c *Container) Registry() *ActionRegistry {
	return c.registry
}

// Config returns the service configuration
func (c *Container) Config() *ServiceConfig {
	return c.cfg
}

// MustGet retrieves a service and panics if not found (for use in initialization)
func (c *Container) MustGet(serviceType ServiceType) interface{} {
	service, err := c.Get(serviceType)
	if err != nil {
		panic(fmt.Sprintf("required service not available: %v", err))
	}
	return service
}

// Close gracefully shuts down all services
func (c *Container) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errs []error
	if err := c.db.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close database: %w", err))
	}

	// if err := c.email.Close(); err != nil {
	// 	errs = append(errs, fmt.Errorf("failed to close email service: %w", err))
	// }

	// Clear services map
	c.services = make(map[ServiceType]interface{})

	if len(errs) > 0 {
		return fmt.Errorf("errors closing container: %v", errs)
	}

	return nil
}

// // HealthCheck returns the status of all services
// func (c *Container) HealthCheck() map[ServiceType]bool {
// 	c.mu.RLock()
// 	defer c.mu.RUnlock()

// 	health := make(map[ServiceType]bool)

// 	for serviceType := range c.services {
// 		// Add custom health checks for each service type if needed
// 		health[serviceType] = true // Default to true, implement actual checks
// 	}

// 	return health
// }
