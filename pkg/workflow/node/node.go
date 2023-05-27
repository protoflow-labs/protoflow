package node

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"github.com/protoflow-labs/protoflow/pkg/workflow/execute"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"io"
	"strings"
)

type Node interface {
	Execute(executor execute.Executor, input execute.Input) (*execute.Result, error)
	NormalizedName() string
	ID() string
	ResourceID() string
	Info(r resource.Resource) (*Info, error)
}

type Info struct {
	MethodProto string
	TypeInfo    *gen.GRPCTypeInfo
}

type BaseNode struct {
	Name       string
	id         string
	resourceID string
}

func (n *BaseNode) NormalizedName() string {
	name := util.ToTitleCase(n.Name)
	if strings.Contains(name, ".") {
		name = strings.Split(name, ".")[1]
	}
	return name
}

func (n *BaseNode) ID() string {
	return n.id
}

func (n *BaseNode) ResourceID() string {
	return n.resourceID
}

func (n *BaseNode) Info(r resource.Resource) (*Info, error) {
	return &Info{}, nil
}

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

type QueryNode struct {
	BaseNode
	Query *gen.Query
}

var activity = &Activity{}

func (s *RESTNode) Execute(executor execute.Executor, input execute.Input) (*execute.Result, error) {
	return executor.Execute(activity.ExecuteRestNode, s, input)
}

func (s *CollectionNode) Execute(executor execute.Executor, input execute.Input) (*execute.Result, error) {
	docs, ok := input.Resource.(*resource.DocstoreResource)
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

	return &execute.Result{
		Data: input.Params,
	}, nil
}

func (s *BucketNode) Execute(executor execute.Executor, input execute.Input) (*execute.Result, error) {
	bucket, ok := input.Resource.(*resource.BlobstoreResource)
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
	return &execute.Result{
		Data: input.Params,
	}, nil
}

func (s *QueryNode) Execute(executor execute.Executor, input execute.Input) (*execute.Result, error) {
	docResource, ok := input.Resource.(*resource.DocstoreResource)
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
	return &execute.Result{
		Data: docs,
	}, nil
}

type ResourceMap map[string]resource.Resource

func NewNode(node *gen.Node) (Node, error) {
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
