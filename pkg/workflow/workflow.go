package workflow

import (
	"context"
	"fmt"
	"github.com/dominikbraun/graph"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/workflow/execute"
	worknode "github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
)

type AdjMap map[string]map[string]graph.Edge[string]

// TODO breadchris can this be a map[string]Resource?
type Instances map[string]resource.Resource

type Workflow struct {
	ID         string
	ProjectID  string
	Graph      graph.Graph[string, string]
	NodeLookup map[string]worknode.Node
	AdjMap     AdjMap
	PreMap     AdjMap
	Resources  map[string]resource.Resource
}

func (w *Workflow) GetNode(id string) (worknode.Node, error) {
	node, ok := w.NodeLookup[id]
	if !ok {
		return nil, fmt.Errorf("node with id %s not found", id)
	}
	return node, nil
}

func (w *Workflow) GetNodeResource(id string) (resource.Resource, error) {
	node, err := w.GetNode(id)
	if err != nil {
		return nil, err
	}
	r, ok := w.Resources[node.ResourceID()]
	if !ok {
		return nil, fmt.Errorf("resource %s not found", node.ResourceID())
	}
	return r, nil
}

// TODO breadchris nodeID should not be needed, the workflow should already be a slice of the graph that is configured to run
func (w *Workflow) Run(ctx context.Context, logger Logger, executor execute.Executor, nodeID string, input rxgo.Observable) (rxgo.Observable, error) {
	var cleanupFuncs []func()
	defer func() {
		for _, cleanup := range cleanupFuncs {
			if cleanup != nil {
				cleanup()
			}
		}
	}()

	// TODO breadchris make a slice of resources that are needed for the workflow
	// TODO breadchris implement resource pool to avoid creating resources for every workflow
	instances := Instances{}
	for id, r := range w.Resources {
		cleanup, err := r.Init()
		if err != nil {
			return nil, errors.Wrapf(err, "error creating resource %s", r.Name())
		}
		cleanupFuncs = append(cleanupFuncs, cleanup)
		instances[id] = r
	}

	vert, err := w.Graph.Vertex(nodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting vertex %s", nodeID)
	}

	connector := NewConnector()
	connector.Add(input)

	err = w.traverseWorkflow(ctx, connector, instances, executor, vert, execute.Input{
		Observable: input,
	})
	if err != nil {
		logger.Error("failed to traverse workflow", "error", err)
		return nil, err
	}
	return connector.Connect(ctx), nil
}

func (w *Workflow) traverseWorkflow(
	ctx context.Context,
	connector *Connector,
	instances Instances,
	executor execute.Executor,
	nodeID string,
	input execute.Input,
) error {
	node, ok := w.NodeLookup[nodeID]
	if !ok {
		return fmt.Errorf("vertex not found: %s", nodeID)
	}

	log.Debug().Str("node", node.NormalizedName()).Msg("injecting dependencies for node")
	err := injectDepsForNode(instances, &input, node)
	if err != nil {
		return errors.Wrapf(err, "error injecting dependencies for node %s", nodeID)
	}

	log.Debug().
		Str("node", node.NormalizedName()).
		Interface("resource", input.Resource).
		Msg("wiring node IO")
	output, err := executor.Execute(node, input)
	if err != nil {
		return errors.Wrapf(err, "error executing node: %s", nodeID)
	}

	connector.Add(output.Observable)

	nextBlockInput := execute.Input{
		Observable: output.Observable,
	}

	for neighbor := range w.AdjMap[nodeID] {
		log.Debug().
			Str("node", node.NormalizedName()).
			Str("neighbor", w.NodeLookup[neighbor].NormalizedName()).
			Msg("traversing workflow")

		err = w.traverseWorkflow(ctx, connector, instances, executor, neighbor, nextBlockInput)
		if err != nil {
			return errors.Wrapf(err, "error traversing workflow %s", neighbor)
		}
	}
	return nil
}

func injectDepsForNode(instances Instances, input *execute.Input, node worknode.Node) error {
	if node.ResourceID() == "" {
		return nil
	}
	r, ok := instances[node.ResourceID()]
	if !ok {
		return fmt.Errorf("resource not found: %s", node.ResourceID())
	}
	input.Resource = r
	return nil
}
