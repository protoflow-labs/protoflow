package execute

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl"
	grpcanal "github.com/protoflow-labs/protoflow/pkg/grpc/manager"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/rs/zerolog/log"
	"io"
	"net/url"
)

type Activity struct{}

type ActivityFunc func(ctx context.Context, n node.Node, input Input) (Result, error)

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

// TODO breadchris this should be workflow.Context, but for the memory executor it needs context.Context
func (a *Activity) ExecuteGRPCNode(ctx context.Context, n node.Node, input Input) (Result, error) {
	gn, ok := n.(*node.GRPCNode)
	if !ok {
		return Result{}, fmt.Errorf("error getting GRPC resource: %s.%s", gn.Service, gn.Method)
	}

	log.Info().Msgf("executing node: %s", gn.Service)

	g, ok := input.Resource.(*resource.GRPCResource)
	if !ok {
		return Result{}, fmt.Errorf("error getting GRPC resource: %s.%s", gn.Service, gn.Method)
	}

	serviceName := gn.Service
	if gn.Package != "" {
		serviceName = gn.Package + "." + serviceName
	}

	host, err := formatHost(g.Host)
	if err != nil {
		return Result{}, errors.Wrapf(err, "error formatting host: %s", g.Host)
	}

	inputStream := bufcurl.NewMemoryInputStream()
	outputStream := bufcurl.NewMemoryOutputStream()

	manager := grpcanal.NewReflectionManager(host)
	cleanup, err := manager.Init()
	if err != nil {
		return Result{}, errors.Wrapf(err, "error initializing reflection manager")
	}
	defer cleanup()

	method, err := manager.ResolveMethod(serviceName, gn.Method)
	if err != nil {
		return Result{}, errors.Wrapf(err, "error resolving method: %s.%s", serviceName, gn.Method)
	}

	go func() {
		// TODO breadchris we are relying on this grpc call to close the output stream. How can the stream be closed by the caller?
		defer outputStream.Close()
		err = manager.ExecuteMethod(ctx, method, inputStream, outputStream)
		if err != nil {
			outputStream.Error(errors.Wrapf(err, "error calling grpc method: %s", host))
		}
	}()
	go func() {
		defer inputStream.Close()
		// TODO breadchris such indentation
		if input.Stream != nil {
			for {
				d, err := input.Stream.Next()
				if err != nil {
					if err != io.EOF {
						outputStream.Error(err)
					}
					break
				}

				log.Debug().Msgf("pushing data to grpc input stream: %s", d)
				// TODO breadchris need to transform the data to the correct type
				// we need to plug in the correct type for the method's input

				inputStream.Push(d)
			}
		} else {
			inputStream.Push(input.Params)
		}
	}()

	var res Result
	switch {
	case method.IsStreamingServer() && method.IsStreamingClient():
		// bidirectional stream: multiple requests, multiple responses
		fallthrough
	case method.IsStreamingServer():
		// server stream: one request, multiple responses
		res.Stream = outputStream
	case method.IsStreamingClient():
		// client stream: multiple requests, one response
		go func() {
			for {
				output, err := outputStream.Next()
				if err != nil {
					if err != io.EOF {
						outputStream.Error(errors.Wrapf(err, "error reading output stream"))
					}
					break
				}
				res.Stream.Push(output)
			}
		}()
	default:
		// unary
		output, err := outputStream.Next()
		if err != nil {
			if err != io.EOF {
				return Result{}, errors.Wrapf(err, "error reading output stream")
			}
			return Result{}, errors.Wrapf(err, "did not receive any output for unary call")
		}
		res.Data = output
	}
	return res, nil
}

func (a *Activity) ExecuteRestNode(ctx context.Context, n node.Node, input Input) (Result, error) {
	gn, ok := n.(*node.RESTNode)
	if !ok {
		return Result{}, fmt.Errorf("error getting REST resource: %s", gn.Name)
	}
	log.Debug().
		Interface("headers", gn.Headers).
		Str("method", gn.Method).
		Str("path", gn.Path).
		Msgf("executing rest")
	res, err := util.InvokeMethodOnUrl(gn.Method, gn.Path, gn.Headers, input.Params)
	if err != nil {
		return Result{}, errors.Wrapf(err, "error invoking method: %s", gn.Method)
	}
	return Result{
		Data: res,
	}, nil
}

