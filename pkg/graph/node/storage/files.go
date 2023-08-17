package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen/storage"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"gocloud.dev/blob"
	"net/url"
	"path"
)

type Folder struct {
	*base.Node
	*storage.Folder
}

var _ graph.Node = &Folder{}

func NewFolder(b *base.Node, node *storage.Folder) *Folder {
	return &Folder{
		Node:   b,
		Folder: node,
	}
}

func (r *Folder) Init() (func(), error) {
	return nil, nil
}

func (r *Folder) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	u, err := url.Parse(r.Url)
	if err != nil {
		return graph.IO{}, errors.Wrapf(err, "error parsing url: %s", r.Url)
	}
	output := input.Observable.Map(func(ctx context.Context, i any) (any, error) {
		var (
			bucketData []byte
			err        error
		)
		switch t := i.(type) {
		case []byte:
			bucketData = t
		case string:
			bucketData = []byte(t)
		default:
			bucketData, err = json.Marshal(t)
			if err != nil {
				return graph.IO{}, errors.Wrapf(err, "error marshaling input params")
			}
		}

		// TODO breadchris this path should have a filename, do we generate a random one? accept metadata on block input?
		b, cleanup, err := r.WithPath(u.Path)
		if err != nil {
			return graph.IO{}, errors.Wrapf(err, fmt.Sprintf("error connecting to bucket: %s", u.Path))
		}
		defer cleanup()

		err = b.WriteAll(context.Background(), u.Path, bucketData, nil)
		return nil, err
	})
	return graph.IO{
		Observable: output,
	}, nil
}

func (r *Folder) WithPath(path string) (*blob.Bucket, func(), error) {
	// remove leading slash
	if path[0] == '/' {
		path = path[1:]
	}
	// TODO breadchris validation of this url working should be done on init
	bucket, err := blob.OpenBucket(context.Background(), r.Url+"?prefix="+path)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open bucket: %v", err)
	}
	return bucket, func() {
		err = bucket.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing blobstore bucket")
		}
	}, nil
}

type File struct {
	*base.Node
	*storage.File
}

var _ graph.Node = &File{}

func NewFile(b *base.Node, node *storage.File) *File {
	return &File{
		Node: b,
		File: node,
	}
}

func (n *File) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	p, err := n.Provider()
	if err != nil {
		return graph.IO{}, errors.Wrapf(err, "error getting provider")
	}

	f, ok := p.(*Folder)
	if !ok {
		return graph.IO{}, errors.Wrapf(err, "error getting folder")
	}
	u, err := url.Parse(f.Url)
	if err != nil {
		return graph.IO{}, errors.Wrapf(err, "error parsing filestore url")
	}
	filepath := path.Join(u.Path, n.File.Path)

	// TODO breadchris verify file exists?
	obs := rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		// TODO breadchris this should be a static type. This is a brittle type that maps to workflow.go:133
		next <- rx.NewItem(map[string]any{
			"path": filepath,
		})
	}})
	return graph.IO{
		Observable: obs,
	}, nil
}
