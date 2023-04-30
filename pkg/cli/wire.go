//go:build wireinject
// +build wireinject

package cli

import (
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/pkg/api"
	"github.com/protoflow-labs/protoflow/pkg/cache"
	"github.com/protoflow-labs/protoflow/pkg/config"
	"github.com/protoflow-labs/protoflow/pkg/generate"
	"github.com/protoflow-labs/protoflow/pkg/k8s"
	"github.com/protoflow-labs/protoflow/pkg/project"
	urfavcli "github.com/urfave/cli/v2"
)

func Wire(cacheConfig cache.Config) (*urfavcli.App, error) {
	panic(wire.Build(
		New,
		config.ProviderSet,
		k8s.ProviderSet,
		project.ProviderSet,
		generate.ProviderSet,
		api.NewHTTPServer,
	))
}
