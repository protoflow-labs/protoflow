package bucket

import (
	"github.com/rs/zerolog/log"
	"go.uber.org/config"
)

const ConfigurationKey = "bucket"

type Config struct {
	Name string `yaml:"name"`
}

func NewDefaultConfig() Config {
	return Config{
		Name: ".protoflow",
	}
}

func NewConfig(config config.Provider) (cfg Config, err error) {
	err = config.Get(ConfigurationKey).Populate(&cfg)
	if err != nil {
		log.Error().Err(err).Msg("failed loading config")
		return
	}
	return
}
