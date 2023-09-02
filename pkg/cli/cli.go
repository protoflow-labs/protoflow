package cli

import (
	"github.com/bufbuild/connect-go"
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/config"
	"github.com/protoflow-labs/protoflow/pkg/generate"
	"github.com/protoflow-labs/protoflow/pkg/util/reload"

	"github.com/protoflow-labs/protoflow/pkg/api"
	logd "github.com/protoflow-labs/protoflow/pkg/log"
	"github.com/protoflow-labs/protoflow/pkg/project"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var ProviderSet = wire.NewSet(
	New,
	project.NewDefaultProject,
	logd.ProviderSet,
	config.ProviderSet,
	project.ProviderSet,
	generate.ProviderSet,
	api.ProviderSet,
)

func New(
	httpHandler *api.HTTPServer,
	project *project.Service,
	l *logd.Log,
) *cli.App {
	return &cli.App{
		Name:        "protoflow",
		Description: "Coding as easy as playing with legos.",
		Flags:       []cli.Flag{},
		Commands: []*cli.Command{
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

					// TODO breadchris for local dev, add live reload into command https://github.com/makiuchi-d/RELOAD/blob/master/RELOAD.go

					log.Info().Int("port", httpPort).Msg("starting http server")
					return httpHandler.Serve(httpPort)
				},
			},
			{
				Name: "load",
				Action: func(ctx *cli.Context) error {
					res, err := project.LoadProject(ctx.Context, connect.NewRequest(&gen.LoadProjectRequest{}))
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

func liveReload() error {
	// TODO breadchris makes this a config that can be set
	c := reload.Config{
		Cmd: []string{"go", "run", "main.go", "studio"},
		// TODO breadchris the patterns and ignores are not quite working
		// ideally we use tilt here
		Patterns: []string{"pkg/**/*.go", "templates/**"},
		Ignores:  []string{"studio/**", "node_modules/**", ".git/**", "examples/**"},
	}
	// TODO breadchris this code needs to be refactored to use observability
	return reload.Reload(c)
}
