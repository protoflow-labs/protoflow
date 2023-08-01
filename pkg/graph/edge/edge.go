package edge

import (
	"github.com/google/uuid"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/graph"
)

func New(edge *gen.Edge) graph.Edge {
	switch t := edge.Type.(type) {
	case *gen.Edge_Provides:
		return NewProvides(edge, t.Provides)
	case *gen.Edge_PublishesTo:
		return NewPublishesTo(edge, t.PublishesTo)
	default:
		return nil
	}
}

// TODO breadchris find a better package for this
func NewEdgeProto(from, to string) *gen.Edge {
	return &gen.Edge{
		Id:   uuid.NewString(),
		From: from,
		To:   to,
	}
}

func NewProvidesProto(from, to string) *gen.Edge {
	e := NewEdgeProto(from, to)
	e.Type = &gen.Edge_Provides{}
	return e
}

func NewPublishesToProto(from, to string) *gen.Edge {
	e := NewEdgeProto(from, to)
	e.Type = &gen.Edge_PublishesTo{}
	return e
}
