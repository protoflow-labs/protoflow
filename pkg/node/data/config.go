package data

import (
	"context"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen/data"
	"github.com/protoflow-labs/protoflow/pkg/node/base"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"go.uber.org/config"
	"gopkg.in/yaml.v3"
)

type ConfigNode struct {
	*base.Node
	*data.Config
}

var _ graph.Node = &ConfigNode{}

func NewConfigNode(b *base.Node, node *data.Config) *ConfigNode {
	return &ConfigNode{
		Node:   b,
		Config: node,
	}
}

func (c *ConfigNode) NewConfigProvider(options ...config.YAMLOption) (config.Provider, error) {
	opts := []config.YAMLOption{
		config.Permissive(),
	}

	for _, o := range options {
		opts = append(opts, o)
	}

	var u map[string]any
	err := yaml.Unmarshal([]byte(c.Config.Value), &u)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal config yaml")
	}
	opts = append(opts, config.Static(u))

	conf, err := config.NewYAML(opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create config provider")
	}
	return conf, nil
}

func (c *ConfigNode) Represent() (string, error) {
	return c.Config.Value, nil
}

func (c *ConfigNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	return graph.Output{}, errors.New("not implemented")
}
