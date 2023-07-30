package storage

//type FileStoreResource struct {
//	*BaseResource
//	*gen.FileStore
//}
//
//var _ graph.Resource = &FileStoreResource{}
//
//func (r *FileStoreResource) Init() (func(), error) {
//	return nil, nil
//}
//
//func (r *FileStoreResource) WithPath(path string) (*blob.Bucket, func(), error) {
//	// remove leading slash
//	if path[0] == '/' {
//		path = path[1:]
//	}
//	// TODO breadchris validation of this url working should be done on init
//	bucket, err := blob.OpenBucket(context.Background(), r.Url+"?prefix="+path)
//	if err != nil {
//		return nil, nil, fmt.Errorf("could not open bucket: %v", err)
//	}
//	return bucket, func() {
//		err = bucket.Close()
//		if err != nil {
//			log.Error().Err(err).Msg("error closing blobstore bucket")
//		}
//	}, nil
//}

//type BucketNode struct {
//	Node
//	*gen.Bucket
//}
//
//var _ graph.Node = &BucketNode{}
//
//func NewBucketNode(node *gen.Node) *BucketNode {
//	return &BucketNode{
//		Node: NewBaseNode(node),
//		Bucket:   node.GetBucket(),
//	}
//}
//
//func (n *BucketNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
//	bucket, ok := input.Resource.(*resource.FileStoreResource)
//	if !ok {
//		return graph.Output{}, fmt.Errorf("error getting blobstore resource: %s", n.Bucket.Path)
//	}
//
//	item, err := input.Observable.First().Get()
//	if err != nil {
//		return graph.Output{}, errors.Wrapf(err, "error getting first item from observable")
//	}
//
//	var (
//		bucketData []byte
//	)
//	switch t := item.V.(type) {
//	case []byte:
//		bucketData = t
//	case string:
//		bucketData = []byte(t)
//	default:
//		bucketData, err = json.Marshal(t)
//		if err != nil {
//			return graph.Output{}, errors.Wrapf(err, "error marshaling input params")
//		}
//	}
//
//	b, cleanup, err := bucket.WithPath(n.Path)
//	if err != nil {
//		return graph.Output{}, errors.Wrapf(err, fmt.Sprintf("error connecting to bucket: %s", n.Path))
//	}
//	defer cleanup()
//
//	err = b.WriteAll(context.Background(), n.Path, bucketData, nil)
//	return graph.Output{
//		Observable: rxgo.Just(map[string]string{
//			"bucket": n.Path,
//		})(),
//	}, nil
//}
//
//type FileNode struct {
//	Node
//	File *gen.File
//}
//
//var _ graph.Node = &FileNode{}
//
//func NewFileNode(node *gen.Node) *FileNode {
//	return &FileNode{
//		Node: NewBaseNode(node),
//		File:     node.GetFile(),
//	}
//}
//
//func (n *FileNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
//	fs, ok := input.Resource.(*resource.FileStoreResource)
//	if !ok {
//		return graph.Output{}, fmt.Errorf("error getting filestore resource: %s", n.File.Path)
//	}
//	u, err := url.Parse(fs.Url)
//	if err != nil {
//		return graph.Output{}, errors.Wrapf(err, "error parsing filestore url")
//	}
//	p := path.Join(u.Path, n.File.Path)
//
//	// TODO breadchris verify file exists?
//	obs := rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
//		// TODO breadchris this should be a static type. This is a brittle type that maps to workflow.go:133
//		next <- rx.NewItem(map[string]any{
//			"path": p,
//		})
//	}})
//
//	return graph.Output{
//		Observable: obs,
//	}, nil
//}
