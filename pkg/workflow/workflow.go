package workflow

import (
	"fmt"

	"github.com/dominikbraun/graph"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/rs/zerolog/log"
)

type Input struct {
	Params interface{}
}

type Result struct {
	Data interface{}
}

type AdjMap map[string]map[string]graph.Edge[string]

type Workflow struct {
	ID          string
	Graph       graph.Graph[string, string]
	BlockLookup map[string]Block
	AdjMap
}

func FromGraph(workflowGraph *gen.Graph) (*Workflow, error) {
	g := graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())

	blockLookup := map[string]Block{}
	for _, node := range workflowGraph.Nodes {
		err := g.AddVertex(node.Id)
		if err != nil {
			return nil, err
		}

		// add block to lookup to be used for execution
		activity, err := NewBlock(node)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating block for node %s", node.Id)
		}
		blockLookup[node.Id] = activity
	}

	for _, edge := range workflowGraph.Edges {
		err := g.AddEdge(edge.From, edge.To)
		if err != nil {
			return nil, err
		}
	}

	adjMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, errors.Wrapf(err, "error getting adjacency map")
	}

	return &Workflow{
		ID:          workflowGraph.Id,
		Graph:       g,
		BlockLookup: blockLookup,
		AdjMap:      adjMap,
	}, nil
}

func (w *Workflow) Run(logger Logger, executor Executor, nodeID string) (*Result, error) {
	vert, err := w.Graph.Vertex(nodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting vertex %s", nodeID)
	}

	res, err := w.traverseWorkflow(logger, executor, vert, Input{
		Params: nil,
	})
	if err != nil {
		logger.Error("Error traversing workflow", "error", err)
		return nil, nil
	}
	return res, nil
}

func (w *Workflow) traverseWorkflow(logger Logger, executor Executor, vert string, input Input) (*Result, error) {
	for neighbor := range w.AdjMap[vert] {
		block, ok := w.BlockLookup[neighbor]
		if !ok {
			return nil, fmt.Errorf("vertex not found: %s", neighbor)
		}

		res, err := block.Execute(executor, input)
		if err != nil {
			return nil, errors.Wrapf(err, "error executing block %s", neighbor)
		}

		nextBlockInput := Input{
			Params: res.Data,
		}

		log.Debug().Interface("result", res).Msg("block result")

		logger.Info("Traversing workflow", "nodeID", neighbor)
		_, err = w.traverseWorkflow(logger, executor, neighbor, nextBlockInput)
		if err != nil {
			return nil, errors.Wrapf(err, "error traversing workflow %s", neighbor)
		}
	}
	return &Result{}, nil
}
