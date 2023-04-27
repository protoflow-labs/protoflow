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
		tag := keyvals[i].(string)
		value := keyvals[i+1]
		if value.(error) != nil {
			errLog = errLog.Err(value.(error))
			continue
		}
		errLog = errLog.Str(tag, fmt.Sprintf("%+v", value))
	}
	errLog.Msg(msg)
}

var _ Logger = (*MemoryLogger)(nil)
