package resource

import (
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
)

type ConfigProviderResource struct {
	*BaseResource
	*gen.ConfigProvider
}

var _ graph.Resource = &ConfigProviderResource{}

func (r *ConfigProviderResource) Init() (func(), error) {
	return nil, nil
}
