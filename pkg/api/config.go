package api

import (
	"github.com/rs/zerolog/log"
	"go.uber.org/config"
)

const ConfigurationKey = "api"

type Config struct {
	StudioProxy string `yaml:"studio_proxy"`
}

// TODO breadchris studio proxy should not be set by default only turn on when in dev mode
func NewDefaultConfig(studioProxy string) Config {
	return Config{
		StudioProxy: studioProxy,
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
