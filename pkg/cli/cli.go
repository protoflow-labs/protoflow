package cli

import (
	"os"

	"github.com/protoflow-labs/protoflow/pkg/api"
	logcfg "github.com/protoflow-labs/protoflow/pkg/log"
	"github.com/protoflow-labs/protoflow/pkg/project"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
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
	worker *workflow.Worker,
	project *project.Service,
) *cli.App {
	setupLogging(logConfig.Level)

	return &cli.App{
		Name:        "protoflow",
		Description: "Coding as easy as playing with legos.",
		Flags:       []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name: "worker",
				Action: func(ctx *cli.Context) error {
					return worker.Run()
				},
			},
			{
				Name: "serve",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "http",
						Usage: "Port for the http server",
					},
				},
				Action: func(ctx *cli.Context) error {
					httpPort := ctx.Int("http")
					if httpPort == 0 {
						httpPort = 8080
					}

					log.Info().Int("port", httpPort).Msg("starting http server")
					return httpHandler.Serve(httpPort)
				},
			},
			{
				Name: "projects",
				Action: func(ctx *cli.Context) error {
					res, err := project.GetProjects(ctx.Context, nil)
					if err != nil {
						return err
					}

					log.Info().Msgf("%+v", res)

					return nil
				},
			},
		},
	}
}
