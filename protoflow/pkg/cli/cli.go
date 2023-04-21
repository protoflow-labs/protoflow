package cli

import (
	logcfg "github.com/breadchris/protoflow/pkg/log"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
)

// TODO breadchris this should be a provided dependency
func setupLogging(level string) {
	logLevel := zerolog.InfoLevel
	if level == "debug" {
		logLevel = zerolog.DebugLevel
	}
	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(logLevel)
}

func New(
	logConfig logcfg.Config,
) *cli.App {
	setupLogging(logConfig.Level)

	return &cli.App{
		Name:        "protoflow",
		Description: "Coding as easy as playing with legos.",
		Flags:       []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name: "serve",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name: "port",
					},
				},
				Action: func(ctx *cli.Context) error {
					port := ctx.Int("port")
					if port != 0 {
						log.Info().Int("port", port).Msg("running on port")
					}
					return nil
				},
			},
		},
	}
}
