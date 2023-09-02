package storage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen/storage"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"gocloud.dev/docstore"
	"gocloud.dev/docstore/memdocstore"
	"io"
	"os"
	"path"
	"strings"
)

type Store struct {
	*base.Node
	*storage.Store
}

var _ graph.Node = &Store{}

func NewStore(b *base.Node, n *storage.Store) *Store {
	return &Store{
		Node:  b,
		Store: n,
	}
}

func (r *Store) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Store) WithCollection(name string) (*docstore.Collection, func(), error) {
	var (
		coll     *docstore.Collection
		err      error
		protoDir string
	)
	if strings.HasPrefix(r.Url, "mem://") {
		// TODO breadchris replace this with bucket.Cache.GetFolder
		protoDir, err = util.ProtoflowHomeDir()
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not get protoflow home dir")
		}

		filename := path.Join(protoDir, name+".json")

		// TODO breadchris "id" is
		coll, err = memdocstore.OpenCollection("id", &memdocstore.Options{
			Filename: filename,
		})
		if err != nil {
			// remove file if it exists
			if os.IsNotExist(err) {
				return nil, nil, errors.Wrapf(err, "could not open memory docstore collection: %s", name)
			}
			err = os.Remove(filename)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "could not remove memory docstore collection: %s", name)
			}
		}
	} else {
		coll, err = docstore.OpenCollection(context.Background(), r.Url+"/"+name)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "could not open docstore collection: %s", name)
		}
	}

	return coll, func() {
		if coll == nil {
			log.Debug().Msg("docstore collection is nil")
			return
		}
		err = coll.Close()
		if err != nil {
			log.Error().Msgf("error closing docstore collection: %+v", err)
		}
	}, nil
}

type Collection struct {
	*base.Node
	c          *storage.Collection
	collection *docstore.Collection
}

func NewCollection(b *base.Node, n *storage.Collection) *Collection {
	return &Collection{
		Node: b,
		c:    n,
	}
}

func (n *Collection) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	p, err := n.Provider()
	if err != nil {
		return graph.IO{}, errors.Wrapf(err, "error getting provider")
	}
	d, ok := p.(*Store)
	if !ok {
		return graph.IO{}, errors.New("error provider is not a docstore")
	}

	collection, cleanup, err := d.WithCollection(n.c.Name)
	if err != nil {
		return graph.IO{}, errors.Wrapf(err, "error connecting to collection")
	}
	n.collection = collection

	insertWithID := func(record map[string]any) (string, error) {
		if record["id"] == nil {
			record["id"] = uuid.NewString()
		}
		err := n.collection.Create(context.Background(), record)
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
	}, func() {
		close(output)
		cleanup()
	})

	return graph.IO{
		Observable: rxgo.FromChannel(output, rxgo.WithPublishStrategy()),
	}, nil
}

type QueryNode struct {
	*base.Node
	*storage.Query
}

var _ graph.Node = &QueryNode{}

func NewQuery(b *base.Node, n *storage.Query) *QueryNode {
	return &QueryNode{
		Node:  b,
		Query: n,
	}
}

func (n *QueryNode) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	p, err := n.Provider()
	if err != nil {
		return graph.IO{}, errors.Wrapf(err, "error getting provider")
	}
	d, ok := p.(*Collection)
	if !ok {
		return graph.IO{}, errors.New("error provider is not a docstore")
	}

	// TODO breadchris if input is specified, use it as the query

	output := make(chan rxgo.Item)
	go func() {
		iter := d.collection.Query().Get(ctx)
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

	return graph.IO{
		Observable: rxgo.FromChannel(output, rxgo.WithPublishStrategy()),
	}, nil
}
