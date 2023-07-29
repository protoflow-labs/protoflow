package node

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
)

// TODO breadchris make this something that can be modularized
func NewNode(node *gen.Node) (graph.Node, error) {
	if node.Id == "" {
		node.Id = uuid.NewString()
	}
	switch node.Config.(type) {
	case *gen.Node_Grpc:
		return NewGRPCNode(node), nil
	case *gen.Node_Input:
		return NewInputNode(node), nil
	case *gen.Node_Function:
		return NewFunctionNode(node), nil
	case *gen.Node_Prompt:
		return NewPromptNode(node), nil
	case *gen.Node_Template:
		return NewTemplateNode(node), nil
	case *gen.Node_Route:
		return NewRouteNode(node), nil
	//case *gen.Node_Collection:
	//	return NewCollectionNode(node), nil
	//case *gen.Node_Bucket:
	//	return NewBucketNode(node), nil
	//case *gen.Node_Query:
	//	return NewQueryNode(node), nil
	//case *gen.Node_Configuration:
	//	return NewConfigNode(node), nil
	//case *gen.Node_File:
	//	return NewFileNode(node), nil
	default:
		return nil, errors.New("no node found")
	}
}
