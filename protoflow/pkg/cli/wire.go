//go:build wireinject
// +build wireinject

package cli

import (
	"github.com/breadchris/protoflow/pkg/api"
	"github.com/breadchris/protoflow/pkg/config"
	"github.com/breadchris/protoflow/pkg/workflow"
	"github.com/google/wire"
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
	urfavcli "github.com/urfave/cli/v2"
)

func Wire(cacheConfig cache.Config) (*urfavcli.App, error) {
	panic(wire.Build(
		New,
		config.ProviderSet,
		workflow.ProviderSet,
		api.NewHTTPServer,
		api.NewGRPCServer,
	))
}
