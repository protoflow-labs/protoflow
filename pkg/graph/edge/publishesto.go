package edge

import (
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/graph"
)

type PublishesTo struct {
	*Base
	*gen.PublishesTo
}

func NewPublishesTo(edge *gen.Edge, p *gen.PublishesTo) graph.Edge {
	return &PublishesTo{
		Base:        NewBase(edge),
		PublishesTo: p,
	}
}

func (p *PublishesTo) Connect(from, to graph.Node) error {
	from.AddSubscriber(NewPublishesToListener(to, p))
	to.AddPublishers(from)
	return nil
}
