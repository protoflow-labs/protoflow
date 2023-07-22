package resource

import (
	"github.com/protoflow-labs/protoflow/gen"
)

type TemplateServiceResource struct {
	*BaseResource
	*gen.TemplateService
}

func (r *TemplateServiceResource) Init() (func(), error) {
	return nil, nil
}
