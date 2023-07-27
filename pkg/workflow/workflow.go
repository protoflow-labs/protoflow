package workflow

import (
	"context"
	"fmt"
	graphlib "github.com/dominikbraun/graph"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const edgeIDKey = "id"

type AdjMap map[string]map[string]graphlib.Edge[string]

// TODO breadchris can this be a map[string]Resource?
type Instances map[string]graph.Resource

// Workflow is a directed graph of nodes that represent a workflow. The builder interface is immutable allowing for extensions.
type Workflow struct {
	ID        string
	ProjectID string
	Graph     graphlib.Graph[string, string]
	// TODO breadchris this should be a deterministic value based on the workflow node slice
	NodeLookup map[string]graph.Node
	AdjMap     AdjMap
	PreMap     AdjMap
	Resources  map[string]graph.Resource
}

type Builder[T graph.ProtoNode] func(message T) (graph.Node, error)

type WorkflowBuilder[T graph.ProtoNode] struct {
	w            *Workflow
	nodeBuilders map[protoreflect.FullName]Builder[T]
	err          error
}

// TODO breadchris is a workflow builder type needed so that you separate the build and workflow steps?

func Default() *WorkflowBuilder[*gen.Node] {
	return NewBuilder[*gen.Node]().WithNodeTypes(&gen.Node{}, node.NewNode)
}

func NewBuilder[T graph.ProtoNode]() *WorkflowBuilder[T] {
	return &WorkflowBuilder[T]{
		nodeBuilders: map[protoreflect.FullName]Builder[T]{},
		w: &Workflow{
			ProjectID: uuid.NewString(),
			// TODO breadchris can reseources be replaced with a graph?
			Resources:  graph.DependencyProvider{},
			NodeLookup: map[string]graph.Node{},
			Graph:      graphlib.New(graphlib.StringHash, graphlib.Directed(), graphlib.PreventCycles()),
		},
	}
}

func builderName(n protoreflect.ProtoMessage) protoreflect.FullName {
	return n.ProtoReflect().Descriptor().FullName()
}

func (w *WorkflowBuilder[T]) newNode(n T) (graph.Node, error) {
	//configOneOf := n.ProtoReflect().Descriptor().Oneofs().Get(0)
	//fd := n.ProtoReflect().WhichOneof(configOneOf)
	//nodeType := fd.Message().FullName()
	bn := builderName(n)
	nodeBuilder, ok := w.nodeBuilders[bn]
	if !ok {
		return nil, errors.Errorf("no builder found for node type %s", bn)
	}
	return nodeBuilder(n)
}

func (w *WorkflowBuilder[T]) WithNodeTypes(n T, builder Builder[T]) *WorkflowBuilder[T] {
	nw := *w
	nw.nodeBuilders[builderName(n)] = builder
	return &nw
}

func (w *WorkflowBuilder[T]) WithNodes(nodes ...T) *WorkflowBuilder[T] {
	nw := *w
	for _, n := range nodes {
		err := nw.addNode(n)
		if err != nil {
			nw.err = errors.Wrapf(err, "error adding node %s", n.GetId())
			return &nw
		}
	}
	return &nw
}

func (w *WorkflowBuilder[T]) WithBuiltNodes(nodes ...graph.Node) *WorkflowBuilder[T] {
	nw := *w
	for _, n := range nodes {
		err := nw.addBuiltNode(n)
		if err != nil {
			nw.err = errors.Wrapf(err, "error adding node %s", n.ID())
			return &nw
		}
	}
	return &nw
}

func (w *WorkflowBuilder[T]) WithBuiltEdges(edges ...graph.Edge) *WorkflowBuilder[T] {
	nw := *w
	for _, e := range edges {
		err := nw.addEdge(&gen.Edge{
			Id:   uuid.NewString(),
			From: e.From.ID(),
			To:   e.To.ID(),
		})
		if err != nil {
			nw.err = errors.Wrapf(err, "error adding edge from %s to %s", e.From.ID(), e.To.ID())
			return &nw
		}
	}
	return &nw
}

func (w *WorkflowBuilder[T]) WithResource(r *gen.Resource) *WorkflowBuilder[T] {
	nw := *w
	err := nw.addResource(r)
	if err != nil {
		nw.err = errors.Wrapf(err, "error adding resource %s", r.Id)
		return &nw
	}
	return &nw
}

func (w *WorkflowBuilder[T]) addNode(n T) error {
	builtNode, err := w.newNode(n)
	if err != nil {
		return errors.Wrapf(err, "error creating block for node %s", n.GetId())
	}
	return w.addBuiltNode(builtNode)
}

func (w *WorkflowBuilder[T]) addBuiltNode(builtNode graph.Node) error {
	if builtNode.ID() == "" {
		return errors.New("node id cannot be empty")
	}
	err := w.w.Graph.AddVertex(builtNode.ID())
	if err != nil {
		return errors.Wrapf(err, "error adding node %s", builtNode.ID())
	}

	w.w.NodeLookup[builtNode.ID()] = builtNode
	r := w.w.Resources[builtNode.ResourceID()]
	if r != nil {
		r.AddNode(builtNode)
	} else {
		// TODO breadchris inputs do not have resources, but should they?
		log.Warn().Str("node", builtNode.NormalizedName()).Msg("no resource found for node")
	}
	return nil
}

// TODO breadchris change to generic edge
func (w *WorkflowBuilder[T]) addEdge(edge *gen.Edge) error {
	return w.w.Graph.AddEdge(edge.From, edge.To, graphlib.EdgeAttribute(edgeIDKey, edge.Id))
}

// TODO breadchris change to generic resource
func (w *WorkflowBuilder[T]) addResource(protoRes *gen.Resource) error {
	r, err := resource.FromProto(protoRes)
	if err != nil {
		return errors.Wrapf(err, "error creating resource for node %s", protoRes.Id)
	}
	w.w.Resources[protoRes.Id] = r
	return nil
}

func (w *WorkflowBuilder[T]) WithProtoProject(project graph.ProtoProject[T]) *WorkflowBuilder[T] {
	nw := *w
	var g = project.GetGraph()
	if g == nil {
		nw.err = errors.New("project graph is nil")
		return &nw
	}
	nw.w.ProjectID = project.GetId()

	for _, protoRes := range project.GetResources() {
		err := nw.addResource(protoRes)
		if err != nil {
			nw.err = errors.Wrapf(err, "error adding resource %s", protoRes.Id)
			return &nw
		}
	}

	for _, n := range g.GetNodes() {
		err := nw.addNode(n)
		if err != nil {
			nw.err = errors.Wrapf(err, "error adding node %s", n.GetId())
			return &nw
		}
	}

	for _, edge := range g.GetEdges() {
		err := nw.addEdge(edge)
		if err != nil {
			nw.err = errors.Wrapf(err, "error adding edge %s", edge.Id)
			return &nw
		}
	}
	return &nw
}

func (w *WorkflowBuilder[T]) Build() (*Workflow, error) {
	nw := *w

	for _, r := range nw.w.Resources {
		if err := r.ResolveDependencies(nw.w.Resources); err != nil {
			return nil, errors.Wrapf(err, "error resolving dependencies for resource %s", r.ID())
		}
	}

	adjMap, err := nw.w.Graph.AdjacencyMap()
	if err != nil {
		return nil, errors.Wrapf(err, "error getting adjacency map")
	}
	nw.w.AdjMap = adjMap

	preMap, err := nw.w.Graph.PredecessorMap()
	if err != nil {
		return nil, errors.Wrapf(err, "error getting predecessor map")
	}
	nw.w.PreMap = preMap

	nw.w.ID = uuid.NewString()
	if nw.err != nil {
		return nil, nw.err
	}
	return nw.w, nil
}

func (w *Workflow) GetNode(id string) (graph.Node, error) {
	node, ok := w.NodeLookup[id]
	if !ok {
		return nil, fmt.Errorf("node with id %s not found", id)
	}
	return node, nil
}

func (w *Workflow) GetNodeResource(id string) (graph.Resource, error) {
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
func (w *Workflow) WireNodes(ctx context.Context, nodeID string, input rxgo.Observable) (rxgo.Observable, error) {
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

	// TODO breadchris use actual output from the workflow
	// wire an input into the workflow so that data can flow between nodes
	_, err = w.wireWorkflow(ctx, connector, instances, vert, graph.Input{
		Observable: input,
	})
	if err != nil {
		log.Error().Err(err).Msgf("failed to traverse workflow")
		return nil, err
	}
	// TODO breadchris this is returning a trace, not the actual output
	return connector.Connect(ctx), nil
}

func (w *Workflow) wireWorkflow(
	ctx context.Context,
	connector *Connector,
	instances Instances,
	nodeID string,
	input graph.Input,
) (*graph.Output, error) {
	node, ok := w.NodeLookup[nodeID]
	if !ok {
		return nil, fmt.Errorf("vertex not found: %s", nodeID)
	}

	log.Debug().Str("node", node.NormalizedName()).Msg("injecting dependencies for node")
	err := injectDepsForNode(instances, &input, node)
	if err != nil {
		return nil, errors.Wrapf(err, "error injecting dependencies for node %s", nodeID)
	}

	log.Debug().
		Str("node", node.NormalizedName()).
		// Interface("resource", input.Resource.Name()).
		Msg("wiring node IO")
	output, err := node.Wire(ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "error executing node: %s", nodeID)
	}

	if output.Observable == nil {
		return nil, fmt.Errorf("node %s returned nil observable", nodeID)
	}

	connector.Add(output.Observable)

	nextBlockInput := graph.Input{
		Observable: output.Observable,
	}

	for neighbor := range w.AdjMap[nodeID] {
		//e, err := w.Graph.Edge(nodeID, neighbor)
		//if err != nil {
		//	return nil, errors.Wrapf(err, "error getting edge between %s and %s", nodeID, neighbor)
		//}
		//edgeID := e.Properties.Attributes[edgeIDKey]
		//nextBlockInput.Observable = output.Observable.Map(func(c context.Context, i any) (any, error) {
		//	vm := otto.New()
		//	vm.Set("input", i)
		//	vm.Run(`input = input.message`)
		//	v, err := vm.Get("input")
		//	if err != nil {
		//		log.Error().Err(err).Msg("error getting input")
		//		return nil, err
		//	}
		//	o, err := v.Export()
		//	if err != nil {
		//		log.Error().Err(err).Msg("error exporting input")
		//		return nil, err
		//	}
		//	switch o.(type) {
		//	case map[string]any:
		//		log.Debug().Interface("input", o).Msg("input is map")
		//	}
		//	return o, nil
		//})

		log.Debug().
			Str("node", node.NormalizedName()).
			Str("neighbor", w.NodeLookup[neighbor].NormalizedName()).
			Msg("traversing workflow")

		// TODO breadchris what to do with the output here? Map over the output to turn into NodeExecution
		_, err = w.wireWorkflow(ctx, connector, instances, neighbor, nextBlockInput)
		if err != nil {
			return nil, errors.Wrapf(err, "error traversing workflow %s", neighbor)
		}
	}
	return nil, nil
}

func injectDepsForNode(instances Instances, input *graph.Input, node graph.Node) error {
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
