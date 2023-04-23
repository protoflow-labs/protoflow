//go:build wireinject
// +build wireinject

package cli

import (
	"github.com/google/wire"
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
	"github.com/protoflow-labs/protoflow-editor/protoflow/pkg/api"
	"github.com/protoflow-labs/protoflow-editor/protoflow/pkg/config"
	"github.com/protoflow-labs/protoflow-editor/protoflow/pkg/workflow"
	urfavcli "github.com/urfave/cli/v2"
)

func Wire(cacheConfig cache.Config) (*urfavcli.App, error) {
	panic(wire.Build(
		New,
		config.ProviderSet,
		workflow.ProviderSet,
		api.NewHTTPServer,
		workflow.NewWorker,
	))
}
