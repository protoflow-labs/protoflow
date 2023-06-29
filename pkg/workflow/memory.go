package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/store"
	"github.com/protoflow-labs/protoflow/pkg/workflow/execute"
	"github.com/rs/zerolog/log"
	"sync"

	"github.com/google/uuid"
)

type MemoryManager struct {
	resourceCleanup sync.Cond
	store           store.Project
}

var _ Manager = (*MemoryManager)(nil)

func NewMemoryManager(store store.Project) *MemoryManager {
	return &MemoryManager{
		store: store,
	}
}

func (m *MemoryManager) saveNodeExecutions(projectID, nodeID string, input interface{}, trace chan *gen.NodeExecution) {
	serInp, err := json.Marshal(input)
	if err != nil {
		log.Error().Err(err).Msg("error serializing input")
		return
	}
	workflowRun := &gen.WorkflowRun{
		Id: uuid.NewString(),
		Request: &gen.RunWorkflowRequest{
			ProjectId: projectID,
			NodeId:    nodeID,
			Input:     string(serInp),
		},
	}
	workflowRunID, err := m.store.CreateWorkflowRun(workflowRun)
	if err != nil {
		log.Error().Err(err).Msg("error creating workflow run")
		return
	}
	for nodeExecution := range trace {
		_, err := m.store.SaveNodeExecution(workflowRunID, nodeExecution)
		if err != nil {
			log.Error().Err(err).Msg("error saving node execution")
		}
	}
}

func (m *MemoryManager) ExecuteWorkflow(ctx context.Context, w *Workflow, nodeID string, input interface{}) (string, error) {
	if w.NodeLookup == nil || w.Graph == nil {
		return "", fmt.Errorf("workflow is not initialized")
	}

	logger := &MemoryLogger{}

	trace := make(chan *gen.NodeExecution)
	go m.saveNodeExecutions(w.ProjectID, nodeID, input, trace)

	memoryCtx := &execute.MemoryContext{Context: ctx}
	executor := execute.NewMemoryExecutor(memoryCtx, execute.WithTrace(trace))

	_, err := w.Run(logger, executor, nodeID, input)
	return uuid.New().String(), err
}

func (m *MemoryManager) ExecuteWorkflowSync(ctx context.Context, w *Workflow, nodeID string, input interface{}) ([]any, error) {
	if w.NodeLookup == nil || w.Graph == nil {
		return nil, fmt.Errorf("workflow is not initialized")
	}

	logger := &MemoryLogger{}

	trace := make(chan *gen.NodeExecution)
	go m.saveNodeExecutions(w.ProjectID, nodeID, input, trace)

	memoryCtx := &execute.MemoryContext{Context: ctx}
	executor := execute.NewMemoryExecutor(memoryCtx, execute.WithTrace(trace))

	return w.Run(logger, executor, nodeID, input)
}

func (m *MemoryManager) CleanupResources() error {
	//TODO implement me
	panic("implement me")
}
