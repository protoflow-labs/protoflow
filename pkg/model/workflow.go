package model

type WorkflowRun struct {
	UUID
	Times

	ProjectID string
	*WorkflowRunJSON

	NodeExecutions []NodeExecution
}

type NodeExecution struct {
	UUID
	Times

	WorkflowRunID string

	*NodeExecutionJSON
}
