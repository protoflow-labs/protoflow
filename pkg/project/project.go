package project

import (
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
)

type Project struct {
	Base     *gen.Project
	Workflow *workflow.Workflow
}

func FromProto(project *gen.Project) (*Project, error) {
	w, err := workflow.FromProject(project)
	if err != nil {
		return nil, err
	}
	return &Project{
		Base:     project,
		Workflow: w,
	}, nil
}
