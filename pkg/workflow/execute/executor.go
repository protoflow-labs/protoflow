package execute

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/reactivex/rxgo/v2"
	"go.temporal.io/sdk/workflow"
)

type Input struct {
	Observable rxgo.Observable
	Resource   resource.Resource
}

type Output struct {
	Observable rxgo.Observable
}

type Executor interface {
	Execute(n node.Node, input Input) (*Output, error)
	Trace(nodeExecution *gen.NodeExecution) error
}

var _ Executor = &TemporalExecutor{}

var activity = &Activity{}

// TODO breadchris each node should have an Execute function, however at the moment, an import cycle would be formed by
// node -> resource -> node. Need to figure out how to avoid this.
func nodeToActivityName(n node.Node) ActivityFunc {
	switch n.(type) {
	case *node.RESTNode:
		return activity.ExecuteRestNode
	case *node.GRPCNode:
		return activity.ExecuteGRPCNode
	case *node.FunctionNode:
		return activity.ExecuteFunctionNode
	case *node.CollectionNode:
		return activity.ExecuteCollectionNode
	case *node.BucketNode:
		return activity.ExecuteBucketNode
	case *node.InputNode:
		return activity.ExecuteInputNode
	case *node.QueryNode:
		return activity.ExecuteQueryNode
	}
	return nil
}

type TemporalExecutor struct {
	ctx workflow.Context
}

func NewTemporalExecutor(ctx workflow.Context) *TemporalExecutor {
	return &TemporalExecutor{
		ctx: ctx,
	}
}

func (e *TemporalExecutor) Execute(n node.Node, input Input) (*Output, error) {
	var result Output
	act := nodeToActivityName(n)
	if act == nil {
		return nil, fmt.Errorf("error getting activity for node: %s", n.NormalizedName())
	}
	err := workflow.ExecuteActivity(e.ctx, act, n, input).Get(e.ctx, &result)
	if err != nil {
		return nil, errors.Wrap(err, "error executing activity")
	}
	return &result, nil
}

func (e *TemporalExecutor) Trace(nodeExecution *gen.NodeExecution) error {
	//TODO implement me
	panic("implement me")
}

type MemoryExecutor struct {
	ctx   *MemoryContext
	trace chan<- *gen.NodeExecution
}

var _ Executor = &MemoryExecutor{}

type MemoryExecutorOption func(*MemoryExecutor)

func WithTrace(trace chan<- *gen.NodeExecution) MemoryExecutorOption {
	return func(executor *MemoryExecutor) {
		executor.trace = trace
	}
}

func NewMemoryExecutor(ctx *MemoryContext, opts ...MemoryExecutorOption) *MemoryExecutor {
	e := &MemoryExecutor{
		ctx: ctx,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *MemoryExecutor) Execute(n node.Node, input Input) (*Output, error) {
	act := nodeToActivityName(n)
	if act == nil {
		return nil, fmt.Errorf("error getting activity for node: %s", n.NormalizedName())
	}

	res, err := act(e.ctx.Context, n, input)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (e *MemoryExecutor) Trace(nodeExecution *gen.NodeExecution) error {
	if e.trace != nil {
		// Do not block on writing to this channel
		go func() {
			e.trace <- nodeExecution
		}()
	}
	return nil
}
