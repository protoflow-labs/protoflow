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
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/data"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

const edgeIDKey = "id"

// Workflow is a directed graph of nodes that represent a workflow. The builder interface is immutable allowing for extensions.
type Workflow struct {
	ID        string
	ProjectID string
	Graph     graphlib.Graph[string, string]
	// TODO breadchris this should be a deterministic value based on the workflow node slice
	NodeLookup map[string]graph.Node
	EdgeLookup map[string]graph.Edge
}

// NodeBuilder is a function that builds a node from a proto message.
type NodeBuilder[T graph.ProtoNode] func(message T) graph.Node

// EdgeBuilder is a function that builds an edge from a proto message.
type EdgeBuilder[U graph.ProtoEdge] func(message U) graph.Edge

// WorkflowBuilder is a builder for a workflow. It is immutable allowing for extensions.
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

func (w *Workflow) EnumerateProviders() []*gen.EnumeratedProvider {
	var providers []*gen.EnumeratedProvider
	for _, n := range w.NodeLookup {
		info := &gen.ProviderInfo{
			State: gen.ProviderState_READY,
			Error: "",
		}
		nl, ok := w.NodeLookup[n.ID()]
		if !ok {
			info.State = gen.ProviderState_ERROR
			info.Error = "node not found"
		}

		providedNodes, err := nl.Provide()
		if len(providedNodes) == 0 && err == nil {
			continue
		}

		if err != nil {
			info.State = gen.ProviderState_ERROR
			info.Error = err.Error()
		}
		providers = append(providers, &gen.EnumeratedProvider{
			Provider: &gen.NodeDetails{
				Id:   nl.ID(),
				Name: nl.Name(),
			},
			Nodes: providedNodes,
			Info:  info,
		})
	}
	return providers
}

func (w *Workflow) GetNode(id string) (graph.Node, error) {
	nl, ok := w.NodeLookup[id]
	if !ok {
		return nil, fmt.Errorf("node with id %s not found", id)
	}
	return nl, nil
}

func (w *Workflow) GetNodeProvider(id string) (graph.Node, error) {
	n, err := w.GetNode(id)
	if err != nil {
		return nil, err
	}
	return n.Provider()
}

