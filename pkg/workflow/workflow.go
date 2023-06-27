package workflow

import (
	"encoding/json"
	"fmt"
	"github.com/dominikbraun/graph"
	"github.com/google/uuid"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/workflow/execute"
	worknode "github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type AdjMap map[string]map[string]graph.Edge[string]

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

// TODO breadchris if there is only one field, set name of message to just the name of the one field.
// this is canonical in grpc
func messageFromTypes(name string, types []protoreflect.MessageDescriptor) (*desc.MessageDescriptor, error) {
	if len(types) == 1 {
		return desc.WrapMessage(types[0])
	}
	mb := builder.NewMessage(name)
	if len(types) == 0 {
		return mb.Build()
	}
	var addedFields []string
	for _, t := range types {
		wt, err := desc.WrapMessage(t)
		if err != nil {
			return nil, errors.Wrapf(err, "error wrapping message %s", name)
		}
		msgBuilder, err := builder.FromMessage(wt)
		if err != nil {
			return nil, errors.Wrapf(err, "error building message %s", name)
		}
		fm := builder.FieldTypeMessage(msgBuilder)

		fieldName := string(t.Name())
		if lo.Contains(addedFields, fieldName) {
			return nil, errors.Errorf("duplicate field %s", name)
		}

		mb = mb.AddField(builder.NewField(fieldName, fm))
	}
	return mb.Build()
}

// TODO breadchris separate "infer" and "collect" type information
func (w *Workflow) GetNodeInfo(n worknode.Node) (*worknode.Info, error) {
	var resp *worknode.Info
	switch n.(type) {
	case *worknode.InputNode:
		children := w.AdjMap[n.ID()]
		if len(children) != 1 {
			// TODO breadchris support multiple children
			return nil, errors.Errorf("input node should have 1 child, got %d", len(children))
		}
		// TODO breadchris only designed for 1 child
		for child := range children {
			n, err := w.GetNode(child)
			if err != nil {
				return nil, errors.Errorf("node %s not found", child)
			}
			return w.GetNodeInfo(n)
		}
	case *worknode.FunctionNode:
		children := w.AdjMap[n.ID()]
		parents := w.PreMap[n.ID()]

		var (
			childInputs   []protoreflect.MessageDescriptor
			parentOutputs []protoreflect.MessageDescriptor
		)

		for child := range children {
			n, err := w.GetNode(child)
			if err != nil {
				return nil, errors.Errorf("node %s not found", child)
			}
			childType, err := w.GetNodeInfo(n)
			if err != nil {
				return nil, err
			}
			childInputs = append(childInputs, childType.Method.MethodDesc.Input())
		}
		for parent := range parents {
			n, err := w.GetNode(parent)
			if err != nil {
				return nil, errors.Errorf("node %s not found", parent)
			}
			parentType, err := w.GetNodeInfo(n)
			if err != nil {
				return nil, err
			}
			parentOutputs = append(parentOutputs, parentType.Method.MethodDesc.Output())
		}
		intputType, err := messageFromTypes(n.NormalizedName()+"Request", parentOutputs)
		if err != nil {
			return nil, errors.Wrapf(err, "error building request message for %s", n.NormalizedName())
		}
		outputType, err := messageFromTypes(n.NormalizedName()+"Response", childInputs)
		if err != nil {
			return nil, errors.Wrapf(err, "error building response message for %s", n.NormalizedName())
		}

		// TODO breadchris how can we determine if the req/res are streaming?
		req := builder.RpcTypeImportedMessage(intputType, false)
		res := builder.RpcTypeImportedMessage(outputType, false)

		// TODO breadchris this is a hack to get the name of the function
		s := builder.NewService("Service")
		b := builder.NewMethod(n.NormalizedName(), req, res)
		s.AddMethod(b)

		m, err := b.Build()
		if err != nil {
			return nil, err
		}

		mthd, err := grpc.NewMethodDescriptor(m.UnwrapMethod())
		if err != nil {
			return nil, err
		}
		return &worknode.Info{
			Method: mthd,
		}, nil
	default:
		res, err := w.GetNodeResource(n.ID())
		if err != nil {
			return nil, err
		}
		return res.Info(n)
	}
	return resp, nil
}

type ResourceMap map[string]resource.Resource

func FromProject(project *gen.Project) (*Workflow, error) {
	g := graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())

	if project.Graph == nil {
		return nil, errors.New("project graph is nil")
	}

	resources := ResourceMap{}
	for _, protoRes := range project.Resources {
		r, err := resource.FromProto(protoRes)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating resource for node %s", protoRes.Id)
		}
		resources[protoRes.Id] = r
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
			log.Warn().Msgf("no resource found for node %s", node.Id)
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

// TODO breadchris can this be a map[string]Resource?
type Instances map[string]resource.Resource

// TODO breadchris nodeID should not be needed, the workflow should already be a slice of the graph that is configured to run
func (w *Workflow) Run(logger Logger, executor execute.Executor, nodeID string, input interface{}) (*execute.Result, error) {
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

	res, err := w.traverseWorkflow(logger, instances, executor, vert, execute.Input{
		Params: input,
	})
	if err != nil {
		logger.Error("failed to traverse workflow", "error", err)
		return nil, err
	}
	return res, nil
}

func (w *Workflow) traverseWorkflow(logger Logger, instances Instances, executor execute.Executor, vert string, input execute.Input) (*execute.Result, error) {
	node, ok := w.NodeLookup[vert]
	if !ok {
		return nil, fmt.Errorf("vertex not found: %s", vert)
	}

	log.Debug().Str("vert", vert).Msg("injecting dependencies for node")
	err := injectDepsForNode(instances, &input, node)
	if err != nil {
		return nil, errors.Wrapf(err, "error injecting dependencies for node %s", vert)
	}

	log.Debug().Str("vert", vert).Interface("resource", input.Resource).Msg("executing node")
	res, err := executor.Execute(node, input)
	if err != nil {
		return nil, errors.Wrapf(err, "error executing node: %s", vert)
	}

	var (
		nextBlockInput execute.Input
	)
	if res.Stream != nil {
		// if the result is a stream, pass the stream to the next block
		nextBlockInput = execute.Input{
			Stream: res.Stream,
		}
	} else {
		// TODO breadchris have this work for streams as well
		traceNodeExec(executor, vert, input, res.Data)

		// otherwise pass the singular result data to the next block
		nextBlockInput = execute.Input{
			Params: res.Data,
		}
	}

	nextResSet := false
	for neighbor := range w.AdjMap[vert] {
		logger.Info("traversing workflow", "nodeID", neighbor)

		// TODO breadchris if there are multiple neighbors, and there is a stream, the stream should be split and passed to each neighbor

		neighborRes, err := w.traverseWorkflow(logger, instances, executor, neighbor, nextBlockInput)
		if err != nil {
			// TODO breadchris if we are inside of a stream, the error should bubble up to the stream entrypoint
			return nil, errors.Wrapf(err, "error traversing workflow %s", neighbor)
		}

		// TODO breadchris how should multiple results be handled?
		if !nextResSet {
			res = neighborRes
			nextResSet = true
		}
	}
	// TODO breadchris need to implement pubsub to handle lifecycle of streams
	//if input.Stream != nil {
	//	for {
	//		time.Sleep(10 * time.Second)
	//		log.Info().Msg("waiting for stream to finish")
	//	}
	//}
	return &execute.Result{
		Data: res.Data,
	}, nil
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

func traceNodeExec(executor execute.Executor, nodeID string, input any, output any) {
	// TODO breadchris clean this up
	inputSer, inputErr := json.Marshal(input)
	outputSer, outputErr := json.Marshal(output)
	if inputErr != nil || outputErr != nil {
		log.Error().
			Err(inputErr).
			Err(outputErr).
			Msg("error serializing node execution")
	}
	err := executor.Trace(&gen.NodeExecution{
		NodeId: nodeID,
		Input:  string(inputSer),
		Output: string(outputSer),
	})
	if err != nil {
		log.Error().Err(err).Msg("error tracing node execution")
	}
}
