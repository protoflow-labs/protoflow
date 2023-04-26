package workflow

import (
	"fmt"

	"github.com/dominikbraun/graph"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/rs/zerolog/log"
)

type Input struct {
	Params       interface{}
	Resources    map[string]any
	Dependencies []string
}

type Result struct {
	Data interface{}
}

type AdjMap map[string]map[string]graph.Edge[string]

type Workflow struct {
	ID         string
	Graph      graph.Graph[string, string]
	NodeLookup map[string]Node
	AdjMap
	Resources map[string]Resource
}

func FromProject(project *gen.Project) (*Workflow, error) {
	g := graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())

	projectResources := map[string]Resource{}
	blockLookup := map[string]*gen.Block{}
	for _, resource := range project.Resources {
		for _, block := range resource.Blocks {
			blockLookup[block.Id] = block
		}

		r, err := ResourceFromProto(resource)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating resource for node %s", resource.Id)
		}
		projectResources[resource.Id] = r
	}

	nodeLookup := map[string]Node{}
	for _, node := range project.Graph.Nodes {
		err := g.AddVertex(node.Id)
		if err != nil {
			return nil, err
		}

		block, ok := blockLookup[node.BlockId]
		if !ok {
			return nil, fmt.Errorf("block not found: %s", node.BlockId)
		}

		// add block to lookup to be used for execution
		builtNode, err := NewNode(node, block)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating block for node %s", node.Id)
		}
		nodeLookup[node.Id] = builtNode
	}

	for _, edge := range project.Graph.Edges {
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
		ID:         project.Id,
		Graph:      g,
		NodeLookup: nodeLookup,
		AdjMap:     adjMap,
		Resources:  projectResources,
	}, nil
}

// TODO breadchris can this be a map[string]Resource?
type Instances map[string]any

func (w *Workflow) Run(logger Logger, executor Executor, nodeID string) (*Result, error) {
	var cleanupFuncs []func()
	defer func() {
		for _, cleanup := range cleanupFuncs {
			cleanup()
		}
	}()

	// TODO breadchris implement resource pool to avoid creating resources for every workflow
	instances := Instances{}
	for id, resource := range w.Resources {
		cleanup, err := resource.Init()
		if err != nil {
			return nil, errors.Wrapf(err, "error creating resource %s", id)
		}
		cleanupFuncs = append(cleanupFuncs, cleanup)
		instances[id] = resource
	}

	vert, err := w.Graph.Vertex(nodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting vertex %s", nodeID)
	}

	res, err := w.traverseWorkflow(logger, instances, executor, vert, Input{
		Params: nil,
	})
	if err != nil {
		logger.Error("Error traversing workflow", "error", err)
		return nil, nil
	}
	return res, nil
}

func (w *Workflow) traverseWorkflow(logger Logger, instances Instances, executor Executor, vert string, input Input) (*Result, error) {
	block, ok := w.NodeLookup[vert]
	if !ok {
		return nil, fmt.Errorf("vertex not found: %s", vert)
	}

	input.Resources = instances
	res, err := block.Execute(executor, input)
	if err != nil {
		return nil, errors.Wrapf(err, "error executing block: %s", vert)
	}

	nextBlockInput := Input{
		Params:    res.Data,
		Resources: instances,
	}

	log.Debug().Interface("result", res).Msg("block result")
	for neighbor := range w.AdjMap[vert] {
		logger.Info("Traversing workflow", "nodeID", neighbor)
		_, err = w.traverseWorkflow(logger, instances, executor, neighbor, nextBlockInput)
		if err != nil {
			return nil, errors.Wrapf(err, "error traversing workflow %s", neighbor)
		}
	}
	return &Result{
		Data: res.Data,
	}, nil
}
