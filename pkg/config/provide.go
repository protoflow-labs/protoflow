package config

import (
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
)

var ProviderSet = wire.NewSet(
	bucket.NewUserCache,
	wire.Bind(new(bucket.Bucket), new(*bucket.LocalBucket)),
	NewProvider,
)
