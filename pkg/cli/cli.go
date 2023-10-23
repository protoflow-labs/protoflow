package cli

import (
	"fmt"
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
					&cli.BoolFlag{
						Name:  "dev",
						Usage: "Start server in dev mode",
					},
					&cli.IntFlag{
						Name:  "port",
						Usage: "Port to start the server",
					},
				},
				Action: func(ctx *cli.Context) error {
					dev := ctx.Bool("dev")
					port := ctx.Int("port")
					if dev {
						return liveReload(port)
					}
					return httpHandler.Start(port)
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

func liveReload(port int) error {
	// TODO breadchris makes this a config that can be set
	c := reload.Config{
		Cmd:      []string{"go", "run", "main.go", "studio", "--port", fmt.Sprintf("%d", port)},
		Targets:  []string{"pkg", "gen"},
		Patterns: []string{"**/*.go"},
	}
	// TODO breadchris this code needs to be refactored to use observability
	return reload.Reload(c)
}
