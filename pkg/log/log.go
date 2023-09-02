package log

import (
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

var (
	once        = &sync.Once{}
	ProviderSet = wire.NewSet(NewLog, NewConfig)
)

type Log struct{}

// NewLog creates a new Log.
func NewLog(config Config) *Log {
	once.Do(func() {
		logLevel := zerolog.InfoLevel
		if config.Level == "debug" {
			logLevel = zerolog.DebugLevel
		}
		log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(logLevel)
	})
	return &Log{}
}
