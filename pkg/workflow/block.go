package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"io"
	"strings"
)

type Input struct {
	Params       interface{}
	Resources    map[string]any
	Dependencies []string
	Stream       bufcurl.OutputStream
}

type Result struct {
	Data   interface{}
	Stream bufcurl.OutputStream
}

type Node interface {
	Execute(executor Executor, input Input) (*Result, error)
	Dependencies() []string
	NormalizedName() string
}

type BaseNode struct {
	Name        string
	ResourceIDs []string
}

func (n *BaseNode) Dependencies() []string {
	return n.ResourceIDs
}

func (n *BaseNode) NormalizedName() string {
	name := util.ToTitleCase(n.Name)
	if strings.Contains(name, ".") {
		name = strings.Split(name, ".")[1]
	}
	return name
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

type FunctionNode struct {
	BaseNode
	Function *gen.Function
}

var _ Node = &FunctionNode{}

type QueryNode struct {
	BaseNode
	Query *gen.Query
}

var _ Node = &FunctionNode{}

var activity = &Activity{}

func (s *GRPCNode) Execute(executor Executor, input Input) (*Result, error) {
	return executor.Execute(activity.ExecuteGRPCNode, s, input)
}

func (s *RESTNode) Execute(executor Executor, input Input) (*Result, error) {
	return executor.Execute(activity.ExecuteRestNode, s, input)
}

func (s *CollectionNode) Execute(executor Executor, input Input) (*Result, error) {
	docs, ok := input.Resources[DocstoreResourceType].(*DocstoreResource)
	if !ok {
		return nil, fmt.Errorf("error getting docstore resource: %s", s.Collection.Name)
	}

	collection, cleanup, err := docs.WithCollection(s.Collection.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting to collection")
	}
	defer cleanup()

	var records []map[string]interface{}

	switch input := input.Params.(type) {
	case map[string]interface{}:
		records = append(records, input)
	case []*map[string]interface{}:
		for _, record := range input {
			records = append(records, *record)
		}
	default:
		return nil, fmt.Errorf("error unsupported input type: %T", input)
	}

	for _, record := range records {
		if record["id"] == nil {
			record["id"] = uuid.NewString()
		}
		err = collection.Create(context.Background(), record)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating doc")
		}
	}

	return &Result{
		Data: input.Params,
	}, nil
}

func (s *BucketNode) Execute(executor Executor, input Input) (*Result, error) {
	bucket, ok := input.Resources[BlobstoreResourceType].(*BlobstoreResource)
	if !ok {
		return nil, fmt.Errorf("error getting blobstore resource: %s", s.Bucket.Path)
	}

	var (
		err        error
		bucketData []byte
	)
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
		return nil, errors.Wrapf(err, fmt.Sprintf("error connecting to bucket: %s", s.Path))
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

func (f *FunctionNode) Execute(executor Executor, input Input) (*Result, error) {
	return executor.Execute(activity.ExecuteFunctionNode, f, input)
}

func (s *QueryNode) Execute(executor Executor, input Input) (*Result, error) {
	docResource, ok := input.Resources[DocstoreResourceType].(*DocstoreResource)
	if !ok {
		return nil, fmt.Errorf("error getting docstore resource: %s", s.Query.Collection)
	}

	d, cleanup, err := docResource.WithCollection(s.Query.Collection)
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting to collection")
	}
	defer cleanup()

	var docs []*map[string]interface{}
	iter := d.Query().Get(context.Background())
	for {
		doc := map[string]interface{}{}
		err = iter.Next(context.Background(), doc)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Wrapf(err, "error iterating over query results")
		}
		docs = append(docs, &doc)
	}
	return &Result{
		Data: docs,
	}, nil
}

type ResourceMap map[string]Resource

func NewNode(resources ResourceMap, node *gen.Node) (Node, error) {
	switch node.Config.(type) {
	case *gen.Node_Grpc:
		return NewGRPCNode(node), nil
	case *gen.Node_Collection:
		return NewCollectionNode(resources, node), nil
	case *gen.Node_Bucket:
		return NewBucketNode(resources, node), nil
	case *gen.Node_Rest:
		return NewRestNode(node), nil
	case *gen.Node_Input:
		return NewInputNode(node), nil
	case *gen.Node_Function:
		return NewFunctionNode(node), nil
	case *gen.Node_Query:
		return NewQueryNode(resources, node), nil
	default:
		return nil, errors.New("no node found")
	}
}

// NewBaseNode creates a new BaseNode from a gen.Node, gen.Node cannot be embedded into BaseNode because proto deserialization will fail on the type
func NewBaseNode(node *gen.Node) BaseNode {
	return BaseNode{
		Name:        util.ToTitleCase(node.Name),
		ResourceIDs: node.ResourceIds,
	}
}

// TODO breadchris we are ignoring blocks that are not set on the node, nodes should have blocks
func NewGRPCNode(node *gen.Node) *GRPCNode {
	return &GRPCNode{
		BaseNode: NewBaseNode(node),
		GRPC:     node.GetGrpc(),
	}
}

func NewRestNode(node *gen.Node) *RESTNode {
	return &RESTNode{
		BaseNode: NewBaseNode(node),
		REST:     node.GetRest(),
	}
}

func NewCollectionNode(resources ResourceMap, node *gen.Node) *CollectionNode {
	// TODO breadchris a node should already have this resource configured
	for id, r := range resources {
		if r.Name() == DocstoreResourceType {
			node.ResourceIds = append(node.ResourceIds, id)
		}
	}
	return &CollectionNode{
		BaseNode:   NewBaseNode(node),
		Collection: node.GetCollection(),
	}
}

func NewBucketNode(resources ResourceMap, node *gen.Node) *BucketNode {
	// TODO breadchris a node should already have this resource configured
	for id, r := range resources {
		if r.Name() == BlobstoreResourceType {
			node.ResourceIds = append(node.ResourceIds, id)
		}
	}
	return &BucketNode{
		BaseNode: NewBaseNode(node),
		Bucket:   node.GetBucket(),
	}
}

func NewInputNode(node *gen.Node) *InputNode {
	return &InputNode{
		BaseNode: NewBaseNode(node),
		Input:    node.GetInput(),
	}
}

func NewFunctionNode(node *gen.Node) *FunctionNode {
	return &FunctionNode{
		BaseNode: NewBaseNode(node),
		Function: node.GetFunction(),
	}
}

func NewQueryNode(resources ResourceMap, node *gen.Node) *QueryNode {
	// TODO breadchris a node should already have this resource configured
	for id, r := range resources {
		if r.Name() == DocstoreResourceType {
			node.ResourceIds = append(node.ResourceIds, id)
		}
	}
	return &QueryNode{
		BaseNode: NewBaseNode(node),
		Query:    node.GetQuery(),
	}
}
