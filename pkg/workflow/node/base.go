package node

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/rs/zerolog/log"
	"strings"
)

type BaseNode struct {
	Name       string
	id         string
	resourceID string
}

// NewBaseNode creates a new BaseNode from a gen.Node, gen.Node cannot be embedded into BaseNode because proto deserialization will fail on the type
func NewBaseNode(node *gen.Node) BaseNode {
	return BaseNode{
		Name:       util.ToTitleCase(node.Name),
		id:         node.Id,
		resourceID: node.ResourceId,
	}
}

func (n *BaseNode) NormalizedName() string {
	name := util.ToTitleCase(n.Name)
	if strings.Contains(name, ".") {
		name = strings.Split(name, ".")[1]
	}
	return name
}

func (n *BaseNode) ID() string {
	return n.id
}

func (n *BaseNode) ResourceID() string {
	return n.resourceID
}

func (n *BaseNode) Info(r graph.Resource) (*graph.Info, error) {
	log.Warn().
		Str("node", n.Name).
		Msg("Info() not implemented")
	return nil, nil
}

func (n *BaseNode) Represent() (string, error) {
	return "", errors.New("not implemented")
}
