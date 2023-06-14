package cli

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/protoflow-labs/protoflow/pkg/temporal"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
	"go.uber.org/config"

	"github.com/breadchris/air/runner"
	"github.com/protoflow-labs/protoflow/pkg/api"
	logcfg "github.com/protoflow-labs/protoflow/pkg/log"
	"github.com/protoflow-labs/protoflow/pkg/project"
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

func liveReload() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var err error
	cfg, err := runner.InitConfig(cfgPath)
	if err != nil {
		return err
	}
	cfg.WithArgs(cmdArgs)
	r, err := runner.NewEngineWithConfig(cfg, debugMode)
	if err != nil {
		return err
	}
	go func() {
		<-sigs
		r.Stop()
	}()

	defer func() {
		if e := recover(); e != nil {
			log.Fatalf("PANIC: %+v", e)
		}
	}()

	r.Run()
	return nil
}

func New(
	logConfig logcfg.Config,
	httpHandler *api.HTTPServer,
	project *project.Service,
	provider config.Provider,
) *cli.App {
	setupLogging(logConfig.Level)

	return &cli.App{
		Name:        "protoflow",
		Description: "Coding as easy as playing with legos.",
		Flags:       []cli.Flag{},
		Commands: []*cli.Command{
			// TODO breadchris how can you provide a command through wire?
			{
				Name: "worker",
				Action: func(ctx *cli.Context) error {
					client, err := temporal.Wire(provider)
					if err != nil {
						return err
					}
					return workflow.NewWorker(client).Run()
				},
			},
			{
				Name: "studio",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "http",
						Usage: "Port for the studio",
					},
					&cli.BoolFlag{
						Name:  "dev",
						Usage: "Start server in dev mode",
					},
				},
				Action: func(ctx *cli.Context) error {
					dev := ctx.Bool("dev")
					if dev {
						return liveReload()
					}

					httpPort := ctx.Int("http")
					if httpPort == 0 {
						httpPort = 8080
					}

					// TODO breadchris for local dev, add live reload into command https://github.com/makiuchi-d/arelo/blob/master/arelo.go

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
