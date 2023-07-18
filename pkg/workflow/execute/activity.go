package execute

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	grpcanal "github.com/protoflow-labs/protoflow/pkg/grpc/manager"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/url"
	"path"
)

type Activity struct{}

type ActivityFunc func(ctx context.Context, n node.Node, input Input) (Output, error)

func formatHost(host string) (string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return "", errors.Wrapf(err, "error parsing url: %s", host)
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return u.String(), nil
}

func (a *Activity) ExecutePromptNode(ctx context.Context, n node.Node, input Input) (Output, error) {
	pn, ok := n.(*node.PromptNode)
	if !ok {
		return Output{}, fmt.Errorf("error getting prompt node: %s", pn.Name)
	}

	r, ok := input.Resource.(*resource.ReasoningEngineResource)
	if !ok {
		return Output{}, fmt.Errorf("error getting reasoning engine resource: %s", pn.Name)
	}

	log.Info().
		Str("name", pn.Name).
		Msg("setting up prompt node")

	outputStream := make(chan rxgo.Item)

	// TODO breadchris how should context be handled here?
	c := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: pn.Prompt.Prompt,
		},
	}

	input.Observable.ForEach(func(item any) {
		log.Debug().
			Str("name", pn.NormalizedName()).
			Interface("item", item).
			Msg("executing prompt node")

		var normalizedItem string
		switch t := item.(type) {
		case string:
			normalizedItem = t
		default:
			c, err := json.Marshal(item)
			if err != nil {
				outputStream <- rx.NewError(errors.Wrapf(err, "error marshalling input: %s", pn.NormalizedName()))
				return
			}
			normalizedItem = string(c)
		}

		c = append(c, openai.ChatCompletionMessage{
			Role:    "user",
			Content: normalizedItem,
		})

		s, err := r.QAClient.Ask(c)
		if err != nil {
			outputStream <- rx.NewError(errors.Wrapf(err, "error executing prompt: %s", pn.NormalizedName()))
		}

		// TODO breadchris react to a function call on s.FunctionCall

		// TODO breadchris this should be a static type. This is a brittle type that maps to workflow.go:133
		outputStream <- rx.NewItem(map[string]any{
			"result": s,
		})
	}, func(err error) {
		outputStream <- rx.NewError(err)
	}, func() {
		close(outputStream)
	})
	res := Output{
		Observable: rxgo.FromChannel(outputStream, rxgo.WithPublishStrategy()),
	}
	return res, nil
}

// TODO breadchris this should be workflow.Context, but for the memory executor it needs context.Context
func (a *Activity) ExecuteGRPCNode(ctx context.Context, n node.Node, input Input) (Output, error) {
	gn, ok := n.(*node.GRPCNode)
	if !ok {
		return Output{}, fmt.Errorf("error getting GRPC resource: %s.%s", gn.Service, gn.Method)
	}

	log.Info().
		Str("service", gn.Service).
		Str("method", gn.Method).
		Msg("setting up grpc node")

	g, ok := input.Resource.(*resource.GRPCResource)
	if !ok {
		return Output{}, fmt.Errorf("error getting GRPC resource: %s.%s", gn.Service, gn.Method)
	}

	serviceName := gn.Service
	if gn.Package != "" {
		serviceName = gn.Package + "." + serviceName
	}

	host, err := formatHost(g.Host)
	if err != nil {
		return Output{}, errors.Wrapf(err, "error formatting host: %s", g.Host)
	}

	manager := grpcanal.NewReflectionManager(host)

	cleanup, err := manager.Init()
	if err != nil {
		return Output{}, errors.Wrapf(err, "error initializing reflection manager")
	}
	defer cleanup()

	method, err := manager.ResolveMethod(serviceName, gn.Method)
	if err != nil {
		return Output{}, errors.Wrapf(err, "error resolving method: %s.%s", serviceName, gn.Method)
	}

	outputStream := make(chan rxgo.Item)
	// TODO breadchris we are relying on this grpc call to close the output stream. How can the stream be closed by the caller?
	if !method.IsStreamingClient() {
		// if the method is not a client stream, we need to send each input observable as a single request
		// TODO breadchris type of this method should be inferred when the workflow is parsed
		input.Observable.ForEach(func(item any) {
			log.Debug().
				Str("name", gn.NormalizedName()).
				Interface("item", item).
				Msg("executing single grpc method")

			err = manager.ExecuteMethod(ctx, method, rx.FromValues(item), outputStream)
			if err != nil {
				outputStream <- rx.NewError(errors.Wrapf(err, "error calling grpc method: %s", host))
			}
			log.Debug().
				Str("name", gn.NormalizedName()).
				Interface("item", item).
				Msg("done executing single grpc method")
		}, func(err error) {
			outputStream <- rx.NewError(err)
		}, func() {
			close(outputStream)
		})
	} else {
		go func() {
			log.Debug().
				Str("name", n.NormalizedName()).
				Msg("executing streaming grpc method")
			defer close(outputStream)
			err = manager.ExecuteMethod(ctx, method, input.Observable, outputStream)
			if err != nil {
				outputStream <- rx.NewError(errors.Wrapf(err, "error calling grpc method: %s", host))
			}
		}()
	}
	res := Output{
		Observable: rxgo.FromChannel(outputStream, rxgo.WithPublishStrategy()),
	}
	return res, nil
}

