package config

import (
	"github.com/breadchris/protoflow/pkg/log"
	"github.com/google/wire"
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
)

var ProviderSet = wire.NewSet(
	cache.NewLocalCache,
	wire.Bind(new(cache.Cache), new(*cache.LocalCache)),

	log.NewConfig,

	NewProvider,
)
