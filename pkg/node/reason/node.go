package reason

import (
	"github.com/protoflow-labs/protoflow/gen/reason"
	"github.com/protoflow-labs/protoflow/pkg/node/base"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
)

func New(b *base.Node, node *reason.Reason) graph.Node {
	switch t := node.Type.(type) {
	case *reason.Reason_Prompt:
		return NewPromptNode(b, t.Prompt)
	case *reason.Reason_Engine:
		return NewEngineNode(b, t.Engine)
	default:
		return nil
	}
}
