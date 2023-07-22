package execute

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"go.temporal.io/sdk/workflow"
)

type Executor interface {
	Execute(n graph.Node, input graph.Input) (*graph.Output, error)
}

var _ Executor = &TemporalExecutor{}

type TemporalExecutor struct {
	ctx workflow.Context
}

func NewTemporalExecutor(ctx workflow.Context) *TemporalExecutor {
	return &TemporalExecutor{
		ctx: ctx,
	}
}

func (e *TemporalExecutor) Execute(n graph.Node, input graph.Input) (*graph.Output, error) {
	var result graph.Output
	// TODO breadchris n.Wire will not work here since n is a pointer and we are changing execution context, i think
	err := workflow.ExecuteActivity(e.ctx, n.Wire, n, input).Get(e.ctx, &result)
	if err != nil {
		return nil, errors.Wrap(err, "error executing activity")
	}
	return &result, nil
}

type MemoryExecutor struct {
	ctx *MemoryContext
}

var _ Executor = &MemoryExecutor{}

type MemoryExecutorOption func(*MemoryExecutor)

func NewMemoryExecutor(ctx *MemoryContext, opts ...MemoryExecutorOption) *MemoryExecutor {
	e := &MemoryExecutor{
		ctx: ctx,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *MemoryExecutor) Execute(n graph.Node, input graph.Input) (*graph.Output, error) {
	res, err := n.Wire(e.ctx.Context, input)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
