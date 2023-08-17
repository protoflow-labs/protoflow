package reason

import (
	"github.com/google/uuid"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/reason"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	"github.com/rs/zerolog/log"
)

func New(b *base.Node, node *reason.Reason) graph.Node {
	switch t := node.Type.(type) {
	case *reason.Reason_Prompt:
		return NewPromptNode(b, t.Prompt)
	case *reason.Reason_Engine:
		return NewEngineNode(b, t.Engine)
	default:
		log.Error().Msgf("unknown reason type %T", t)
		return nil
	}
}

func NewProto(name string, d *reason.Reason) *gen.Node {
	return &gen.Node{
		Id:   uuid.NewString(),
		Name: name,
		Type: &gen.Node_Reason{
			Reason: d,
		},
	}
}
