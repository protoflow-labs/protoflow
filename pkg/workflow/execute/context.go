package execute

import (
	"context"
	"go.temporal.io/sdk/workflow"
)

type MemoryContext struct {
	context.Context
}

func (m MemoryContext) Done() workflow.Channel {
	panic("implement me")
}

var _ workflow.Context = (*MemoryContext)(nil)
