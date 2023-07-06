package resource

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
)

type ReasoningEngineResource struct {
	*BaseResource
	*gen.ReasoningEngine
}

func (r *ReasoningEngineResource) Init() (func(), error) {
	return nil, nil
}

func (r *ReasoningEngineResource) Info(n node.Node) (*node.Info, error) {
	_, ok := n.(*node.PromptNode)
	if !ok {
		return nil, errors.New("node is not a prompt node")
	}
	return nil, nil
}
