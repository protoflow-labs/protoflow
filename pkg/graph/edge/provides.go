package edge

import (
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/graph"
)

type Provides struct {
	*Base
	*gen.Provides
}

func NewProvides(edge *gen.Edge, p *gen.Provides) graph.Edge {
	return &Provides{
		Base:     NewBase(edge),
		Provides: p,
	}
}

func (p *Provides) Connect(from, to graph.Node) error {
	return to.SetProvider(from)
}
