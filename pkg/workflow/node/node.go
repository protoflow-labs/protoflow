package node

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
)

// TODO breadchris make this something that can be modularized
func NewNode(node *gen.Node) (graph.Node, error) {
	switch node.Config.(type) {
	case *gen.Node_Grpc:
		return NewGRPCNode(node), nil
	case *gen.Node_Collection:
		return NewCollectionNode(node), nil
	case *gen.Node_Bucket:
		return NewBucketNode(node), nil
	case *gen.Node_Rest:
		return NewRestNode(node), nil
	case *gen.Node_Input:
		return NewInputNode(node), nil
	case *gen.Node_Function:
		return NewFunctionNode(node), nil
	case *gen.Node_Query:
		return NewQueryNode(node), nil
	case *gen.Node_Prompt:
		return NewPromptNode(node), nil
	case *gen.Node_Configuration:
		return NewConfigNode(node), nil
	case *gen.Node_Secret:
		return NewSecretNode(node), nil
	case *gen.Node_Template:
		return NewTemplateNode(node), nil
	case *gen.Node_Route:
		return NewRouteNode(node), nil
	case *gen.Node_File:
		return NewFileNode(node), nil
	default:
		return nil, errors.New("no node found")
	}
}

// NewBaseNode creates a new BaseNode from a gen.Node, gen.Node cannot be embedded into BaseNode because proto deserialization will fail on the type
func NewBaseNode(node *gen.Node) BaseNode {
	return BaseNode{
		Name:       util.ToTitleCase(node.Name),
		id:         node.Id,
		resourceID: node.ResourceId,
	}
}

func NewPromptNode(node *gen.Node) *PromptNode {
	return &PromptNode{
		BaseNode: NewBaseNode(node),
		Prompt:   node.GetPrompt(),
	}
}

func NewSecretNode(node *gen.Node) *SecretNode {
	return &SecretNode{
		BaseNode: NewBaseNode(node),
		Secret:   node.GetSecret(),
	}
}

func NewTemplateNode(node *gen.Node) *TemplateNode {
	return &TemplateNode{
		BaseNode: NewBaseNode(node),
		Template: node.GetTemplate(),
	}
}

func NewFileNode(node *gen.Node) *FileNode {
	return &FileNode{
		BaseNode: NewBaseNode(node),
		File:     node.GetFile(),
	}
}

func NewRestNode(node *gen.Node) *RESTNode {
	return &RESTNode{
		BaseNode: NewBaseNode(node),
		REST:     node.GetRest(),
	}
}

func NewCollectionNode(node *gen.Node) *CollectionNode {
	return &CollectionNode{
		BaseNode:   NewBaseNode(node),
		Collection: node.GetCollection(),
	}
}

func NewBucketNode(node *gen.Node) *BucketNode {
	return &BucketNode{
		BaseNode: NewBaseNode(node),
		Bucket:   node.GetBucket(),
	}
}

func NewQueryNode(node *gen.Node) *QueryNode {
	return &QueryNode{
		BaseNode: NewBaseNode(node),
		Query:    node.GetQuery(),
	}
}
