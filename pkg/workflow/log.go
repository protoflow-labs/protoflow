package workflow

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"runtime"
)

type Logger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
}

type MemoryLogger struct{}

func fileLine() string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		return fmt.Sprintf("%s:%d", file, line)
	} else {
		return ""
	}
}

func (m MemoryLogger) Debug(msg string, keyvals ...interface{}) {
	log.Debug().Interface("keyvals", keyvals).Msg(msg)
}

func (m MemoryLogger) Info(msg string, keyvals ...interface{}) {
	log.Info().Str("file", fileLine()).Interface("keyvals", keyvals).Msg(msg)
}

func (m MemoryLogger) Warn(msg string, keyvals ...interface{}) {
	log.Warn().Interface("keyvals", keyvals).Msg(msg)
}

func (m MemoryLogger) Error(msg string, keyvals ...interface{}) {
	errLog := log.Error().Str("file", fileLine())
	for i := 0; i < len(keyvals); i += 2 {
		errLog = errLog.Str(keyvals[i].(string), fmt.Sprintf("%+v", keyvals[i+1]))
	}
	errLog.Msg(msg)
}

var _ Logger = (*MemoryLogger)(nil)
