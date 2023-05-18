package workflow

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	grpcanal "github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"github.com/rs/zerolog/log"
	"io"
	"net/url"
)

type Activity struct{}

func getInvokeOptions(base, serviceName, methodName string, outputStream bufcurl.OutputStream) grpcanal.InvokeOptions {
	grpcURL := fmt.Sprintf("%s/%s/%s", base, serviceName, methodName)
	return grpcanal.InvokeOptions{
		OutputStream:          outputStream,
		TLSConfig:             bufcurl.TLSSettings{},
		URL:                   grpcURL,
		Protocol:              "grpc",
		Headers:               nil,
		UserAgent:             "",
		ReflectProtocol:       "grpc-v1",
		ReflectHeaders:        nil,
		UnixSocket:            "",
		HTTP2PriorKnowledge:   true,
		NoKeepAlive:           false,
		KeepAliveTimeSeconds:  0,
		ConnectTimeoutSeconds: 0,
	}
}

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
func (a *Activity) ExecuteGRPCNode(ctx context.Context, node *GRPCNode, input Input) (Result, error) {
	log.Info().Msgf("executing node: %s", node.Service)

	g, ok := input.Resources[GRPCResourceType].(*GRPCResource)
	if !ok {
		return Result{}, fmt.Errorf("error getting GRPC resource: %s.%s", node.Service, node.Method)
	}

	serviceName := node.Service
	if node.Package != "" {
		serviceName = node.Package + "." + serviceName
	}

	host, err := formatHost(g.Host)
	if err != nil {
		return Result{}, errors.Wrapf(err, "error formatting host: %s", g.Host)
	}

	inputStream := bufcurl.NewMemoryInputStream()
	outputStream := bufcurl.NewMemoryOutputStream()
	go func() {
		// TODO breadchris we are relying on this grpc call to close the output stream. How can the stream be closed by the caller?
		defer outputStream.Close()
		err := grpcanal.ExecuteCurl(ctx, getInvokeOptions(host, serviceName, node.Method, outputStream), inputStream)
		if err != nil {
			outputStream.Error(err)
		}
	}()
	go func() {
		inputStream.Push(input.Params)
		inputStream.Close()
	}()

	var data any
	for {
		output, err := outputStream.Next()
		if err != nil {
			if err != io.EOF {
				return Result{}, errors.Wrapf(err, "error reading output stream")
			}
			break
		}

		// TODO breadchris whatever the last output is, is the data. Streaming is not supported yet.
		data = output
	}
	return Result{
		Data: data,
	}, nil
}

func (a *Activity) ExecuteRestNode(ctx context.Context, node *RESTNode, input Input) (Result, error) {
	log.Debug().
		Interface("headers", node.Headers).
		Str("method", node.Method).
		Str("path", node.Path).
		Msgf("executing rest")
	res, err := util.InvokeMethodOnUrl(node.Method, node.Path, node.Headers, input.Params)
	if err != nil {
		return Result{}, errors.Wrapf(err, "error invoking method: %s", node.Method)
	}
	return Result{
		Data: res,
	}, nil
}

func (a *Activity) ExecuteFunctionNode(ctx context.Context, node *FunctionNode, input Input) (Result, error) {
	log.Debug().Msgf("executing function node: %s.%s", node.Function.Runtime, node.Name)
	g, ok := input.Resources[LanguageServiceType].(*LanguageServiceResource)
	if !ok {
		return Result{}, fmt.Errorf("error getting language service resource: %s.%s", node.Function.Runtime, node.Name)
	}

	// provide the grpc resource to the grpc node call. Is this the best place for this? Should this be provided on injection? Probably.
	input.Resources[GRPCResourceType] = g.GRPCResource

	// TODO breadchris how is the method name formatted?
	serviceName := node.Function.Runtime + "Service"
	methodName := util.ToTitleCase(node.Name)

	grpcNode := &GRPCNode{
		GRPC: &gen.GRPC{
			Package: "protoflow",
			Service: serviceName,
			Method:  methodName,
		},
	}
	return a.ExecuteGRPCNode(ctx, grpcNode, input)
}