func (a *Activity) ExecuteRestNode(ctx context.Context, n node.Node, input Input) (Output, error) {
	gn, ok := n.(*node.RESTNode)
	if !ok {
		return Output{}, fmt.Errorf("error getting REST resource: %s", gn.Name)
	}
	log.Debug().
		Interface("headers", gn.Headers).
		Str("method", gn.Method).
		Str("path", gn.Path).
		Msgf("executing rest")
	// TODO breadchris turn this into streamable because why not
	item, err := input.Observable.First().Get()
	if err != nil {
		return Output{}, errors.Wrapf(err, "error getting first item from observable")
	}
	res, err := util.InvokeMethodOnUrl(gn.Method, gn.Path, gn.Headers, item.V)
	if err != nil {
		return Output{Observable: rxgo.Empty()}, nil
	}
	return Output{
		Observable: rxgo.Just(res)(),
	}, nil
}

func (a *Activity) ExecuteFunctionNode(ctx context.Context, n node.Node, input Input) (Output, error) {
	gn, ok := n.(*node.FunctionNode)
	if !ok {
		return Output{}, fmt.Errorf("error getting Function resource: %s", gn.Name)
	}
	log.Debug().Str("name", gn.Name).Msg("setting up function")
	g, ok := input.Resource.(*resource.LanguageServiceResource)
	if !ok {
		return Output{}, fmt.Errorf("error getting language service resource: %s", gn.Name)
	}

	// provide the grpc resource to the grpc gn call. Is this the best place for this? Should this be provided on injection? Probably.
	input.Resource = g.GRPCResource

	grpcNode := g.ToGRPC(gn)
	return a.ExecuteGRPCNode(ctx, grpcNode, input)
}

func (a *Activity) ExecuteCollectionNode(ctx context.Context, n node.Node, input Input) (Output, error) {
	gn, ok := n.(*node.CollectionNode)
	if !ok {
		return Output{}, fmt.Errorf("error getting Collection resource: %s", gn.Name)
	}
	docs, ok := input.Resource.(*resource.DocstoreResource)
	if !ok {
		return Output{}, fmt.Errorf("error getting docstore resource: %s", gn.Collection.Name)
	}

	collection, cleanup, err := docs.WithCollection(gn.Collection.Name)
	if err != nil {
		return Output{}, errors.Wrapf(err, "error connecting to collection")
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

	return Output{
		Observable: rxgo.FromChannel(output, rxgo.WithPublishStrategy()),
	}, nil
}

func (a *Activity) ExecuteInputNode(ctx context.Context, n node.Node, input Input) (Output, error) {
	return Output{
		Observable: input.Observable,
	}, nil
}

func (a *Activity) ExecuteFileNode(ctx context.Context, n node.Node, input Input) (Output, error) {
	f, ok := n.(*node.FileNode)
	if !ok {
		return Output{}, fmt.Errorf("error getting File node: %s", n.NormalizedName())
	}
	fs, ok := input.Resource.(*resource.FileStoreResource)
	if !ok {
		return Output{}, fmt.Errorf("error getting filestore resource: %s", f.File.Path)
	}
	u, err := url.Parse(fs.Url)
	if err != nil {
		return Output{}, errors.Wrapf(err, "error parsing filestore url")
	}
	p := path.Join(u.Path, f.File.Path)

	// TODO breadchris verify file exists?
	obs := rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		// TODO breadchris this should be a static type. This is a brittle type that maps to workflow.go:133
		next <- rx.NewItem(map[string]any{
			"path": p,
		})
	}})

	return Output{
		Observable: obs,
	}, nil
}

func (a *Activity) ExecuteBucketNode(ctx context.Context, n node.Node, input Input) (Output, error) {
	gn, ok := n.(*node.BucketNode)
	if !ok {
		return Output{}, fmt.Errorf("error getting Collection resource: %s", gn.Name)
	}
	bucket, ok := input.Resource.(*resource.FileStoreResource)
	if !ok {
		return Output{}, fmt.Errorf("error getting blobstore resource: %s", gn.Bucket.Path)
	}

	item, err := input.Observable.First().Get()
	if err != nil {
		return Output{}, errors.Wrapf(err, "error getting first item from observable")
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
			return Output{}, errors.Wrapf(err, "error marshaling input params")
		}
	}

	b, cleanup, err := bucket.WithPath(gn.Path)
	if err != nil {
		return Output{}, errors.Wrapf(err, fmt.Sprintf("error connecting to bucket: %s", gn.Path))
	}
	defer cleanup()

	err = b.WriteAll(context.Background(), gn.Path, bucketData, nil)
	return Output{
		Observable: rxgo.Just(map[string]string{
			"bucket": gn.Path,
		})(),
	}, nil
}

func (a *Activity) ExecuteQueryNode(ctx context.Context, n node.Node, input Input) (Output, error) {
	s, ok := n.(*node.QueryNode)
	if !ok {
		return Output{}, fmt.Errorf("error getting query resource: %s", s.Query.Collection)
	}
	docResource, ok := input.Resource.(*resource.DocstoreResource)
	if !ok {
		return Output{}, fmt.Errorf("error getting docstore resource: %s", s.Query.Collection)
	}

	d, cleanup, err := docResource.WithCollection(s.Query.Collection)
	if err != nil {
		return Output{}, errors.Wrapf(err, "error connecting to collection")
	}

	output := make(chan rxgo.Item)
	go func() {
		defer cleanup()
		iter := d.Query().Get(ctx)
		for {
			doc := map[string]interface{}{}
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

	return Output{
		Observable: rxgo.FromChannel(output, rxgo.WithPublishStrategy()),
	}, nil
}
