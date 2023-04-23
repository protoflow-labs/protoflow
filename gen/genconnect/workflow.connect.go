// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: workflow.proto

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
	// ManagerServiceName is the fully-qualified name of the ManagerService service.
	ManagerServiceName = "workflow.ManagerService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ManagerServiceCreateWorkflowProcedure is the fully-qualified name of the ManagerService's
	// CreateWorkflow RPC.
	ManagerServiceCreateWorkflowProcedure = "/workflow.ManagerService/CreateWorkflow"
	// ManagerServiceStartWorkflowProcedure is the fully-qualified name of the ManagerService's
	// StartWorkflow RPC.
	ManagerServiceStartWorkflowProcedure = "/workflow.ManagerService/StartWorkflow"
)

// ManagerServiceClient is a client for the workflow.ManagerService service.
type ManagerServiceClient interface {
	CreateWorkflow(context.Context, *connect_go.Request[gen.Workflow]) (*connect_go.Response[gen.ID], error)
	StartWorkflow(context.Context, *connect_go.Request[gen.WorkflowEntrypoint]) (*connect_go.Response[gen.Run], error)
}

// NewManagerServiceClient constructs a client for the workflow.ManagerService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewManagerServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ManagerServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &managerServiceClient{
		createWorkflow: connect_go.NewClient[gen.Workflow, gen.ID](
			httpClient,
			baseURL+ManagerServiceCreateWorkflowProcedure,
			opts...,
		),
		startWorkflow: connect_go.NewClient[gen.WorkflowEntrypoint, gen.Run](
			httpClient,
			baseURL+ManagerServiceStartWorkflowProcedure,
			opts...,
		),
	}
}

// managerServiceClient implements ManagerServiceClient.
type managerServiceClient struct {
	createWorkflow *connect_go.Client[gen.Workflow, gen.ID]
	startWorkflow  *connect_go.Client[gen.WorkflowEntrypoint, gen.Run]
}

// CreateWorkflow calls workflow.ManagerService.CreateWorkflow.
func (c *managerServiceClient) CreateWorkflow(ctx context.Context, req *connect_go.Request[gen.Workflow]) (*connect_go.Response[gen.ID], error) {
	return c.createWorkflow.CallUnary(ctx, req)
}

// StartWorkflow calls workflow.ManagerService.StartWorkflow.
func (c *managerServiceClient) StartWorkflow(ctx context.Context, req *connect_go.Request[gen.WorkflowEntrypoint]) (*connect_go.Response[gen.Run], error) {
	return c.startWorkflow.CallUnary(ctx, req)
}

// ManagerServiceHandler is an implementation of the workflow.ManagerService service.
type ManagerServiceHandler interface {
	CreateWorkflow(context.Context, *connect_go.Request[gen.Workflow]) (*connect_go.Response[gen.ID], error)
	StartWorkflow(context.Context, *connect_go.Request[gen.WorkflowEntrypoint]) (*connect_go.Response[gen.Run], error)
}

// NewManagerServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewManagerServiceHandler(svc ManagerServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(ManagerServiceCreateWorkflowProcedure, connect_go.NewUnaryHandler(
		ManagerServiceCreateWorkflowProcedure,
		svc.CreateWorkflow,
		opts...,
	))
	mux.Handle(ManagerServiceStartWorkflowProcedure, connect_go.NewUnaryHandler(
		ManagerServiceStartWorkflowProcedure,
		svc.StartWorkflow,
		opts...,
	))
	return "/workflow.ManagerService/", mux
}

// UnimplementedManagerServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedManagerServiceHandler struct{}

func (UnimplementedManagerServiceHandler) CreateWorkflow(context.Context, *connect_go.Request[gen.Workflow]) (*connect_go.Response[gen.ID], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("workflow.ManagerService.CreateWorkflow is not implemented"))
}

func (UnimplementedManagerServiceHandler) StartWorkflow(context.Context, *connect_go.Request[gen.WorkflowEntrypoint]) (*connect_go.Response[gen.Run], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("workflow.ManagerService.StartWorkflow is not implemented"))
}