func (a *Activity) ExecuteFunctionNode(ctx context.Context, n node.Node, input Input) (Result, error) {
	gn, ok := n.(*node.FunctionNode)
	if !ok {
		return Result{}, fmt.Errorf("error getting Function resource: %s", gn.Name)
	}
	log.Debug().Msgf("executing function gn: %s", gn.Name)
	g, ok := input.Resource.(*resource.LanguageServiceResource)
	if !ok {
		return Result{}, fmt.Errorf("error getting language service resource: %s", gn.Name)
	}

	// provide the grpc resource to the grpc gn call. Is this the best place for this? Should this be provided on injection? Probably.
	input.Resource = g.GRPCResource

	grpcNode := g.ToGRPC(gn)
	return a.ExecuteGRPCNode(ctx, grpcNode, input)
}

func (a *Activity) ExecuteCollectionNode(ctx context.Context, n node.Node, input Input) (Result, error) {
	gn, ok := n.(*node.CollectionNode)
	if !ok {
		return Result{}, fmt.Errorf("error getting Collection resource: %s", gn.Name)
	}
	docs, ok := input.Resource.(*resource.DocstoreResource)
	if !ok {
		return Result{}, fmt.Errorf("error getting docstore resource: %s", gn.Collection.Name)
	}

	collection, cleanup, err := docs.WithCollection(gn.Collection.Name)
	if err != nil {
		return Result{}, errors.Wrapf(err, "error connecting to collection")
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
		return Result{}, fmt.Errorf("error unsupported input type: %T", input)
	}

	for _, record := range records {
		if record["id"] == nil {
			record["id"] = uuid.NewString()
		}
		err = collection.Create(context.Background(), record)
		if err != nil {
			return Result{}, errors.Wrapf(err, "error creating doc")
		}
	}

	return Result{
		Data: input.Params,
	}, nil
}

func (a *Activity) ExecuteInputNode(ctx context.Context, n node.Node, input Input) (Result, error) {
	return Result{
		Data: input.Params,
	}, nil
}

func (a *Activity) ExecuteBucketNode(ctx context.Context, n node.Node, input Input) (Result, error) {
	gn, ok := n.(*node.BucketNode)
	if !ok {
		return Result{}, fmt.Errorf("error getting Collection resource: %s", gn.Name)
	}
	bucket, ok := input.Resource.(*resource.BlobstoreResource)
	if !ok {
		return Result{}, fmt.Errorf("error getting blobstore resource: %s", gn.Bucket.Path)
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
			return Result{}, errors.Wrapf(err, "error marshaling input params")
		}
	}

	b, cleanup, err := bucket.WithPath(gn.Path)
	if err != nil {
		return Result{}, errors.Wrapf(err, fmt.Sprintf("error connecting to bucket: %s", gn.Path))
	}
	defer cleanup()

	err = b.WriteAll(context.Background(), gn.Path, bucketData, nil)
	return Result{
		Data: input.Params,
	}, nil
}

func (a *Activity) ExecuteQueryNode(ctx context.Context, n node.Node, input Input) (Result, error) {
	s, ok := n.(*node.QueryNode)
	if !ok {
		return Result{}, fmt.Errorf("error getting query resource: %s", s.Query.Collection)
	}
	docResource, ok := input.Resource.(*resource.DocstoreResource)
	if !ok {
		return Result{}, fmt.Errorf("error getting docstore resource: %s", s.Query.Collection)
	}

	d, cleanup, err := docResource.WithCollection(s.Query.Collection)
	if err != nil {
		return Result{}, errors.Wrapf(err, "error connecting to collection")
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
			return Result{}, errors.Wrapf(err, "error iterating over query results")
		}
		docs = append(docs, &doc)
	}
	return Result{
		Data: docs,
	}, nil
}
