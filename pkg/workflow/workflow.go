package workflow

import (
	"fmt"

	"github.com/dominikbraun/graph"
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
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

func FromProject(project *gen.Project, cache cache.Cache) (*Workflow, error) {
	g := graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())

	// TODO breadchris this should not be hardcoded, this should be provided to the service when it is created?
	projectResources := map[string]Resource{
		"js": &LanguageServiceResource{
			LanguageService: &gen.LanguageService{
				Runtime: gen.Runtime_NODE,
				Host:    "localhost:8086",
			},
			cache: cache,
		},
		"docs": &DocstoreResource{
			Docstore: &gen.Docstore{
				Url: "mem://",
			},
		},
		"bucket": &BlobstoreResource{
			Blobstore: &gen.Blobstore{
				Url: "file:///home/breadchris/.protoflow",
			},
		},
	}

	// TODO breadchris blocks will be used in the future to associate with nodes, but for now they are not used
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

		// add block to lookup to be used for execution
		builtNode, err := NewNode(projectResources, node)
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
type Instances map[string]Resource

// TODO breadchris nodeID should not be needed, the workflow should already be a slice of the graph that is configured to run
func (w *Workflow) Run(logger Logger, executor Executor, nodeID string, input string) (*Result, error) {
	var cleanupFuncs []func()
	defer func() {
		for _, cleanup := range cleanupFuncs {
			if cleanup != nil {
				cleanup()
			}
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
		Params: input,
	})
	if err != nil {
		logger.Error("Error traversing workflow", "error", err)
		return nil, err
	}
	return res, nil
}

func (w *Workflow) traverseWorkflow(logger Logger, instances Instances, executor Executor, vert string, input Input) (*Result, error) {
	node, ok := w.NodeLookup[vert]
	if !ok {
		return nil, fmt.Errorf("vertex not found: %s", vert)
	}

	log.Debug().Str("vert", vert).Msg("injecting dependencies for node")
	input.Resources = map[string]any{}
	for _, resourceID := range node.Dependencies() {
		resource, ok := instances[resourceID]
		if !ok {
			return nil, fmt.Errorf("resource not found: %s", resourceID)
		}
		if _, ok := input.Resources[resource.Name()]; ok {
			logger.Warn("Resource type already exists in input", "resource", resource.Name())
			continue
		}
		input.Resources[resource.Name()] = resource
	}

	log.Debug().Str("vert", vert).Msg("executing node")
	res, err := node.Execute(executor, input)
	if err != nil {
		return nil, errors.Wrapf(err, "error executing node: %s", vert)
	}

	nextBlockInput := Input{
		Params: res.Data,
	}

	log.Debug().Interface("result", res).Msg("node result")
	nextResSet := false
	for neighbor := range w.AdjMap[vert] {
		logger.Info("Traversing workflow", "nodeID", neighbor)
		neighborRes, err := w.traverseWorkflow(logger, instances, executor, neighbor, nextBlockInput)
		if err != nil {
			return nil, errors.Wrapf(err, "error traversing workflow %s", neighbor)
		}

		if !nextResSet {
			res = neighborRes
			nextResSet = true
		}
	}
	return &Result{
		Data: res.Data,
	}, nil
}
