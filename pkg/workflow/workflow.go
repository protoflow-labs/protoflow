package workflow

import (
	"context"
	"fmt"
	"github.com/dominikbraun/graph"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/workflow/execute"
	worknode "github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/protobuf/reflect/protoreflect"
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
	case *worknode.BucketNode:
		// TODO breadchris how do you handle file permissions?
		reqMsg := builder.NewMessage("Request")
		reqMsg = reqMsg.AddField(builder.NewField("path", builder.FieldTypeString()))
		reqMsg = reqMsg.AddField(builder.NewField("data", builder.FieldTypeBytes()))
		req := builder.RpcTypeMessage(reqMsg, true)

		resMsg := builder.NewMessage("Response")
		resMsg = resMsg.AddField(builder.NewField("path", builder.FieldTypeString()))
		// TODO breadchris what does this type mean if it streaming or not? sync vs async?
		res := builder.RpcTypeMessage(resMsg, false)

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
		resp = &worknode.Info{
			Method: mthd,
		}
	case *worknode.PromptNode:
		reqMsg := builder.NewMessage("Request")
		reqMsg = reqMsg.AddField(builder.NewField("message", builder.FieldTypeString()))
		req := builder.RpcTypeMessage(reqMsg, true)

		resMsg := builder.NewMessage("Response")
		resMsg = resMsg.AddField(builder.NewField("result", builder.FieldTypeString()))
		res := builder.RpcTypeMessage(resMsg, false)

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
		resp = &worknode.Info{
			Method: mthd,
		}
	case *worknode.FunctionNode:
		children := w.AdjMap[n.ID()]
		parents := w.PreMap[n.ID()]

		var (
			childInputs     []protoreflect.MessageDescriptor
			parentOutputs   []protoreflect.MessageDescriptor
			streamingChild  bool
			streamingParent bool
		)

		// TODO breadchris if two function nodes are connected, you can't infer the type
		// make sure an infinite loop doesn't happen
		for child := range children {
			n, err := w.GetNode(child)
			if err != nil {
				return nil, errors.Errorf("node %s not found", child)
			}

			switch n.(type) {
			case *worknode.FunctionNode:
				log.Warn().
					Str("parent", n.ID()).
					Str("child", child).
					Msg("function node connected to function node not supported yet")
				continue
			}

			childType, err := w.GetNodeInfo(n)
			if err != nil {
				return nil, err
			}
			if childType == nil {
				return nil, errors.Wrapf(err, "error getting node info for %s", n.NormalizedName())
			}
			if childType.Method.MethodDesc.IsStreamingClient() {
				streamingChild = true
			}
			childInputs = append(childInputs, childType.Method.MethodDesc.Input())
		}
		for parent := range parents {
			n, err := w.GetNode(parent)
			if err != nil {
				return nil, errors.Errorf("node %s not found", parent)
			}

			switch n.(type) {
			case *worknode.FunctionNode:
				log.Warn().
					Str("parent", parent).
					Str("child", n.ID()).
					Msg("function node connected to function node not supported yet")
				continue
			}

			parentType, err := w.GetNodeInfo(n)
			if err != nil {
				return nil, err
			}
			if parentType == nil {
				return nil, errors.Wrapf(err, "error getting node info for %s", n.NormalizedName())
			}
			if parentType.Method.MethodDesc.IsStreamingServer() {
				streamingParent = true
			}
			parentOutputs = append(parentOutputs, parentType.Method.MethodDesc.Output())
		}
		inputType, err := messageFromTypes(n.NormalizedName()+"Request", parentOutputs)
		if err != nil {
			return nil, errors.Wrapf(err, "error building request message for %s", n.NormalizedName())
		}
		outputType, err := messageFromTypes(n.NormalizedName()+"Response", childInputs)
		if err != nil {
			return nil, errors.Wrapf(err, "error building response message for %s", n.NormalizedName())
		}

		req := builder.RpcTypeImportedMessage(inputType, streamingParent)
		res := builder.RpcTypeImportedMessage(outputType, streamingChild)

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

type Connector struct {
	observers []rxgo.Observable
}

func NewConnector() *Connector {
	return &Connector{
		observers: []rxgo.Observable{},
	}
}

func (c *Connector) Add(o rxgo.Observable) {
	c.observers = append(c.observers, o)
}

func (c *Connector) Connect(ctx context.Context) rxgo.Observable {
	// TODO breadchris is this publish startegy needed here
	o := rxgo.Merge(c.observers, rxgo.WithPublishStrategy())
	// TODO breadchris figure out what to do with disposed and cancel
	// disposed, cancel := output.Observable.Connect(ctx)
	for _, obs := range c.observers {
		obs.Connect(ctx)
	}
	o.Connect(ctx)
	return o
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

	// TODO breadchris figure out what to do with disposed and cancel
	// disposed, cancel := output.Observable.Connect(ctx)

	// rx.LogObserver(node.NormalizedName(), output.Observable)
	// output.Observable.Connect(ctx)

	nextBlockInput := execute.Input{
		Observable: output.Observable,
	}

	for neighbor := range w.AdjMap[nodeID] {
		log.Debug().
			Str("node", node.NormalizedName()).
			Str("neighbor", w.NodeLookup[neighbor].NormalizedName()).
			Msg("traversing workflow")

		// TODO breadchris if there are multiple neighbors, and there is a stream, the stream should be split and passed to each neighbor
		// TODO breadchris how should we handle multiple neighbors? output is being overwritten
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
