package graph

import (
	"github.com/protoflow-labs/protoflow/gen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// TODO breadchris figure out how to make this generic
func ConvertProto(p *gen.Project) *protoProject[*gen.Node, *gen.Edge] {
	pp := &protoProject[*gen.Node, *gen.Edge]{
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
	GetId() string
	GetFrom() string
	GetTo() string
}

type ProtoGraph[T ProtoNode, U ProtoEdge] interface {
	SetNodes([]T)
	SetEdges([]U)
	GetNodes() []T
	GetEdges() []U
}

type protoGraph[T ProtoNode, U ProtoEdge] struct {
	nodes []T
	edges []U
}

func (g *protoGraph[T, U]) GetNodes() []T {
	return g.nodes
}

func (g *protoGraph[T, U]) GetEdges() []U {
	return g.edges
}

func (g *protoGraph[T, U]) SetNodes(ts []T) {
	g.nodes = ts
}

func (g *protoGraph[T, U]) SetEdges(edges []U) {
	g.edges = edges
}

type ProtoProject[T ProtoNode, U ProtoEdge] interface {
	GetId() string
	GetGraph() ProtoGraph[T, U]
}

type protoProject[T ProtoNode, U ProtoEdge] struct {
	id    string
	graph protoGraph[T, U]
}

func (p *protoProject[T, U]) SetId(id string) {
	p.id = id
}

func (p *protoProject[T, U]) GetId() string {
	return p.id
}

func (p *protoProject[T, U]) GetGraph() ProtoGraph[T, U] {
	return &p.graph
}
