package resource

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
)

type TemplateServiceResource struct {
	*BaseResource
	*gen.TemplateService
}

func (r *TemplateServiceResource) Init() (func(), error) {
	return nil, nil
}

func (r *TemplateServiceResource) Info(n node.Node) (*node.Info, error) {
	_, ok := n.(*node.TemplateNode)
	if !ok {
		return nil, errors.New("node is not a template node")
	}
	return nil, nil
}
