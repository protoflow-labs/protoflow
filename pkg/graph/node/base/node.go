package base

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"github.com/rs/zerolog/log"
	"strings"
)

type Node struct {
	Name string
	id   string

	provider    graph.Node
	dependents  []graph.Node
	subscribers []graph.Listener
	publishers  []graph.Listener
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

func (n *Node) Type() (*graph.Info, error) {
	log.Warn().
		Str("node", n.Name).
		Msg("Type() not implemented")
	// TODO breadchris this doesn't make node generic, figure out a way to make this generic
	return graph.NewInfoFromType("node", &gen.Node{})
}

func (n *Node) Method() (*graph.Info, error) {
	log.Warn().
		Str("node", n.Name).
		Msg("Method() not implemented")
	return nil, errors.New("not implemented")
}

func (n *Node) Init() (func(), error) {
	return func() {}, nil
}

// TODO breadchris this should be more robust and take into consideration type types of edges into the node
func (n *Node) Provider() (graph.Node, error) {
	if n.provider == nil {
		return nil, errors.New("no provider")
	}
	return n.provider, nil
}

func (n *Node) SetProvider(p graph.Node) error {
	if n.provider != nil {
		return errors.New("provider already set")
	}
	n.provider = p
	return nil
}

func (n *Node) Dependents() []graph.Node {
	return n.dependents
}

func (n *Node) AddDependent(p graph.Node) {
	n.dependents = append(n.dependents, p)
}

func (n *Node) Subscribers() []graph.Listener {
	return n.subscribers
}

func (n *Node) AddSubscriber(p graph.Listener) {
	n.subscribers = append(n.subscribers, p)
}

func (n *Node) Publishers() []graph.Listener {
	return n.publishers
}

func (n *Node) AddPublishers(p graph.Listener) {
	n.publishers = append(n.publishers, p)
}
