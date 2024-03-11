package server

import (
	"github.com/rs/zerolog/log"
	"go.uber.org/config"
)

type Config struct {
	Port        int    `yaml:"port"`
	StudioProxy string `yaml:"studio_proxy"`
}

// TODO breadchris studio proxy should not be set by default only turn on when in dev mode
func NewDefaultConfig() Config {
	return Config{
		StudioProxy: "",
		Port:        8000,
	}
}

func NewConfig(config config.Provider) (cfg Config, err error) {
	err = config.Get("server").Populate(&cfg)
	if err != nil {
		log.Error().Err(err).Msg("failed loading config")
		return
	}
	return
}
