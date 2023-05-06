package config

import (
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/pkg/cache"
	"github.com/protoflow-labs/protoflow/pkg/log"
)

var ProviderSet = wire.NewSet(
	cache.NewUserCache,
	wire.Bind(new(cache.Cache), new(*cache.LocalCache)),

	log.NewConfig,

	NewProvider,
)
