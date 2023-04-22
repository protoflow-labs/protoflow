package cli

import (
	"github.com/breadchris/protoflow/pkg/api"
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
	httpHandler *api.HTTPServer,
	grpcHandler *api.GRPCServer,
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
						Name:  "http",
						Usage: "Port for the http server",
					},
					&cli.IntFlag{
						Name:  "grpc",
						Usage: "Port for the grpc server",
					},
				},
				Action: func(ctx *cli.Context) error {
					httpPort := ctx.Int("http")
					if httpPort == 0 {
						httpPort = 8080
					}

					grpcPort := ctx.Int("grpc")
					if grpcPort == 0 {
						grpcPort = 50051
					}

					go func() {
						log.Info().Int("port", grpcPort).Msg("starting http server")
						if err := grpcHandler.Serve(grpcPort); err != nil {
							log.Error().Err(err).Msg("error serving grpc")
							return
						}
					}()
					log.Info().Int("port", httpPort).Msg("starting http server")
					return httpHandler.Serve(httpPort)
				},
			},
		},
	}
}
