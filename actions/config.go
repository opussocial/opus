package actions

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// ConfigLoader provides configuration loading with optional defaults
type ConfigLoader struct {
	defaultConfigs map[string]interface{}
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{
		defaultConfigs: make(map[string]interface{}),
	}
}

// SetDefault sets a default configuration for a specific path
func (cl *ConfigLoader) SetDefault(path string, defaultConfig interface{}) {
	cl.defaultConfigs[path] = defaultConfig
}
// LoadConfig loads configuration from file with optional default
func (cl *ConfigLoader) LoadConfig(path string, target interface{}) error {
	// 1. Apply defaults immediately (if any)
	if def, ok := cl.defaultConfigs[path]; ok && def != nil {
		switch d := def.(type) {
		case *ServiceConfig:
			*target.(*ServiceConfig) = *d
		case *ProviderConfig:
			*target.(*ProviderConfig) = *d
		default:
			return fmt.Errorf("unsupported configuration type for default")
		}
	}

	// 2. Try loading file (optional!)
	data, err := os.ReadFile(path)
	if err != nil {
		// If missing → defaults already applied → OK
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	// 3. Merge: YAML overrides defaults already set in `target`
	if err := yaml.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

type PubSubOptions struct {
	Concurrency int `yaml:"concurrency"`
	RetryPolicy struct {
		MaxAttempts int           `yaml:"max_attempts"`
		Backoff     time.Duration `yaml:"backoff"`
	} `yaml:"retry_policy"`
}

type Subscriber struct {
	ID     string `yaml:"id"`
	Action string `yaml:"action"`
	Topic  string `yaml:"topic"`
}

// Validate checks the provider configuration for errors
func (mc *ProviderConfig) Validate() error {
	if mc.Provider.Name == "" {
		return fmt.Errorf("provider name is required")
	}

	// Validate HTTP API routes
	for i, route := range mc.Provider.Http.Api.Routes {
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
	for i, route := range mc.Provider.Http.Views.Routes {
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

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type SmtpConfig struct {
	// Driver string
	Host     string
	Port     int
	Sender   string
	Password string
}

type ProviderConfig struct {
	Provider ProviderDefinition `yaml:"provider"`
}

// ProviderDefinition contains core provider metadata
type ProviderDefinition struct {
	Name        string       `yaml:"name"`
	Depends     []string     `yaml:"depends"`
	Http        HttpConfig   `yaml:"http"`
	Subscribers []Subscriber `yaml:"subscribers"`
}

// HttpConfig contains Http route definitions
type HttpConfig struct {
	Api   HttpRouteGroup `yaml:"api"`
	Views HttpRouteGroup `yaml:"views"`
}

type HttpRouteGroup struct {
	Protected bool        `yaml:"protected"`
	Document  string      `yaml:"document"`
	Template  string      `yaml:"template"`
	Prefix    string      `yaml:"prefix"`
	Routes    []HttpRoute `yaml:"routes"`
}

type HTTPRequestPayload struct {
	Method string `yaml:"method"`
	Path   string `yaml:"path"`
}

func (p *HTTPRequestPayload) Validate() error {
	return nil
}

// HttpRoute defines an Http endpoint
type HttpRoute struct {
	Document string `yaml:"document"`
	Provider   string `yaml:"provider"`

	Session   bool   `yaml:"session"`
	Protected bool   `yaml:"protected"`
	Method    string `yaml:"method"`
	Pattern   string `yaml:"pattern"`

	Payload string `yaml:"payload"`
	Action  string `yaml:"action"`

	Template string `yaml:"template"`
}

// LoadServiceConfig loads and validates service configuration with optional default
func LoadServiceConfig(path string, defaults ...*ServiceConfig) (*ServiceConfig, error) {
	loader := NewConfigLoader()
	
	if len(defaults) > 0 {
		loader.SetDefault(path, defaults[0])
	}

	var cfg ServiceConfig
	if err := loader.LoadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// LoadProviderConfig loads and validates provider configuration with optional default
func LoadProviderConfig(path string, defaults ...*ProviderConfig) (*ProviderConfig, error) {
	loader := NewConfigLoader()
	
	if len(defaults) > 0 {
		loader.SetDefault(path, defaults[0])
	}

	var cfg ProviderConfig
	if err := loader.LoadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
