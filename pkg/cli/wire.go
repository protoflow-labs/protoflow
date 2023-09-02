//go:build wireinject
// +build wireinject

package cli

import (
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	urfavcli "github.com/urfave/cli/v2"
)

func Wire(cacheConfig bucket.Config) (*urfavcli.App, error) {
	panic(wire.Build(
		ProviderSet,
	))
}
