package bufcurl

import (
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl/verbose"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
	"io"
	"os"
)

type Container struct {
}

func (c Container) Env(key string) string {
	//TODO implement me
	panic("implement me")
}

func (c Container) ForEachEnv(f func(string, string)) {
	//TODO implement me
	panic("implement me")
}

func (c Container) Stdin() io.Reader {
	//TODO implement me
	panic("implement me")
}

func (c Container) Stdout() io.Writer {
	//TODO implement me
	panic("implement me")
}

func (c Container) Stderr() io.Writer {
	return os.Stderr
}

func (c Container) NumArgs() int {
	//TODO implement me
	panic("implement me")
}

func (c Container) Arg(i int) string {
	//TODO implement me
	panic("implement me")
}

func (c Container) AppName() string {
	//TODO implement me
	panic("implement me")
}

func (c Container) ConfigDirPath() string {
	//TODO implement me
	panic("implement me")
}

func (c Container) CacheDirPath() string {
	//TODO implement me
	panic("implement me")
}

func (c Container) DataDirPath() string {
	//TODO implement me
	panic("implement me")
}

func (c Container) Port() (uint16, error) {
	//TODO implement me
	panic("implement me")
}

func (c Container) Logger() *zap.Logger {
	//TODO implement me
	panic("implement me")
}

type ZeroLogPrinter struct{}

func (z *ZeroLogPrinter) Printf(format string, args ...interface{}) {
	log.Trace().Msgf(format, args...)
}

func (c Container) VerbosePrinter() verbose.Printer {
	return &ZeroLogPrinter{}
}
