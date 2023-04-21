package model

import (
	"github.com/breadchris/protoflow/gen/workflow"
	"gorm.io/datatypes"
)

type Workflow struct {
	Base

	Protoflow datatypes.JSONType[workflow.Workflow]
}

type WorkflowRun struct {
	Base

	Workflow Workflow
}
