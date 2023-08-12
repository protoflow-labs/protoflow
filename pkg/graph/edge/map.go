package edge

import (
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/graph"
)

type Map struct {
	*Base
	*gen.Map
}

func NewMap(edge *gen.Edge, p *gen.Map) graph.Edge {
	return &Map{
		Base: NewBase(edge),
		Map:  p,
	}
}

func (p *Map) Connect(from, to graph.Node) error {
	from.AddSubscriber(NewMapListener(to, p))
	to.AddPublishers(NewMapListener(from, p))
	return nil
}
