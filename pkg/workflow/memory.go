package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/store"
	"github.com/protoflow-labs/protoflow/pkg/workflow/execute"
	"github.com/reactivex/rxgo/v2"
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

func (m *MemoryManager) saveNodeExecutions(projectID, nodeID string, input any, trace rxgo.Observable) {
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
	trace.ForEach(func(i any) {
		// TODO breadchris trace should be a generic observable
		if nodeExec, ok := i.(*gen.NodeExecution); ok {
			_, err := m.store.SaveNodeExecution(workflowRunID, nodeExec)
			if err != nil {
				log.Error().Err(err).Msg("error saving node execution")
			}
		} else {
			log.Error().
				Interface("item", i).
				Msg("error saving node execution, not a node execution")
		}
	}, func(err error) {
		log.Error().Err(err).Msg("trace error")
	}, func() {
		log.Debug().Msg("trace complete")
	})
}

func (m *MemoryManager) ExecuteWorkflow(ctx context.Context, w *Workflow, nodeID string, input rxgo.Observable) (rxgo.Observable, error) {
	if w.NodeLookup == nil || w.Graph == nil {
		return nil, fmt.Errorf("workflow is not initialized")
	}

	logger := &MemoryLogger{}

	memoryCtx := &execute.MemoryContext{Context: ctx}
	executor := execute.NewMemoryExecutor(memoryCtx)

	return w.Run(ctx, logger, executor, nodeID, input)
}

func (m *MemoryManager) CleanupResources() error {
	//TODO implement me
	panic("implement me")
}
