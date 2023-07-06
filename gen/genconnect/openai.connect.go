// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: openai.proto

package genconnect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	gen "github.com/protoflow-labs/protoflow/gen"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// OpenAIServiceName is the fully-qualified name of the OpenAIService service.
	OpenAIServiceName = "openai.OpenAIService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// OpenAIServicePromptProcedure is the fully-qualified name of the OpenAIService's Prompt RPC.
	OpenAIServicePromptProcedure = "/openai.OpenAIService/Prompt"
)

// OpenAIServiceClient is a client for the openai.OpenAIService service.
type OpenAIServiceClient interface {
	Prompt(context.Context, *connect_go.Request[gen.PromptRequest]) (*connect_go.Response[gen.PromptResponse], error)
}

// NewOpenAIServiceClient constructs a client for the openai.OpenAIService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewOpenAIServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) OpenAIServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &openAIServiceClient{
		prompt: connect_go.NewClient[gen.PromptRequest, gen.PromptResponse](
			httpClient,
			baseURL+OpenAIServicePromptProcedure,
			opts...,
		),
	}
}

// openAIServiceClient implements OpenAIServiceClient.
type openAIServiceClient struct {
	prompt *connect_go.Client[gen.PromptRequest, gen.PromptResponse]
}

// Prompt calls openai.OpenAIService.Prompt.
func (c *openAIServiceClient) Prompt(ctx context.Context, req *connect_go.Request[gen.PromptRequest]) (*connect_go.Response[gen.PromptResponse], error) {
	return c.prompt.CallUnary(ctx, req)
}

// OpenAIServiceHandler is an implementation of the openai.OpenAIService service.
type OpenAIServiceHandler interface {
	Prompt(context.Context, *connect_go.Request[gen.PromptRequest]) (*connect_go.Response[gen.PromptResponse], error)
}

// NewOpenAIServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewOpenAIServiceHandler(svc OpenAIServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	openAIServicePromptHandler := connect_go.NewUnaryHandler(
		OpenAIServicePromptProcedure,
		svc.Prompt,
		opts...,
	)
	return "/openai.OpenAIService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case OpenAIServicePromptProcedure:
			openAIServicePromptHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedOpenAIServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedOpenAIServiceHandler struct{}

func (UnimplementedOpenAIServiceHandler) Prompt(context.Context, *connect_go.Request[gen.PromptRequest]) (*connect_go.Response[gen.PromptResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("openai.OpenAIService.Prompt is not implemented"))
}