package cli

import (
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
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
	p *project.Service,
	l *logd.Log,
) *cli.App {
	return &cli.App{
		Name:        "protoflow",
		Description: "Coding as easy as playing with legos.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "project",
				Usage: "Project directory. Defaults to current working directory.",
			},
		},
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
				Name: "generate",
				Subcommands: []*cli.Command{
					{
						Name: "init",
						Action: func(ctx *cli.Context) error {
							pDir := ctx.String("project")

							// TODO breadchris not friendly
							b, err := bucket.NewUserCache(bucket.Config{
								Name: pDir,
							})
							if err != nil {
								return err
							}

							proj, err := project.NewDefaultProject(b)
							if err != nil {
								return err
							}

							w, err := project.FromProto(proj)
							if err != nil {
								return err
							}

							g, err := generate.NewGenerate(generate.Config{
								ProjectPath: pDir,
							})
							if err != nil {
								return err
							}

							err = g.Init(&project.Project{
								Base:     proj,
								Workflow: w,
							})
							if err != nil {
								return err
							}
							return nil
						},
					},
					{
						Name: "service",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name: "name",
							},
							&cli.StringFlag{
								Name:    "language",
								Aliases: []string{"l"},
								Usage:   "Language to generate",
							},
						},
						Action: func(ctx *cli.Context) error {
							return nil
						},
					},
					{
						Name: "method",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name: "name",
							},
							&cli.StringFlag{
								Name:    "language",
								Aliases: []string{"l"},
								Usage:   "Language to generate",
							},
						},
						Action: func(ctx *cli.Context) error {
							return nil
						},
					},
				},
				Action: func(ctx *cli.Context) error {
					return nil
				},
			},
			{
				Name: "load",
				Action: func(ctx *cli.Context) error {
					res, err := p.LoadProject(ctx.Context, connect.NewRequest(&gen.LoadProjectRequest{}))
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
