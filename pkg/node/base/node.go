package base

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/rs/zerolog/log"
	"strings"
)

type Node struct {
	Name string
	id   string

	predecessors []graph.Node
	successors   []graph.Node
}

// NodeFromProto creates a new Node from a gen.Node, gen.Node cannot be embedded into Node because proto deserialization will fail on the type
func NodeFromProto(node *gen.Node) *Node {
	return &Node{
		Name: util.ToTitleCase(node.Name),
		id:   node.Id,
	}
}

func NewNode(name string) *Node {
	return &Node{
		Name: name,
		id:   uuid.NewString(),
	}
}

func (n *Node) NormalizedName() string {
	name := util.ToTitleCase(n.Name)
	if strings.Contains(name, ".") {
		name = strings.Split(name, ".")[1]
	}
	return name
}

func (n *Node) ID() string {
	return n.id
}

func (n *Node) Info() (*graph.Info, error) {
	log.Warn().
		Str("node", n.Name).
		Msg("Info() not implemented")
	return nil, nil
}

func (n *Node) Represent() (string, error) {
	return "", errors.New("not implemented")
}

func (n *Node) Init() (func(), error) {
	return func() {}, nil
}

func (n *Node) AddPredecessor(node graph.Node) {
	n.predecessors = append(n.predecessors, node)
}

func (n *Node) AddSuccessor(node graph.Node) {
	n.successors = append(n.successors, node)
}

func (n *Node) Predecessors() []graph.Node {
	return n.predecessors
}

func (n *Node) Successors() []graph.Node {
	return n.successors
}

// TODO breadchris this should be more robust and take into consideration type types of edges into the node
func (n *Node) Provider() (graph.Node, error) {
	if len(n.predecessors) == 0 {
		return nil, errors.New("no provider")
	}
	return n.predecessors[0], nil
}
