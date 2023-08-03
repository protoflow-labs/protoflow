package workflow

import (
	"context"
	"fmt"
	graphlib "github.com/dominikbraun/graph"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/edge"
	"github.com/protoflow-labs/protoflow/pkg/graph/node"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const edgeIDKey = "id"

type AdjMap map[string]map[string]graphlib.Edge[string]

// TODO breadchris can this be a map[string]Resource?
type Instances map[string]graph.Node

// Workflow is a directed graph of nodes that represent a workflow. The builder interface is immutable allowing for extensions.
type Workflow struct {
	ID        string
	ProjectID string
	Graph     graphlib.Graph[string, string]
	// TODO breadchris this should be a deterministic value based on the workflow node slice
	NodeLookup map[string]graph.Node
	EdgeLookup map[string]graph.Edge
}

type NodeBuilder[T graph.ProtoNode] func(message T) graph.Node
type EdgeBuilder[U graph.ProtoEdge] func(message U) graph.Edge

type WorkflowBuilder[T graph.ProtoNode, U graph.ProtoEdge] struct {
	w            *Workflow
	nodeBuilders map[protoreflect.FullName]NodeBuilder[T]
	edgeBuilders map[protoreflect.FullName]EdgeBuilder[U]
	err          error
}

// TODO breadchris is a workflow builder type needed so that you separate the build and workflow steps?

func Default() *WorkflowBuilder[*gen.Node, *gen.Edge] {
	return NewBuilder[*gen.Node, *gen.Edge]().
		WithNodeTypes(&gen.Node{}, node.New).
		WithEdgeTypes(&gen.Edge{}, edge.New)
}

func NewBuilder[T graph.ProtoNode, U graph.ProtoEdge]() *WorkflowBuilder[T, U] {
	return &WorkflowBuilder[T, U]{
		nodeBuilders: map[protoreflect.FullName]NodeBuilder[T]{},
		edgeBuilders: map[protoreflect.FullName]EdgeBuilder[U]{},
		w: &Workflow{
			ProjectID:  uuid.NewString(),
			NodeLookup: map[string]graph.Node{},
			EdgeLookup: map[string]graph.Edge{},
			Graph:      graphlib.New(graphlib.StringHash, graphlib.Directed(), graphlib.PreventCycles()),
		},
	}
}

func builderName(n protoreflect.ProtoMessage) protoreflect.FullName {
	return n.ProtoReflect().Descriptor().FullName()
}

func (w *WorkflowBuilder[T, U]) newNode(n T) (graph.Node, error) {
	//configOneOf := n.ProtoReflect().Descriptor().Oneofs().Get(0)
	//fd := n.ProtoReflect().WhichOneof(configOneOf)
	//nodeType := fd.Message().FullName()
	bn := builderName(n)
	nodeBuilder, ok := w.nodeBuilders[bn]
	if !ok {
		return nil, errors.Errorf("no builder found for node type %s", bn)
	}
	builtNode := nodeBuilder(n)
	if builtNode == nil {
		return nil, errors.Errorf("node builder for node type %+v returned nil", n)
	}
	return builtNode, nil
}

func (w *WorkflowBuilder[T, U]) newEdge(n U) (graph.Edge, error) {
	bn := builderName(n)
	edgeBuilder, ok := w.edgeBuilders[bn]
	if !ok {
		return nil, errors.Errorf("no builder found for edge type %s", bn)
	}
	builtEdge := edgeBuilder(n)
	if builtEdge == nil {
		return nil, errors.Errorf("edge builder for edge type %+v returned nil", n)
	}
	return builtEdge, nil
}

func (w *WorkflowBuilder[T, U]) WithNodeTypes(n T, builder NodeBuilder[T]) *WorkflowBuilder[T, U] {
	nw := *w
	nw.nodeBuilders[builderName(n)] = builder
	return &nw
}

func (w *WorkflowBuilder[T, U]) WithEdgeTypes(n U, builder EdgeBuilder[U]) *WorkflowBuilder[T, U] {
	nw := *w
	nw.edgeBuilders[builderName(n)] = builder
	return &nw
}

func (w *WorkflowBuilder[T, U]) WithNodes(nodes ...T) *WorkflowBuilder[T, U] {
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

func (w *WorkflowBuilder[T, U]) WithBuiltNodes(nodes ...graph.Node) *WorkflowBuilder[T, U] {
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

func (w *WorkflowBuilder[T, U]) WithBuiltEdges(edges ...graph.Edge) *WorkflowBuilder[T, U] {
	nw := *w
	for _, e := range edges {
		err := nw.addBuiltEdge(e)
		if err != nil {
			nw.err = errors.Wrapf(err, "error adding edge from %s to %s", e.From(), e.To())
			return &nw
		}
	}
	return &nw
}

func (w *WorkflowBuilder[T, U]) addNode(n T) error {
	builtNode, err := w.newNode(n)
	if err != nil {
		return errors.Wrapf(err, "error creating block for node %s", n.GetId())
	}
	return w.addBuiltNode(builtNode)
}

func (w *WorkflowBuilder[T, U]) addBuiltNode(builtNode graph.Node) error {
	if builtNode.ID() == "" {
		return errors.New("node id cannot be empty")
	}
	err := w.w.Graph.AddVertex(builtNode.ID())
	if err != nil {
		return errors.Wrapf(err, "error adding node %s", builtNode.ID())
	}

	w.w.NodeLookup[builtNode.ID()] = builtNode
	return nil
}

func (w *WorkflowBuilder[T, U]) addEdge(edge U) error {
	builtEdge, err := w.newEdge(edge)
	if err != nil {
		return errors.Wrapf(err, "error creating edge for edge %s", edge.GetId())
	}
	return w.addBuiltEdge(builtEdge)
}

func (w *WorkflowBuilder[T, U]) addBuiltEdge(builtEdge graph.Edge) error {
	if builtEdge.ID() == "" {
		return errors.New("edge id cannot be empty")
	}
	err := w.w.Graph.AddEdge(builtEdge.From(), builtEdge.To(), graphlib.EdgeAttribute(edgeIDKey, builtEdge.ID()))
	if err != nil {
		return errors.Wrapf(err, "error adding edge %s", builtEdge.ID())
	}
	w.w.EdgeLookup[builtEdge.ID()] = builtEdge
	return nil
}

func (w *WorkflowBuilder[T, U]) WithProtoProject(project graph.ProtoProject[T, U]) *WorkflowBuilder[T, U] {
	nw := *w
	var g = project.GetGraph()
	if g == nil {
		nw.err = errors.New("project graph is nil")
		return &nw
	}
	nw.w.ProjectID = project.GetId()

	for _, n := range g.GetNodes() {
		err := nw.addNode(n)
		if err != nil {
			nw.err = errors.Wrapf(err, "error adding node %s", n.GetId())
			return &nw
		}
	}

	for _, e := range g.GetEdges() {
		if e.GetFrom() == "" || e.GetTo() == "" {
			nw.err = errors.Errorf("edge %s has empty from (%s) or to (%s)", e.GetId(), e.GetFrom(), e.GetTo())
			return &nw
		}
		err := nw.addEdge(e)
		if err != nil {
			nw.err = errors.Wrapf(err, "error adding edge %s", e.GetId())
			return &nw
		}
	}
	return &nw
}

func (w *WorkflowBuilder[T, U]) Build() (*Workflow, error) {
	nw := *w

	// TODO breadchris resolve dependencies
	// need to look at edge types nw.w.Graph.Edge

	sucMap, err := nw.w.Graph.AdjacencyMap()
	if err != nil {
		return nil, errors.Wrapf(err, "error getting adjacency map")
	}

	for nodeID, sucs := range sucMap {
		from, _ := nw.w.NodeLookup[nodeID]
		for sID, e := range sucs {
			to, _ := nw.w.NodeLookup[sID]

			eID := e.Properties.Attributes[edgeIDKey]
			builtEdge, _ := w.w.EdgeLookup[eID]

			err = builtEdge.Connect(from, to)
			if err != nil {
				return nil, errors.Wrapf(err, "error connecting edge %s", eID)
			}
		}
	}

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

func (w *Workflow) GetNodeProvider(id string) (graph.Node, error) {
	n, err := w.GetNode(id)
	if err != nil {
		return nil, err
	}
	return n.Provider()
}

// TODO breadchris nodeID should not be needed, the workflow should already be a slice of the graph that is configured to run
// WireNodes wires the nodes in the workflow together and returns an observable that can be subscribed to. Nodes are executed when an event is received on the input observable.
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
	for id, r := range w.NodeLookup {
		cleanup, err := r.Init()
		if err != nil {
			return nil, errors.Wrapf(err, "error creating resource %s", r.NormalizedName())
		}
		cleanupFuncs = append(cleanupFuncs, cleanup)
		instances[id] = r
	}

	connector := NewConnector()
	connector.Add(input)

	// TODO breadchris use actual output from the workflow
	// wire an input into the workflow so that data can flow between nodes
	n, ok := w.NodeLookup[nodeID]
	if !ok {
		return nil, fmt.Errorf("vertex not found: %s", nodeID)
	}

	_, err := w.wireWorkflow(ctx, connector, instances, n, graph.IO{
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
	node graph.Node,
	input graph.IO,
) (*graph.IO, error) {
	log.Debug().
		Str("node", node.NormalizedName()).
		// Interface("resource", input.Resource.Name()).
		Msg("wiring node IO")
	nextInput, err := node.Wire(ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "error executing node: %s", node.NormalizedName())
	}

	if nextInput.Observable == nil {
		return nil, fmt.Errorf("node %s returned nil observable", node.NormalizedName())
	}

	connector.Add(nextInput.Observable)

	for _, listener := range node.Subscribers() {
		child := listener.GetNode()

		nextInput, err = listener.Transform(ctx, nextInput)
		if err != nil {
			return nil, errors.Wrapf(err, "error transforming input for node %s", child.NormalizedName())
		}

		// TODO breadchris for any publishers to this node, add the child as a subscriber
		//for _, parent := range child.Publishers() {
		//
		//}

		log.Debug().
			Str("node", node.NormalizedName()).
			Str("child", child.NormalizedName()).
			Msg("traversing workflow")

		// TODO breadchris what to do with the output here? Map over the output to turn into NodeExecution
		_, err = w.wireWorkflow(ctx, connector, instances, child, nextInput)
		if err != nil {
			return nil, errors.Wrapf(err, "error traversing workflow %s", child)
		}
	}
	return nil, nil
}
