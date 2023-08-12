package data

import (
	"context"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen/data"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
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

func NewConfigProto(value any) *data.Data {
	b, err := yaml.Marshal(value)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal config value")
		return nil
	}
	return &data.Data{
		Type: &data.Data_Config{
			Config: &data.Config{
				Value: string(b),
			},
		},
	}
}

//func (c *ConfigNode) Type() (*graph.Type, error) {
//	return graph.NewInfoFromType("config", &data.Config{})
//}

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

func (c *ConfigNode) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	var u map[string]any
	err := yaml.Unmarshal([]byte(c.Config.Value), &u)
	if err != nil {
		return graph.IO{}, errors.Wrapf(err, "failed to unmarshal config yaml")
	}

	return graph.IO{
		Observable: rxgo.Just(u)(),
	}, nil
}
