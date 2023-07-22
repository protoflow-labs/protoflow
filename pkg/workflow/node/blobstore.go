package node

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/reactivex/rxgo/v2"
	"net/url"
	"path"
)

type BucketNode struct {
	BaseNode
	*gen.Bucket
}

var _ graph.Node = &BucketNode{}

func NewBucketNode(node *gen.Node) *BucketNode {
	return &BucketNode{
		BaseNode: NewBaseNode(node),
		Bucket:   node.GetBucket(),
	}
}

func (n *BucketNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	bucket, ok := input.Resource.(*resource.FileStoreResource)
	if !ok {
		return graph.Output{}, fmt.Errorf("error getting blobstore resource: %s", n.Bucket.Path)
	}

	item, err := input.Observable.First().Get()
	if err != nil {
		return graph.Output{}, errors.Wrapf(err, "error getting first item from observable")
	}

	var (
		bucketData []byte
	)
	switch t := item.V.(type) {
	case []byte:
		bucketData = t
	case string:
		bucketData = []byte(t)
	default:
		bucketData, err = json.Marshal(t)
		if err != nil {
			return graph.Output{}, errors.Wrapf(err, "error marshaling input params")
		}
	}

	b, cleanup, err := bucket.WithPath(n.Path)
	if err != nil {
		return graph.Output{}, errors.Wrapf(err, fmt.Sprintf("error connecting to bucket: %s", n.Path))
	}
	defer cleanup()

	err = b.WriteAll(context.Background(), n.Path, bucketData, nil)
	return graph.Output{
		Observable: rxgo.Just(map[string]string{
			"bucket": n.Path,
		})(),
	}, nil
}

type FileNode struct {
	BaseNode
	File *gen.File
}

var _ graph.Node = &FileNode{}

func NewFileNode(node *gen.Node) *FileNode {
	return &FileNode{
		BaseNode: NewBaseNode(node),
		File:     node.GetFile(),
	}
}

func (n *FileNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	fs, ok := input.Resource.(*resource.FileStoreResource)
	if !ok {
		return graph.Output{}, fmt.Errorf("error getting filestore resource: %s", n.File.Path)
	}
	u, err := url.Parse(fs.Url)
	if err != nil {
		return graph.Output{}, errors.Wrapf(err, "error parsing filestore url")
	}
	p := path.Join(u.Path, n.File.Path)

	// TODO breadchris verify file exists?
	obs := rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		// TODO breadchris this should be a static type. This is a brittle type that maps to workflow.go:133
		next <- rx.NewItem(map[string]any{
			"path": p,
		})
	}})

	return graph.Output{
		Observable: obs,
	}, nil
}
