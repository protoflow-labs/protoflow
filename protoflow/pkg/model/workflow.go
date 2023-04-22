package model

type Workflow struct {
	Base

	Protoflow ProtoJSON[interface{}]
}

type WorkflowRun struct {
	Base

	Workflow Workflow
}
