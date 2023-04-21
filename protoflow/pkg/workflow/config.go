package workflow

import "go.uber.org/config"

type TemporalConfig struct {
	Namespace string `yaml:"namespace"`
	Host      string `yaml:"host"`
}

type Config struct {
	Temporal TemporalConfig `yaml:"temporal"`
}

func NewDefaultConfig() Config {
	return Config{
		Temporal: TemporalConfig{
			Namespace: "${TEMPORAL_NAMESPACE:\"protoflow\"}",
			Host:      "${TEMPORAL_HOST:\"localhost:7233\"}",
		},
	}
}

func NewConfig(provider config.Provider) (Config, error) {
	var c Config
	err := provider.Get("workflow").Populate(&c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}
