package node

import (
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
)

type ConfigNode struct {
	BaseNode
	Config *gen.Config
}

var _ graph.Node = &ConfigNode{}

func NewConfigNode(node *gen.Node) *ConfigNode {
	return &ConfigNode{
		BaseNode: NewBaseNode(node),
		Config:   node.GetConfiguration(),
	}
}

func (c *ConfigNode) Represent() (string, error) {
	return c.Config.Value, nil
}
