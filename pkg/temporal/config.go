package temporal

import "go.uber.org/config"

type Config struct {
	Namespace string `yaml:"namespace"`
	Host      string `yaml:"host"`
}

func NewDefaultConfig() Config {
	return Config{
		Namespace: "${TEMPORAL_NAMESPACE:\"protoflow\"}",
		Host:      "${TEMPORAL_HOST:\"localhost:7233\"}",
	}
}

func NewConfig(provider config.Provider) (Config, error) {
	var c Config
	err := provider.Get("temporal").Populate(&c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}
