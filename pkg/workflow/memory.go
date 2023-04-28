package workflow

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type MemoryManager struct {
	resourceCleanup sync.Cond
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{}
}

func (m *MemoryManager) ExecuteWorkflow(ctx context.Context, w *Workflow, nodeID string) (string, error) {
	if w.NodeLookup == nil || w.Graph == nil {
		return "", fmt.Errorf("workflow is not initialized")
	}

	logger := &MemoryLogger{}

	memoryCtx := &MemoryContext{Context: ctx}
	executor := NewMemoryExecutor(memoryCtx)

	_, err := w.Run(logger, executor, nodeID)
	return uuid.New().String(), err
}

func (m *MemoryManager) ExecuteWorkflowSync(ctx context.Context, w *Workflow, nodeID string) (*Result, error) {
	if w.NodeLookup == nil || w.Graph == nil {
		return nil, fmt.Errorf("workflow is not initialized")
	}

	logger := &MemoryLogger{}

	memoryCtx := &MemoryContext{Context: ctx}
	executor := NewMemoryExecutor(memoryCtx)

	return w.Run(logger, executor, nodeID)
}

func (m *MemoryManager) CleanupResources() error {
	//TODO implement me
	panic("implement me")
}
