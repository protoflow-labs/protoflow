package workflow

import (
	"github.com/dominikbraun/graph"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	worknode "github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/rs/zerolog/log"
)

func FromProject(project *gen.Project) (*Workflow, error) {
	g := graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())

	if project.Graph == nil {
		return nil, errors.New("project graph is nil")
	}

	resources := resource.DependencyProvider{}
	for _, protoRes := range project.Resources {
		r, err := resource.FromProto(protoRes)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating resource for node %s", protoRes.Id)
		}
		resources[protoRes.Id] = r
	}

	for _, r := range resources {
		if err := r.ResolveDependencies(resources); err != nil {
			return nil, errors.Wrapf(err, "error resolving dependencies for resource %s", r.ID())
		}
	}

	nodeLookup := map[string]worknode.Node{}
	for _, node := range project.Graph.Nodes {
		if node.Id == "" {
			return nil, errors.New("node id cannot be empty")
		}
		err := g.AddVertex(node.Id)
		if err != nil {
			return nil, err
		}

		// add block to lookup to be used for execution
		builtNode, err := worknode.NewNode(node)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating block for node %s", node.Id)
		}
		nodeLookup[node.Id] = builtNode
		r := resources[node.ResourceId]
		if r != nil {
			r.AddNode(builtNode)
		} else {
			// TODO breadchris inputs do not have resources, but should they?
			log.Warn().Str("node", builtNode.NormalizedName()).Msg("no resource found for node")
		}
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

	preMap, err := g.PredecessorMap()
	if err != nil {
		return nil, errors.Wrapf(err, "error getting predecessor map")
	}

	return &Workflow{
		// TODO breadchris this should be a deterministic value based on the workflow node slice
		ID:         uuid.NewString(),
		ProjectID:  project.Id,
		Graph:      g,
		NodeLookup: nodeLookup,
		AdjMap:     adjMap,
		PreMap:     preMap,
		Resources:  resources,
	}, nil
}
