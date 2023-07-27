package workflow

import (
	"context"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/store"
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

func (m *MemoryManager) saveNodeExecutions(projectID, nodeID string, trace rxgo.Observable) error {
	workflowRun := &gen.WorkflowRun{
		Id: uuid.NewString(),
		Request: &gen.RunWorkflowRequest{
			ProjectId: projectID,
			NodeId:    nodeID,
		},
	}
	workflowRunID, err := m.store.CreateWorkflowRun(workflowRun)
	if err != nil {
		return errors.Wrap(err, "error creating workflow run")
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
	return nil
}

func (m *MemoryManager) ExecuteWorkflow(ctx context.Context, w *Workflow, nodeID string, input rxgo.Observable) (rxgo.Observable, error) {
	return w.WireNodes(ctx, nodeID, input)

	//logger := &MemoryLogger{}
	//
	//memoryCtx := &execute.MemoryContext{Context: ctx}
	//executor := execute.NewMemoryExecutor(memoryCtx)

	//if err != nil {
	//	return nil, errors.Wrapf(err, "error running workflow")
	//}
	//err = m.saveNodeExecutions(w.ProjectID, nodeID, obs)
	//if err != nil {
	//	return nil, errors.Wrapf(err, "error saving node executions")
	//}
	//return obs, nil
}

func (m *MemoryManager) CleanupResources() error {
	//TODO implement me
	panic("implement me")
}
