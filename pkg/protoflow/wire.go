//go:build wireinject
// +build wireinject

package protoflow

import (
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/config"
	"github.com/protoflow-labs/protoflow/pkg/generate"
	"github.com/protoflow-labs/protoflow/pkg/project"
	"github.com/protoflow-labs/protoflow/pkg/server"
)

// TODO breadchris should not need a bucket config, should be able to pass config in-memory
func Wire(cacheConfig bucket.Config, defaultProject *gen.Project) (*Protoflow, error) {
	panic(wire.Build(
		config.ProviderSet,
		project.ProviderSet,
		generate.ProviderSet,
		server.ProviderSet,
		New,
	))
}
