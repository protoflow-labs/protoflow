package edge

import (
	"github.com/protoflow-labs/protoflow/gen"
)

type Base struct {
	Edge *gen.Edge
}

func (b *Base) From() string {
	return b.Edge.From
}

func (b *Base) To() string {
	return b.Edge.To
}

func (b *Base) ID() string {
	return b.Edge.Id
}

func (b *Base) CanWire() bool {
	return true
}

func NewBase(edge *gen.Edge) *Base {
	return &Base{
		Edge: edge,
	}
}
