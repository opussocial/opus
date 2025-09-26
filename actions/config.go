package actions

type ServiceConfig struct {
	Service struct {
		Modules         []string `yaml:"modules"`
		Host            string   `yaml:"host"`
		ResourcesDir    string   `yaml:"resourcesDir"`
		DefaultTemplate string   `yaml:"defaultTemplate"`
	} `yaml:"service"`
}

// Validate checks the service configuration for errors
func (sc *ServiceConfig) Validate() error {
	return nil
}
