package workflow

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
)

type Node interface {
	Execute(executor Executor, input Input) (*Result, error)
}

type BaseNode struct {
	ResourceIDs []string
}

type GRPCNode struct {
	BaseNode
	*gen.GRPC
}

var _ Node = &GRPCNode{}

type RESTNode struct {
	BaseNode
	*gen.REST
}

var _ Node = &RESTNode{}

type CollectionNode struct {
	BaseNode
	*gen.Collection
}

var _ Node = &CollectionNode{}

type BucketNode struct {
	BaseNode
	*gen.Bucket
}

var _ Node = &BucketNode{}

type InputNode struct {
	BaseNode
	*gen.Input
}

var _ Node = &InputNode{}

var activity = &Activity{}

func (s *GRPCNode) Execute(executor Executor, input Input) (*Result, error) {
	return executor.Execute(activity.ExecuteGRPCNode, s, input)
}

func (s *RESTNode) Execute(executor Executor, input Input) (*Result, error) {
	return executor.Execute(activity.ExecuteRestNode, s, input)
}

func (s *CollectionNode) Execute(executor Executor, input Input) (*Result, error) {
	docs, err := getResource[DocstoreResource](input.Resources)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting docstore resource")
	}

	d, cleanup, err := docs.WithKeyField(s.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting to collection")
	}
	defer cleanup()

	err = d.Create(context.Background(), input.Params)
	return &Result{
		Data: input.Params,
	}, nil
}

func (s *BucketNode) Execute(executor Executor, input Input) (*Result, error) {
	bucket, err := getResource[BlobstoreResource](input.Resources)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting bucket resource")
	}

	var bucketData []byte
	switch input.Params.(type) {
	case []byte:
		bucketData = input.Params.([]byte)
	case string:
		bucketData = []byte(input.Params.(string))
	default:
		bucketData, err = json.Marshal(input.Params)
		if err != nil {
			return nil, errors.Wrapf(err, "error marshaling input params")
		}
	}

	b, cleanup, err := bucket.WithPath(s.Path)
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting to bucket")
	}
	defer cleanup()

	err = b.WriteAll(context.Background(), s.Path, bucketData, nil)
	return &Result{
		Data: input.Params,
	}, nil
}

func (s *InputNode) Execute(executor Executor, input Input) (*Result, error) {
	return &Result{
		Data: input.Params,
	}, nil
}

func NewNode(node *gen.Node, block *gen.Block) (Node, error) {
	switch block.Type.(type) {
	case *gen.Block_Grpc:
		return NewGRPCNode(node, block), nil
	case *gen.Block_Collection:
		return NewCollectionNode(node, block), nil
	case *gen.Block_Bucket:
		return NewBucketNode(node, block), nil
	case *gen.Block_Rest:
		return NewRestNode(node, block), nil
	case *gen.Block_Input:
		return NewInputNode(node, block), nil
	default:
		return nil, errors.New("no node found")
	}
}

func NewGRPCNode(node *gen.Node, block *gen.Block) *GRPCNode {
	b := block.GetGrpc()
	n := node.GetGrpc()
	if n != nil {
		if n.Service != "" {
			b.Service = n.Service
		}
		if n.Method != "" {
			b.Method = n.Method
		}
	}
	return &GRPCNode{
		BaseNode: BaseNode{
			ResourceIDs: node.ResourceIds,
		},
		GRPC: b,
	}
}

func NewRestNode(node *gen.Node, block *gen.Block) *RESTNode {
	b := block.GetRest()
	n := node.GetRest()
	if n != nil {
		if n.Path != "" {
			b.Path = n.Path
		}
		if n.Method != "" {
			b.Method = n.Method
		}
	}
	return &RESTNode{
		BaseNode: BaseNode{
			ResourceIDs: node.ResourceIds,
		},
		REST: b,
	}
}

func NewCollectionNode(node *gen.Node, block *gen.Block) *CollectionNode {
	b := block.GetCollection()
	n := node.GetCollection()
	if n != nil {
		if n.Name != "" {
			b.Name = n.Name
		}
	}
	return &CollectionNode{
		BaseNode: BaseNode{
			ResourceIDs: node.ResourceIds,
		},
		Collection: b,
	}
}

func NewBucketNode(node *gen.Node, block *gen.Block) *BucketNode {
	b := block.GetBucket()
	n := node.GetBucket()
	if n != nil {
		if n.Path != "" {
			b.Path = n.Path
		}
	}
	return &BucketNode{
		BaseNode: BaseNode{
			ResourceIDs: node.ResourceIds,
		},
		Bucket: b,
	}
}

func NewInputNode(node *gen.Node, block *gen.Block) *InputNode {
	b := block.GetInput()
	n := node.GetInput()
	if n != nil {
		if n.Fields != nil {
			b.Fields = n.Fields
		}
	}
	return &InputNode{
		BaseNode: BaseNode{
			ResourceIDs: node.ResourceIds,
		},
		Input: b,
	}
}
