package generate

import "go.uber.org/config"

type Config struct {
	ProjectPath string `yaml:"project_path"`
}

func NewDefaultConfig() Config {
	return Config{
		ProjectPath: "${PROJECT_PATH:\"\"}",
	}
}

func NewConfig(provider config.Provider) (Config, error) {
	var c Config
	err := provider.Get("generate").Populate(&c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}
