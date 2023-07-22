package resource

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"go.uber.org/config"
	"gopkg.in/yaml.v3"
)

type ConfigProviderResource struct {
	*BaseResource
	*gen.ConfigProvider
}

var _ graph.Resource = &ConfigProviderResource{}

func (r *ConfigProviderResource) Init() (func(), error) {
	return nil, nil
}

func (r *ConfigProviderResource) NewConfigProvider(options ...config.YAMLOption) (config.Provider, error) {
	opts := []config.YAMLOption{
		config.Permissive(),
	}

	for _, o := range options {
		opts = append(opts, o)
	}

	for _, n := range r.Nodes() {
		repr, err := n.Represent()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to represent node")
		}
		var u map[string]any
		err = yaml.Unmarshal([]byte(repr), &u)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal config yaml")
		}
		opts = append(opts, config.Static(u))
	}

	c, err := config.NewYAML(opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create config provider")
	}
	return c, nil
}
