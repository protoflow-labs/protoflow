package db

import "go.uber.org/config"

type Config struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

func NewDefaultConfig() Config {
	return Config{
		Driver: "sqlite",
		DSN:    "db.sqlite",
	}
}

func NewConfig(provider config.Provider) (Config, error) {
	var c Config
	err := provider.Get("db").Populate(&c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}
