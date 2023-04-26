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

type RESTNode struct {
	BaseNode
	*gen.REST
}

type CollectionNode struct {
	BaseNode
	*gen.Collection
}

type BucketNode struct {
	BaseNode
	*gen.Bucket
}

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
	err = docs.Collection.Create(context.Background(), input.Params)
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
	err = bucket.Bucket.WriteAll(context.Background(), s.Path, bucketData, nil)
	return &Result{
		Data: input.Params,
	}, nil
}

func NewNode(node *gen.Node, block *gen.Block) (Node, error) {
	switch block.Type.(type) {
	case *gen.Block_Grpc:
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
			GRPC: b,
		}, nil
	default:
		return nil, errors.New("no node found")
	}
}
