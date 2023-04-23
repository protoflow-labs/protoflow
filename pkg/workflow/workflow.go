package workflow

import (
	"fmt"
	"time"

	"github.com/protoflow-labs/protoflow/gen"

	"github.com/hmdsefi/gograph"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"go.temporal.io/sdk/workflow"
)

type Result struct {
	Data    string
	IsError bool
	Error   string
}

type Workflow struct {
	Graph       gograph.Graph[string]
	BlockLookup map[string]Block
}

type WorkflowGraph struct {
	Nodes       []*gen.Node
	Edges       []*gen.Edge
	BlockLookup map[string]Block
}

func NewWorkflowFromProtoflow(workflowGraph *gen.Workflow) (*Workflow, error) {
	graph := gograph.New[string](gograph.Directed())

	blockLookup := map[string]Block{}
	vertexLookup := map[string]*gograph.Vertex[string]{}
	for _, node := range workflowGraph.Nodes {
		v := gograph.NewVertex(node.Id)
		graph.AddVertex(v)

		// add vertex to lookup to be used for edges
		vertexLookup[node.Id] = v

		// add block to lookup to be used for execution
		activity, err := NewBlock(node)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating block for node %s", node.Id)
		}
		blockLookup[node.Id] = activity
	}

	for _, edge := range workflowGraph.Edges {
		_, err := graph.AddEdge(vertexLookup[edge.From], vertexLookup[edge.To])
		if err != nil {
			return nil, err
		}
	}

	return &Workflow{
		Graph:       graph,
		BlockLookup: blockLookup,
	}, nil
}

func Run(ctx workflow.Context, w *Workflow, nodeID string) (string, error) {
	if w.BlockLookup == nil || w.Graph == nil {
		return "", fmt.Errorf("workflow is not initialized")
	}

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	// Adding context to a workflow
	// ctx = workflow.WithValue(ctx, AccountIDContextKey, dslWorkflow.AccountID)

	logger.Info("Starting workflow", "workflowID", workflow.GetInfo(ctx).WorkflowExecution.ID, "nodeID", nodeID)

	vert := w.Graph.GetVertexByID(nodeID)
	_, err := w.traverseWorkflow(ctx, vert)
	if err != nil {
		logger.Error("Error traversing workflow", "error", err)
		return "", nil
	}
	return "", nil
}

func (w *Workflow) traverseWorkflow(ctx workflow.Context, vert *gograph.Vertex[string]) (*Result, error) {
	logger := workflow.GetLogger(ctx)
	if vert == nil {
		return nil, errors.New("vertex is nil")
	}

	for _, neighbor := range vert.Neighbors() {
		block, ok := w.BlockLookup[neighbor.Label()]
		if !ok {
			return nil, fmt.Errorf("vertex not found: %s", neighbor.Label())
		}

		res, err := block.Execute(ctx)
		if err != nil {
			logger.Error("Error executing block", "error", err)
			return nil, errors.Wrapf(err, "error executing block %s", neighbor.Label())
		}

		log.Debug().Interface("result", res).Msg("block result")

		logger.Info("Traversing workflow", "nodeID", neighbor.Label())
		_, err = w.traverseWorkflow(ctx, neighbor)
		if err != nil {
			return nil, errors.Wrapf(err, "error traversing workflow %s", neighbor.Label())
		}
	}
	return &Result{}, nil
}
