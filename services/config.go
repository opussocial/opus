package services

import (
	"fmt"
	"os"
	"time"
	"gopkg.in/yaml.v3"
  // "database/sql"
  "gitlab.com/pedrokoblitz/opus-go/internal/actions"
  "gitlab.com/pedrokoblitz/opus-go/internal/adapters"
)

func NewActionRegistry() *actions.ActionRegistry {
  return actions.NewActionRegistry()
}

func NewSQLAdapter(dsn string) (adapters.DatabaseAdapter, error) {
	return adapters.NewSQLAdapter(dsn)
}

func NewThemeAdapter(basePath, document string) *adapters.ThemeAdapter {
  return adapters.NewThemeAdapter(basePath, document)
}

type PubSubOptions struct {
  Concurrency int `yaml:"concurrency"`
  RetryPolicy struct {
    MaxAttempts int           `yaml:"max_attempts"`
    Backoff     time.Duration `yaml:"backoff"`
  } `yaml:"retry_policy"`
}

type Subscriber struct {
  ID      string `yaml:"id"`
  Action  string `yaml:"action"`
  Topic  string `yaml:"topic"`
}

type ServiceConfig struct {
    Service struct {
        Modules []string `yaml:"modules"`
        Host    string   `yaml:"host"`
        ResourcesDir    string   `yaml:"resourcesDir"`
        DefaultTemplate    string   `yaml:"defaultTemplate"`
        HTTP    struct {
            Port int `yaml:"port"`
        } `yaml:"http"`
        // TODO: move to env
        Config struct {
            JwtSecret  string `yaml:"jwt-secret"`
            SigningKey string `yaml:"signing-key"`
        } `yaml:"config"`
        Database DatabaseConfig `yaml:"database"`
        Smtp adapters.EmailConfig `yaml:"smtp"`
    } `yaml:"service"`
}
// Validate checks the service configuration for errors
func (sc *ServiceConfig) Validate() error {
  if sc.Service.Host == "" {
    return fmt.Errorf("service host is required")
  }

  if sc.Service.HTTP.Port == 0 {
    return fmt.Errorf("HTTP port is required")
  }

  if sc.Service.Config.JwtSecret == "" {
    return fmt.Errorf("JWT secret is required")
  }

  if sc.Service.Config.SigningKey == "" {
    return fmt.Errorf("signing key is required")
  }

  // Validate database configuration
  if sc.Service.Database.Driver == "" {
    return fmt.Errorf("database driver is required")
  }
  if sc.Service.Database.Host == "" {
    return fmt.Errorf("database host is required")
  }
  if sc.Service.Database.Port == 0 {
    return fmt.Errorf("database port is required")
  }
  if sc.Service.Database.User == "" {
    return fmt.Errorf("database user is required")
  }
  if sc.Service.Database.Database == "" {
    return fmt.Errorf("database name is required")
  }
  return nil
}

// Validate checks the module configuration for errors
func (mc *ModuleConfig) Validate() error {
  if mc.Module.Name == "" {
    return fmt.Errorf("module name is required")
  }

  // Validate HTTP API routes
  for i, route := range mc.Module.Http.Api.Routes {
    if route.Payload == "" && route.Action == "" {
      return fmt.Errorf("view route %d: at least one of payload, action or template is required", i+1)
    }    
    if route.Method == "" {
      return fmt.Errorf("API route %d: method is required", i+1)
    }
    if route.Pattern == "" {
      return fmt.Errorf("API route %d: pattern is required", i+1)
    }
    if route.Action == "" {
      return fmt.Errorf("API route %d: action is required", i+1)
    }
  }

  // Validate HTTP View routes
  for i, route := range mc.Module.Http.Views.Routes {
    if route.Template == "" {
      return fmt.Errorf("view route %d: at least one template is required", i+1)
    }
    if route.Method == "" {
      route.Method = "GET"
    }
    if route.Pattern == "" {
      return fmt.Errorf("view route %d: pattern is required", i+1)
    }
  }
  return nil
}

func (s ServiceConfig)	Dsn() string {	
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"root",
		"localhost",
		3307,
		"opus_test",
	)
}

type DatabaseConfig struct {
  Driver string
  Host string
  Port int
  User string
  Password string
  Database string
}

type SmtpConfig struct {
  // Driver string
  Host string
  Port int
  Sender string
  Password string
}

type ModuleConfig struct {
    Module ModuleDefinition `yaml:"module"`
}

// ModuleDefinition contains core module metadata
type ModuleDefinition struct {
	Name    string   `yaml:"name"`
	Depends []string `yaml:"depends"`
	Http HttpConfig `yaml:"http"`
  Subscribers []Subscriber `yaml:"subscribers"`
}

// HttpConfig contains Http route definitions
type HttpConfig struct {
	Api HttpRouteGroup     `yaml:"api"`
	Views HttpRouteGroup `yaml:"views"`
}

type HttpRouteGroup struct {
  Protected bool `yaml:"protected"`
  Document string     `yaml:"document"`
  Template string     `yaml:"template"`
  Prefix string     `yaml:"prefix"`
	Routes []HttpRoute `yaml:"routes"`
}

type HTTPRequestPayload struct {
	Method string     `yaml:"method"`
	Path string     `yaml:"path"`
}

func (p *HTTPRequestPayload) Validate() error {
	return nil
}

// HttpRoute defines an Http endpoint
type HttpRoute struct {
  Document string     `yaml:"document"`
	Module  string `yaml:"module"`

  Session bool `yaml:"session"`
  Protected bool `yaml:"protected"`
	Method  string `yaml:"method"`
	Pattern string `yaml:"pattern"`

	Payload string `yaml:"payload"`
	Action  string `yaml:"action"`

	Template string `yaml:"template"`
}

// LoadServiceConfig loads and validates module configuration
func LoadServiceConfig(path string) (*ServiceConfig, error) {
	var cfg ServiceConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read service config: %w", err)
	}

	yaml.Unmarshal(data, &cfg)

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid service config: %w", err)
	}

	return &cfg, nil
}

// LoadModuleConfig loads and validates module configuration
func LoadModuleConfig(path string) (*ModuleConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read module config: %w", err)
	}

	var cfg ModuleConfig
	yaml.Unmarshal(data, &cfg)

	// fmt.Println(cfg)
	// if err := cfg.Validate(); err != nil {
	// 	return nil, fmt.Errorf("invalid module config: %w", err)
	// }

	return &cfg, nil
}
