package node

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/reactivex/rxgo/v2"
	"io"
)

type CollectionNode struct {
	BaseNode
	Collection *gen.Collection
}

var _ graph.Node = &CollectionNode{}

func NewCollectionNode(node *gen.Node) *CollectionNode {
	return &CollectionNode{
		BaseNode:   NewBaseNode(node),
		Collection: node.GetCollection(),
	}
}

func (n *CollectionNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	docs, ok := input.Resource.(*resource.DocstoreResource)
	if !ok {
		return graph.Output{}, fmt.Errorf("error getting docstore resource: %s", n.Collection.Name)
	}

	collection, cleanup, err := docs.WithCollection(n.Collection.Name)
	if err != nil {
		return graph.Output{}, errors.Wrapf(err, "error connecting to collection")
	}

	insertWithID := func(record map[string]any) (string, error) {
		if record["id"] == nil {
			record["id"] = uuid.NewString()
		}
		err = collection.Create(context.Background(), record)
		if err != nil {
			return "", errors.Wrapf(err, "error creating doc")
		}
		return record["id"].(string), nil
	}

	output := make(chan rxgo.Item)
	input.Observable.ForEach(func(item any) {
		var (
			id  string
			err error
		)
		switch i := item.(type) {
		case map[string]interface{}:
			id, err = insertWithID(i)
			output <- rx.NewItem(id)
		case []*map[string]interface{}:
			for _, record := range i {
				id, err = insertWithID(*record)
				if err != nil {
					break
				}
				output <- rx.NewItem(id)
			}
		case string:
			id, err = insertWithID(map[string]interface{}{
				"input": i,
			})
			output <- rx.NewItem(id)
		default:
			err = fmt.Errorf("error unsupported input type: %T", input)
		}
		if err != nil {
			output <- rx.NewError(errors.Wrapf(err, "error inserting record"))
		}
	}, func(err error) {
		output <- rx.NewError(err)
		// TODO breadchris cleanup and close here too?
	}, func() {
		cleanup()
		close(output)
	})

	return graph.Output{
		Observable: rxgo.FromChannel(output, rxgo.WithPublishStrategy()),
	}, nil
}

type QueryNode struct {
	BaseNode
	Query *gen.Query
}

var _ graph.Node = &QueryNode{}

func NewQueryNode(node *gen.Node) *QueryNode {
	return &QueryNode{
		BaseNode: NewBaseNode(node),
		Query:    node.GetQuery(),
	}
}

func (n *QueryNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	docResource, ok := input.Resource.(*resource.DocstoreResource)
	if !ok {
		return graph.Output{}, fmt.Errorf("error getting docstore resource: %s", n.Query.Collection)
	}

	d, cleanup, err := docResource.WithCollection(n.Query.Collection)
	if err != nil {
		return graph.Output{}, errors.Wrapf(err, "error connecting to collection")
	}

	output := make(chan rxgo.Item)
	go func() {
		defer cleanup()
		iter := d.Query().Get(ctx)
		for {
			doc := map[string]any{}
			err = iter.Next(ctx, doc)
			if err != nil {
				if err != io.EOF {
					output <- rx.NewError(errors.Wrapf(err, "error iterating over query results"))
				}
				close(output)
				break
			}
			output <- rx.NewItem(doc)
		}
	}()

	return graph.Output{
		Observable: rxgo.FromChannel(output, rxgo.WithPublishStrategy()),
	}, nil
}
