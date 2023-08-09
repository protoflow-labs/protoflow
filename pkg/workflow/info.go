package workflow

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/data"
)

// TODO breadchris separate "infer" and "collect" type information
// TODO breadchris node.Type should be passed workflow since a node's info can depend on other nodes
func (w *Workflow) GetNodeInfo(n graph.Node) (*graph.Info, error) {
	var resp *graph.Info
	switch n.(type) {
	case *data.InputNode:
		listeners := n.Subscribers()
		if len(listeners) != 1 {
			// TODO breadchris support multiple listeners
			return nil, errors.Errorf("input node should have 1 child, got %d", len(listeners))
		}
		// TODO breadchris only designed for 1 child
		for _, l := range listeners {
			child := l.GetNode()
			return w.GetNodeInfo(child)
		}
	//case *code.FunctionNode:
	//	var (
	//		childInputs     []protoreflect.MessageDescriptor
	//		parentOutputs   []protoreflect.MessageDescriptor
	//		streamingChild  bool
	//		streamingParent bool
	//	)
	//
	//	// TODO breadchris if there is no publisher to determine the input type, try to see if the type already exists
	//	if len(n.Publishers()) == 0 {
	//		return n.Type()
	//	}
	//
	//	// TODO breadchris if two function nodes are connected, you can't infer the type
	//	// make sure an infinite loop doesn't happen
	//	for _, listener := range n.Subscribers() {
	//		child := listener.GetNode()
	//		switch n.(type) {
	//		case *code.FunctionNode:
	//			log.Warn().
	//				Str("parent", n.NormalizedName()).
	//				Str("child", child.NormalizedName()).
	//				Msg("function node connected to function node not supported yet")
	//			continue
	//		}
	//
	//		childType, err := w.GetNodeInfo(n)
	//		if err != nil {
	//			return nil, err
	//		}
	//		if childType == nil {
	//			return nil, errors.Wrapf(err, "error getting node info for %s", n.NormalizedName())
	//		}
	//		if childType.Method.MethodDesc.IsStreamingClient() {
	//			streamingChild = true
	//		}
	//		childInputs = append(childInputs, childType.Method.MethodDesc.Input())
	//	}
	//	for _, parent := range n.Publishers() {
	//		switch n.(type) {
	//		case *code.FunctionNode:
	//			log.Warn().
	//				Str("parent", parent.NormalizedName()).
	//				Str("child", n.NormalizedName()).
	//				Msg("function node connected to function node not supported yet")
	//			continue
	//		}
	//
	//		parentType, err := w.GetNodeInfo(n)
	//		if err != nil {
	//			return nil, err
	//		}
	//		if parentType == nil {
	//			return nil, errors.Wrapf(err, "error getting node info for %s", n.NormalizedName())
	//		}
	//		if parentType.Method.MethodDesc.IsStreamingServer() {
	//			streamingParent = true
	//		}
	//		parentOutputs = append(parentOutputs, parentType.Method.MethodDesc.Output())
	//	}
	//	inputType, err := graph.MessageFromTypes(n.NormalizedName()+"Request", parentOutputs)
	//	if err != nil {
	//		return nil, errors.Wrapf(err, "error building request message for %s", n.NormalizedName())
	//	}
	//	outputType, err := graph.MessageFromTypes(n.NormalizedName()+"Response", childInputs)
	//	if err != nil {
	//		return nil, errors.Wrapf(err, "error building response message for %s", n.NormalizedName())
	//	}
	//
	//	req := builder.RpcTypeImportedMessage(inputType, streamingParent)
	//	res := builder.RpcTypeImportedMessage(outputType, streamingChild)
	//
	//	// TODO breadchris this is a hack to get the name of the function
	//	s := builder.NewService("Service")
	//	b := builder.NewMethod(n.NormalizedName(), req, res)
	//	s.AddMethod(b)
	//
	//	m, err := b.Build()
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	mthd, err := grpc.NewMethodDescriptor(m.UnwrapMethod())
	//	if err != nil {
	//		return nil, err
	//	}
	//	return &graph.Info{
	//		Method: mthd,
	//	}, nil
	//case *worknode.BucketNode:
	//	// TODO breadchris how do you handle file permissions?
	//	reqMsg := builder.NewMessage("Request")
	//	reqMsg = reqMsg.AddField(builder.NewField("path", builder.FieldTypeString()))
	//	reqMsg = reqMsg.AddField(builder.NewField("data", builder.FieldTypeBytes()))
	//	req := builder.RpcTypeMessage(reqMsg, true)
	//
	//	resMsg := builder.NewMessage("Response")
	//	resMsg = resMsg.AddField(builder.NewField("path", builder.FieldTypeString()))
	//	// TODO breadchris what does this type mean if it streaming or not? sync vs async?
	//	res := builder.RpcTypeMessage(resMsg, false)
	//
	//	s := builder.NewService("Service")
	//	b := builder.NewMethod(n.NormalizedName(), req, res)
	//	s.AddMethod(b)
	//
	//	m, err := b.Build()
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	mthd, err := grpc.NewMethodDescriptor(m.UnwrapMethod())
	//	if err != nil {
	//		return nil, err
	//	}
	//	resp = &graph.Type{
	//		Method: mthd,
	//	}
	//case *worknode.FileNode:
	//	// TODO breadchris how do you handle file permissions?
	//	reqMsg := builder.NewMessage("Request")
	//	req := builder.RpcTypeMessage(reqMsg, false)
	//
	//	resMsg := builder.NewMessage("Response")
	//	resMsg = resMsg.AddField(builder.NewField("path", builder.FieldTypeString()))
	//	// TODO breadchris what does this type mean if it streaming or not? sync vs async?
	//	res := builder.RpcTypeMessage(resMsg, false)
	//
	//	s := builder.NewService("Service")
	//	b := builder.NewMethod(n.NormalizedName(), req, res)
	//	s.AddMethod(b)
	//
	//	m, err := b.Build()
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	mthd, err := grpc.NewMethodDescriptor(m.UnwrapMethod())
	//	if err != nil {
	//		return nil, err
	//	}
	//	resp = &graph.Type{
	//		Method: mthd,
	//	}
	default:
		return n.Type()
	}
	return resp, nil
}
