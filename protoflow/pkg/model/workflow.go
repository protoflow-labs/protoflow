package model

import (
	"github.com/breadchris/protoflow/gen/workflow"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/datatypes"
)

type WorkflowDoc struct {
	workflow.Workflow
}

func (w *WorkflowDoc) UnmarshalJSON(data []byte) error {
	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	if err := unmarshaler.Unmarshal(rawRespBody, out); err != nil {
		return err
	}
}

type Workflow struct {
	Base

	Protoflow datatypes.JSONType[workflow.Workflow]
}

type WorkflowRun struct {
	Base

	Workflow Workflow
}
