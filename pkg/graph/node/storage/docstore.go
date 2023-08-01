package storage

//type DocstoreResource struct {
//	*BaseResource
//	*gen.DocStore
//}
//
//var _ graph.Resource = &DocstoreResource{}
//
//func (r *DocstoreResource) Init() (func(), error) {
//	return nil, nil
//}
//
//func (r *DocstoreResource) WithCollection(name string) (*docstore.Collection, func(), error) {
//	var (
//		coll     *docstore.Collection
//		err      error
//		protoDir string
//	)
//	if strings.HasPrefix(r.Url, "mem://") {
//		// TODO breadchris replace this with bucket.Cache.GetFolder
//		protoDir, err = util.ProtoflowHomeDir()
//		if err != nil {
//			return nil, nil, errors.Wrap(err, "could not get protoflow home dir")
//		}
//
//		filename := path.Join(protoDir, name+".json")
//
//		// TODO breadchris "id" is
//		coll, err = memdocstore.OpenCollection("id", &memdocstore.Options{
//			Filename: filename,
//		})
//		if err != nil {
//			// remove file if it exists
//			if os.IsNotExist(err) {
//				return nil, nil, errors.Wrapf(err, "could not open memory docstore collection: %s", name)
//			}
//			err = os.Remove(filename)
//			if err != nil {
//				return nil, nil, errors.Wrapf(err, "could not remove memory docstore collection: %s", name)
//			}
//		}
//	} else {
//		coll, err = docstore.OpenCollection(context.Background(), r.Url+"/"+name)
//		if err != nil {
//			return nil, nil, errors.Wrapf(err, "could not open docstore collection: %s", name)
//		}
//	}
//
//	return coll, func() {
//		if coll == nil {
//			log.Debug().Msg("docstore collection is nil")
//			return
//		}
//		err = coll.Close()
//		if err != nil {
//			log.Error().Msgf("error closing docstore collection: %+v", err)
//		}
//	}, nil
//}

//type CollectionNode struct {
//	Node
//	Collection *gen.Collection
//}
//
//var _ graph.Node = &CollectionNode{}
//
//func NewCollectionNode(node *gen.Node) *CollectionNode {
//	return &CollectionNode{
//		Node:   NewBaseNode(node),
//		Collection: node.GetCollection(),
//	}
//}
//
//func (n *CollectionNode) Wire(ctx context.Context, input graph.Input) (graph.Input, error) {
//	docs, ok := input.Resource.(*resource.DocstoreResource)
//	if !ok {
//		return graph.Input{}, fmt.Errorf("error getting docstore resource: %s", n.Collection.Name)
//	}
//
//	collection, cleanup, err := docs.WithCollection(n.Collection.Name)
//	if err != nil {
//		return graph.Input{}, errors.Wrapf(err, "error connecting to collection")
//	}
//
//	insertWithID := func(record map[string]any) (string, error) {
//		if record["id"] == nil {
//			record["id"] = uuid.NewString()
//		}
//		err = collection.Create(context.Background(), record)
//		if err != nil {
//			return "", errors.Wrapf(err, "error creating doc")
//		}
//		return record["id"].(string), nil
//	}
//
//	output := make(chan rxgo.Item)
//	input.Observable.ForEach(func(item any) {
//		var (
//			id  string
//			err error
//		)
//		switch i := item.(type) {
//		case map[string]interface{}:
//			id, err = insertWithID(i)
//			output <- rx.NewItem(id)
//		case []*map[string]interface{}:
//			for _, record := range i {
//				id, err = insertWithID(*record)
//				if err != nil {
//					break
//				}
//				output <- rx.NewItem(id)
//			}
//		case string:
//			id, err = insertWithID(map[string]interface{}{
//				"input": i,
//			})
//			output <- rx.NewItem(id)
//		default:
//			err = fmt.Errorf("error unsupported input type: %T", input)
//		}
//		if err != nil {
//			output <- rx.NewError(errors.Wrapf(err, "error inserting record"))
//		}
//	}, func(err error) {
//		output <- rx.NewError(err)
//		// TODO breadchris cleanup and close here too?
//	}, func() {
//		cleanup()
//		close(output)
//	})
//
//	return graph.Input{
//		Observable: rxgo.FromChannel(output, rxgo.WithPublishStrategy()),
//	}, nil
//}
//
//type QueryNode struct {
//	Node
//	Query *gen.Query
//}
//
//var _ graph.Node = &QueryNode{}
//
//func NewQueryNode(node *gen.Node) *QueryNode {
//	return &QueryNode{
//		Node: NewBaseNode(node),
//		Query:    node.GetQuery(),
//	}
//}
//
//func (n *QueryNode) Wire(ctx context.Context, input graph.Input) (graph.Input, error) {
//	docResource, ok := input.Resource.(*resource.DocstoreResource)
//	if !ok {
//		return graph.Input{}, fmt.Errorf("error getting docstore resource: %s", n.Query.Collection)
//	}
//
//	d, cleanup, err := docResource.WithCollection(n.Query.Collection)
//	if err != nil {
//		return graph.Input{}, errors.Wrapf(err, "error connecting to collection")
//	}
//
//	output := make(chan rxgo.Item)
//	go func() {
//		defer cleanup()
//		iter := d.Query().Get(ctx)
//		for {
//			doc := map[string]any{}
//			err = iter.Next(ctx, doc)
//			if err != nil {
//				if err != io.EOF {
//					output <- rx.NewError(errors.Wrapf(err, "error iterating over query results"))
//				}
//				close(output)
//				break
//			}
//			output <- rx.NewItem(doc)
//		}
//	}()
//
//	return graph.Input{
//		Observable: rxgo.FromChannel(output, rxgo.WithPublishStrategy()),
//	}, nil
//}
