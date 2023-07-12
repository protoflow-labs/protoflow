package resource

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"go.uber.org/config"
	"gopkg.in/yaml.v3"
)

type ConfigProviderResource struct {
	*BaseResource
	*gen.ConfigProvider
}

var _ Resource = &ConfigProviderResource{}

func (r *ConfigProviderResource) Init() (func(), error) {
	return nil, nil
}

func (r *ConfigProviderResource) Build(options ...config.YAMLOption) (config.Provider, error) {
	opts := []config.YAMLOption{
		config.Permissive(),
	}

	for _, o := range options {
		opts = append(opts, o)
	}

	for _, n := range r.nodes {
		c, ok := n.(*node.ConfigNode)
		if !ok {
			return nil, errors.New("node is not a config provider node")
		}

		var u map[string]interface{}
		err := yaml.Unmarshal([]byte(c.Config.Value), &u)
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
