package workflow

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/node/code"
	"github.com/protoflow-labs/protoflow/pkg/node/data"
	"github.com/protoflow-labs/protoflow/pkg/node/reason"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// TODO breadchris separate "infer" and "collect" type information
func (w *Workflow) GetNodeInfo(n graph.Node) (*graph.Info, error) {
	var resp *graph.Info
	switch n.(type) {
	case *data.InputNode:
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
	case *reason.PromptNode:
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
		resp = &graph.Info{
			Method: mthd,
		}
	case *code.FunctionNode:
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
			case *code.FunctionNode:
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
			case *code.FunctionNode:
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
		return &graph.Info{
			Method: mthd,
		}, nil
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
	//	resp = &graph.Info{
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
	//	resp = &graph.Info{
	//		Method: mthd,
	//	}
	default:
		return n.Info()
	}
	return resp, nil
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
