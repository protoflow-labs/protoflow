//go:build wireinject
// +build wireinject

package cli

import (
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/pkg/api"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/config"
	"github.com/protoflow-labs/protoflow/pkg/generate"
	"github.com/protoflow-labs/protoflow/pkg/project"
	urfavcli "github.com/urfave/cli/v2"
)

func Wire(cacheConfig bucket.Config) (*urfavcli.App, error) {
	panic(wire.Build(
		New,
		config.ProviderSet,
		project.ProviderSet,
		generate.ProviderSet,
		api.ProviderSet,
	))
}
