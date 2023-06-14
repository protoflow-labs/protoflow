package node

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl"
	grpcanal "github.com/protoflow-labs/protoflow/pkg/grpc/manager"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"github.com/protoflow-labs/protoflow/pkg/workflow/execute"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/rs/zerolog/log"
	"io"
	"net/url"
)

type Activity struct{}

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
func (a *Activity) ExecuteGRPCNode(ctx context.Context, node *GRPCNode, input execute.Input) (execute.Result, error) {
	log.Info().Msgf("executing node: %s", node.Service)

	g, ok := input.Resource.(*resource.GRPCResource)
	if !ok {
		return execute.Result{}, fmt.Errorf("error getting GRPC resource: %s.%s", node.Service, node.Method)
	}

	serviceName := node.Service
	if node.Package != "" {
		serviceName = node.Package + "." + serviceName
	}

	host, err := formatHost(g.Host)
	if err != nil {
		return execute.Result{}, errors.Wrapf(err, "error formatting host: %s", g.Host)
	}

	inputStream := bufcurl.NewMemoryInputStream()
	outputStream := bufcurl.NewMemoryOutputStream()

	manager := grpcanal.NewReflectionManager(host)
	cleanup, err := manager.Init()
	if err != nil {
		return execute.Result{}, errors.Wrapf(err, "error initializing reflection manager")
	}
	defer cleanup()

	method, err := manager.ResolveMethod(serviceName, node.Method)
	if err != nil {
		return execute.Result{}, errors.Wrapf(err, "error resolving method: %s.%s", serviceName, node.Method)
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
				inputStream.Push(d)
			}
		} else {
			// TODO breadchris figure out how backpacks work
			err = util.SetInterfaceField(input.Params, BackpackKey, input.Backpack)
			if err != nil {
				outputStream.Error(errors.Wrapf(err, "error setting backpack field"))
				return
			}
			inputStream.Push(input.Params)
		}
	}()

	var res execute.Result
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
				return execute.Result{}, errors.Wrapf(err, "error reading output stream")
			}
			return execute.Result{}, errors.Wrapf(err, "did not receive any output for unary call")
		}
		res.Data = output
	}
	return res, nil
}

func (a *Activity) ExecuteRestNode(ctx context.Context, node *RESTNode, input execute.Input) (execute.Result, error) {
	log.Debug().
		Interface("headers", node.Headers).
		Str("method", node.Method).
		Str("path", node.Path).
		Msgf("executing rest")
	res, err := util.InvokeMethodOnUrl(node.Method, node.Path, node.Headers, input.Params)
	if err != nil {
		return execute.Result{}, errors.Wrapf(err, "error invoking method: %s", node.Method)
	}
	return execute.Result{
		Data: res,
	}, nil
}

func (a *Activity) ExecuteFunctionNode(ctx context.Context, node *FunctionNode, input execute.Input) (execute.Result, error) {
	log.Debug().Msgf("executing function node: %s", node.Name)
	g, ok := input.Resource.(*resource.LanguageServiceResource)
	if !ok {
		return execute.Result{}, fmt.Errorf("error getting language service resource: %s", node.Name)
	}

	// provide the grpc resource to the grpc node call. Is this the best place for this? Should this be provided on injection? Probably.
	input.Resource = g.GRPCResource

	grpcNode := node.ToGRPC(g)
	return a.ExecuteGRPCNode(ctx, grpcNode, input)
}