// WireNodes wires the nodes in the workflow together and returns an observable that can be subscribed to. Nodes are executed when an event is received on the input observable.
func (w *Workflow) WireNodes(
	ctx context.Context,
	nodeID string,
	input rxgo.Observable,
	manager *Manager,
) (*graph.IO, error) {
	getEdge := func(from, to string) (graph.Edge, error) {
		e, err := w.Graph.Edge(from, to)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting edge from %s to %s", from, to)
		}
		eID := e.Properties.Attributes[edgeIDKey]
		builtEdge, ok := w.EdgeLookup[eID]
		if !ok {
			log.Error().Msgf("edge %s not found", eID)
			return nil, errors.Errorf("edge %s not found", eID)
		}
		return builtEdge, nil
	}

	// TODO breadchris I like the idea of this, but needs to be cleaned up, should probably be closer to the creation of the workflow?
	in := data.NewInputNode(base.NewNode("input"), data.NewInputProto().GetInput(), data.WithObservable(&graph.IO{
		Observable: input,
	}))
	w.NodeLookup[in.ID()] = in

	e := edge.New(edge.NewMapProto(in.ID(), nodeID))
	w.EdgeLookup[e.ID()] = e

	n := w.NodeLookup[nodeID]
	err := e.Connect(in, n)
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting edge %s", e.ID())
	}

	// TODO breadchris move workflow slicing into its own function
	// search forward and backward in the graph to create a subgraph of nodes that are connected to the node we want to wire
	subgraph := graphlib.New(graphlib.StringHash, graphlib.Directed(), graphlib.PreventCycles())
	_ = subgraph.AddVertex(in.ID())
	_ = subgraph.AddEdge(in.ID(), nodeID)

	discovered := []string{nodeID}
	err = util.BFS(w.Graph, nodeID, true, func(from, to string) bool {
		e, err := getEdge(from, to)
		if err != nil {
			log.Error().Err(err).Msgf("error getting edge from %s to %s", from, to)
			return false
		}
		log.Debug().
			Str("from", w.NodeLookup[from].NormalizedName()).
			Str("to", w.NodeLookup[to].NormalizedName()).
			Str("edge", e.ID()).
			Msg("workflow forward bfs")
		_ = subgraph.AddVertex(from)
		_ = subgraph.AddVertex(to)
		_ = subgraph.AddEdge(from, to)
		discovered = append(discovered, to)
		return false
	})
	if err != nil {
		return nil, errors.Wrapf(err, "error traversing graph forward")
	}
	for _, discID := range discovered {
		err = util.BFS(w.Graph, discID, false, func(to, from string) bool {
			e, err := getEdge(from, to)
			if err != nil {
				log.Error().Err(err).Msgf("error getting edge from %s to %s", from, to)
				return false
			}
			log.Debug().
				Str("from", w.NodeLookup[from].NormalizedName()).
				Str("to", w.NodeLookup[to].NormalizedName()).
				Str("edge", e.ID()).
				Msg("workflow backward bfs")
			_ = subgraph.AddVertex(from)
			_ = subgraph.AddVertex(to)
			_ = subgraph.AddEdge(from, to)
			return false
		})
		if err != nil {
			return nil, errors.Wrapf(err, "error traversing graph backward")
		}
	}

	depList, err := graphlib.TopologicalSort(subgraph)
	if err != nil {
		return nil, errors.Wrapf(err, "error sorting graph")
	}

	var deps []string
	for _, nID := range depList {
		n := w.NodeLookup[nID]
		deps = append(deps, n.NormalizedName())
	}
	log.Debug().Msgf("sorted nodes: %s", strings.Join(deps, ", "))

	ctx, connector := manager.Start(ctx)

	// since the graph is sorted topologically, we can wire the nodes in order
	for _, nID := range depList {
		n := w.NodeLookup[nID]
		// TODO breadchris how do we handle lifecycle of nodes? cleaning up resources,
		err = w.wireWorkflow(ctx, connector, n)
		if err != nil {
			log.Error().Err(err).Msgf("failed to traverse workflow")
			return nil, err
		}
	}
	obs := connector.Connect(ctx)
	obs.ForEach(func(item any) {}, func(err error) {
		e := manager.Stop()
		if e != nil {
			log.Error().Err(e).Msgf("error, manager is stopping workflow")
		}
	}, func() {
		log.Debug().Msgf("completed, manager is stopping workflow")
		err = manager.Stop()
		if err != nil {
			log.Error().Err(err).Msgf("error stopping manager")
		}
	})

	// TODO breadchris this is returning a trace, not the actual output
	return &graph.IO{
		Observable: connector.Connect(ctx),
	}, nil
}

func (w *Workflow) wireWorkflow(
	ctx context.Context,
	connector *Connector,
	node graph.Node,
) error {
	// TODO breadchris should be able to skip this if there is only one publisher
	var inputObs []rxgo.Observable
	for _, publisher := range node.Publishers() {
		io, ok := connector.Get(publisher.GetNode().ID())
		if !ok {
			return fmt.Errorf("publisher %s not found", publisher.GetNode().ID())
		}

		var err error
		io, err = publisher.Transform(ctx, io)
		if err != nil {
			return errors.Wrapf(err, "error transforming publisher %s", publisher.GetNode().ID())
		}
		inputObs = append(inputObs, io.Observable)
	}

	// TODO breadchris probably want to know what nodes are being combined
	i := rxgo.CombineLatest(func(i ...any) any {
		combined := map[string]any{}
		for _, j := range i {
			switch t := j.(type) {
			case map[string]any:
				combined = util.Merge(combined, t)
			default:
				// TODO breadchris how can inputs be merged
				log.Warn().Msgf("unexpected type %T when merging inputs", t)
				return j
				//b, err := json.Marshal(t)
				//if err != nil {
				//	log.Error().Err(err).Msgf("error marshalling input %T", t)
				//	continue
				//}
				//var m map[string]any
				//err = json.Unmarshal(b, &m)
				//if err != nil {
				//	log.Error().Err(err).Msgf("error unmarshalling input %T", t)
				//	continue
				//}
				//combined = util.Merge(combined, m)
			}
		}
		return combined
	}, inputObs)
	input := graph.IO{
		Observable: i,
	}

	log.Debug().
		Str("node", node.NormalizedName()).
		Msg("wiring node IO")
	output, err := node.Wire(ctx, input)
	if err != nil {
		return errors.Wrapf(err, "error executing node: %s", node.NormalizedName())
	}

	if output.Observable == nil {
		log.Warn().Msgf("node %s returned nil observable", node.NormalizedName())
		return nil
	}

	connector.Add(node.ID(), &output)
	return nil
}
