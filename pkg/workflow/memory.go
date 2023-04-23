package workflow

import (
	"context"
	"fmt"
)

type Manager interface {
	ExecuteWorkflow(ctx context.Context, w *Workflow, nodeID string) (string, error)
}

type MemoryManager struct {
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{}
}

func (m *MemoryManager) ExecuteWorkflow(ctx context.Context, w *Workflow, nodeID string) (string, error) {
	if w.BlockLookup == nil || w.Graph == nil {
		return "", fmt.Errorf("workflow is not initialized")
	}

	logger := &MemoryLogger{}

	memoryCtx := &MemoryContext{Context: ctx}
	executor := NewMemoryExecutor(memoryCtx)

	return w.Run(logger, executor, nodeID)
}
