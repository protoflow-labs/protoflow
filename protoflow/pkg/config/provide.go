package config

import (
	"github.com/google/wire"
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
	"github.com/protoflow-labs/protoflow-editor/protoflow/pkg/log"
)

var ProviderSet = wire.NewSet(
	cache.NewLocalCache,
	wire.Bind(new(cache.Cache), new(*cache.LocalCache)),

	log.NewConfig,

	NewProvider,
)
