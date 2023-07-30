package graph

import (
	"github.com/protoflow-labs/protoflow/gen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// TODO breadchris figure out how to make this generic
func ConvertProto(p *gen.Project) *protoProject[*gen.Node] {
	pp := &protoProject[*gen.Node]{
		id: p.GetId(),
	}
	for _, n := range p.GetGraph().GetNodes() {
		pp.graph.nodes = append(pp.graph.nodes, n)
	}
	for _, e := range p.GetGraph().GetEdges() {
		pp.graph.edges = append(pp.graph.edges, e)
	}
	return pp
}

type ProtoNode interface {
	protoreflect.ProtoMessage
	GetId() string
}

type ProtoEdge interface {
	protoreflect.ProtoMessage
	GetFrom() string
	GetTo() string
}

type ProtoGraph[T ProtoNode] interface {
	SetNodes([]T)
	SetEdges([]*gen.Edge)
	GetNodes() []T
	GetEdges() []*gen.Edge
}

type protoGraph[T ProtoNode] struct {
	nodes []T
	edges []*gen.Edge
}

func (g *protoGraph[T]) GetNodes() []T {
	return g.nodes
}

func (g *protoGraph[T]) GetEdges() []*gen.Edge {
	return g.edges
}

func (g *protoGraph[T]) SetNodes(ts []T) {
	g.nodes = ts
}

func (g *protoGraph[T]) SetEdges(edges []*gen.Edge) {
	g.edges = edges
}

type ProtoProject[T ProtoNode] interface {
	GetId() string
	GetGraph() ProtoGraph[T]
}

type protoProject[T ProtoNode] struct {
	id    string
	graph protoGraph[T]
}

func (p *protoProject[T]) SetId(id string) {
	p.id = id
}

func (p *protoProject[T]) GetId() string {
	return p.id
}

func (p *protoProject[T]) GetGraph() ProtoGraph[T] {
	return &p.graph
}
