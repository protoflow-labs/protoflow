// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: project.proto

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
	// ProjectServiceName is the fully-qualified name of the ProjectService service.
	ProjectServiceName = "project.ProjectService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ProjectServiceGetProjectProcedure is the fully-qualified name of the ProjectService's GetProject
	// RPC.
	ProjectServiceGetProjectProcedure = "/project.ProjectService/GetProject"
	// ProjectServiceGetProjectsProcedure is the fully-qualified name of the ProjectService's
	// GetProjects RPC.
	ProjectServiceGetProjectsProcedure = "/project.ProjectService/GetProjects"
	// ProjectServiceCreateProjectProcedure is the fully-qualified name of the ProjectService's
	// CreateProject RPC.
	ProjectServiceCreateProjectProcedure = "/project.ProjectService/CreateProject"
	// ProjectServiceDeleteProjectProcedure is the fully-qualified name of the ProjectService's
	// DeleteProject RPC.
	ProjectServiceDeleteProjectProcedure = "/project.ProjectService/DeleteProject"
	// ProjectServiceGetResourcesProcedure is the fully-qualified name of the ProjectService's
	// GetResources RPC.
	ProjectServiceGetResourcesProcedure = "/project.ProjectService/GetResources"
	// ProjectServiceSaveProjectProcedure is the fully-qualified name of the ProjectService's
	// SaveProject RPC.
	ProjectServiceSaveProjectProcedure = "/project.ProjectService/SaveProject"
	// ProjectServiceCreateResourceProcedure is the fully-qualified name of the ProjectService's
	// CreateResource RPC.
	ProjectServiceCreateResourceProcedure = "/project.ProjectService/CreateResource"
	// ProjectServiceRunWorklowProcedure is the fully-qualified name of the ProjectService's RunWorklow
	// RPC.
	ProjectServiceRunWorklowProcedure = "/project.ProjectService/RunWorklow"
	// ProjectServiceRunNodeProcedure is the fully-qualified name of the ProjectService's RunNode RPC.
	ProjectServiceRunNodeProcedure = "/project.ProjectService/RunNode"
)

// ProjectServiceClient is a client for the project.ProjectService service.
type ProjectServiceClient interface {
	GetProject(context.Context, *connect_go.Request[gen.GetProjectRequest]) (*connect_go.Response[gen.GetProjectResponse], error)
	GetProjects(context.Context, *connect_go.Request[gen.GetProjectsRequest]) (*connect_go.Response[gen.GetProjectsResponse], error)
	CreateProject(context.Context, *connect_go.Request[gen.CreateProjectRequest]) (*connect_go.Response[gen.CreateProjectResponse], error)
	DeleteProject(context.Context, *connect_go.Request[gen.DeleteProjectRequest]) (*connect_go.Response[gen.DeleteProjectResponse], error)
	GetResources(context.Context, *connect_go.Request[gen.GetResourcesRequest]) (*connect_go.Response[gen.GetResourcesResponse], error)
	SaveProject(context.Context, *connect_go.Request[gen.SaveProjectRequest]) (*connect_go.Response[gen.SaveProjectResponse], error)
	CreateResource(context.Context, *connect_go.Request[gen.CreateResourceRequest]) (*connect_go.Response[gen.CreateResourceResponse], error)
	RunWorklow(context.Context, *connect_go.Request[gen.RunWorkflowRequest]) (*connect_go.Response[gen.RunOutput], error)
	RunNode(context.Context, *connect_go.Request[gen.RunNodeRequest]) (*connect_go.Response[gen.RunOutput], error)
}

// NewProjectServiceClient constructs a client for the project.ProjectService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewProjectServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ProjectServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &projectServiceClient{
		getProject: connect_go.NewClient[gen.GetProjectRequest, gen.GetProjectResponse](
			httpClient,
			baseURL+ProjectServiceGetProjectProcedure,
			opts...,
		),
		getProjects: connect_go.NewClient[gen.GetProjectsRequest, gen.GetProjectsResponse](
			httpClient,
			baseURL+ProjectServiceGetProjectsProcedure,
			opts...,
		),
		createProject: connect_go.NewClient[gen.CreateProjectRequest, gen.CreateProjectResponse](
			httpClient,
			baseURL+ProjectServiceCreateProjectProcedure,
			opts...,
		),
		deleteProject: connect_go.NewClient[gen.DeleteProjectRequest, gen.DeleteProjectResponse](
			httpClient,
			baseURL+ProjectServiceDeleteProjectProcedure,
			opts...,
		),
		getResources: connect_go.NewClient[gen.GetResourcesRequest, gen.GetResourcesResponse](
			httpClient,
			baseURL+ProjectServiceGetResourcesProcedure,
			opts...,
		),
		saveProject: connect_go.NewClient[gen.SaveProjectRequest, gen.SaveProjectResponse](
			httpClient,
			baseURL+ProjectServiceSaveProjectProcedure,
			opts...,
		),
		createResource: connect_go.NewClient[gen.CreateResourceRequest, gen.CreateResourceResponse](
			httpClient,
			baseURL+ProjectServiceCreateResourceProcedure,
			opts...,
		),
		runWorklow: connect_go.NewClient[gen.RunWorkflowRequest, gen.RunOutput](
			httpClient,
			baseURL+ProjectServiceRunWorklowProcedure,
			opts...,
		),
		runNode: connect_go.NewClient[gen.RunNodeRequest, gen.RunOutput](
			httpClient,
			baseURL+ProjectServiceRunNodeProcedure,
			opts...,
		),
	}
}

// projectServiceClient implements ProjectServiceClient.
type projectServiceClient struct {
	getProject     *connect_go.Client[gen.GetProjectRequest, gen.GetProjectResponse]
	getProjects    *connect_go.Client[gen.GetProjectsRequest, gen.GetProjectsResponse]
	createProject  *connect_go.Client[gen.CreateProjectRequest, gen.CreateProjectResponse]
	deleteProject  *connect_go.Client[gen.DeleteProjectRequest, gen.DeleteProjectResponse]
	getResources   *connect_go.Client[gen.GetResourcesRequest, gen.GetResourcesResponse]
	saveProject    *connect_go.Client[gen.SaveProjectRequest, gen.SaveProjectResponse]
	createResource *connect_go.Client[gen.CreateResourceRequest, gen.CreateResourceResponse]
	runWorklow     *connect_go.Client[gen.RunWorkflowRequest, gen.RunOutput]
	runNode        *connect_go.Client[gen.RunNodeRequest, gen.RunOutput]
}

// GetProject calls project.ProjectService.GetProject.
func (c *projectServiceClient) GetProject(ctx context.Context, req *connect_go.Request[gen.GetProjectRequest]) (*connect_go.Response[gen.GetProjectResponse], error) {
	return c.getProject.CallUnary(ctx, req)
}

// GetProjects calls project.ProjectService.GetProjects.
func (c *projectServiceClient) GetProjects(ctx context.Context, req *connect_go.Request[gen.GetProjectsRequest]) (*connect_go.Response[gen.GetProjectsResponse], error) {
	return c.getProjects.CallUnary(ctx, req)
}

// CreateProject calls project.ProjectService.CreateProject.
func (c *projectServiceClient) CreateProject(ctx context.Context, req *connect_go.Request[gen.CreateProjectRequest]) (*connect_go.Response[gen.CreateProjectResponse], error) {
	return c.createProject.CallUnary(ctx, req)
}

// DeleteProject calls project.ProjectService.DeleteProject.
func (c *projectServiceClient) DeleteProject(ctx context.Context, req *connect_go.Request[gen.DeleteProjectRequest]) (*connect_go.Response[gen.DeleteProjectResponse], error) {
	return c.deleteProject.CallUnary(ctx, req)
}

// GetResources calls project.ProjectService.GetResources.
func (c *projectServiceClient) GetResources(ctx context.Context, req *connect_go.Request[gen.GetResourcesRequest]) (*connect_go.Response[gen.GetResourcesResponse], error) {
	return c.getResources.CallUnary(ctx, req)
}

// SaveProject calls project.ProjectService.SaveProject.
func (c *projectServiceClient) SaveProject(ctx context.Context, req *connect_go.Request[gen.SaveProjectRequest]) (*connect_go.Response[gen.SaveProjectResponse], error) {
	return c.saveProject.CallUnary(ctx, req)
}

// CreateResource calls project.ProjectService.CreateResource.
func (c *projectServiceClient) CreateResource(ctx context.Context, req *connect_go.Request[gen.CreateResourceRequest]) (*connect_go.Response[gen.CreateResourceResponse], error) {
	return c.createResource.CallUnary(ctx, req)
}

// RunWorklow calls project.ProjectService.RunWorklow.
func (c *projectServiceClient) RunWorklow(ctx context.Context, req *connect_go.Request[gen.RunWorkflowRequest]) (*connect_go.Response[gen.RunOutput], error) {
	return c.runWorklow.CallUnary(ctx, req)
}

// RunNode calls project.ProjectService.RunNode.
func (c *projectServiceClient) RunNode(ctx context.Context, req *connect_go.Request[gen.RunNodeRequest]) (*connect_go.Response[gen.RunOutput], error) {
	return c.runNode.CallUnary(ctx, req)
}

// ProjectServiceHandler is an implementation of the project.ProjectService service.
type ProjectServiceHandler interface {
	GetProject(context.Context, *connect_go.Request[gen.GetProjectRequest]) (*connect_go.Response[gen.GetProjectResponse], error)
	GetProjects(context.Context, *connect_go.Request[gen.GetProjectsRequest]) (*connect_go.Response[gen.GetProjectsResponse], error)
	CreateProject(context.Context, *connect_go.Request[gen.CreateProjectRequest]) (*connect_go.Response[gen.CreateProjectResponse], error)
	DeleteProject(context.Context, *connect_go.Request[gen.DeleteProjectRequest]) (*connect_go.Response[gen.DeleteProjectResponse], error)
	GetResources(context.Context, *connect_go.Request[gen.GetResourcesRequest]) (*connect_go.Response[gen.GetResourcesResponse], error)
	SaveProject(context.Context, *connect_go.Request[gen.SaveProjectRequest]) (*connect_go.Response[gen.SaveProjectResponse], error)
	CreateResource(context.Context, *connect_go.Request[gen.CreateResourceRequest]) (*connect_go.Response[gen.CreateResourceResponse], error)
	RunWorklow(context.Context, *connect_go.Request[gen.RunWorkflowRequest]) (*connect_go.Response[gen.RunOutput], error)
	RunNode(context.Context, *connect_go.Request[gen.RunNodeRequest]) (*connect_go.Response[gen.RunOutput], error)
}

// NewProjectServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewProjectServiceHandler(svc ProjectServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(ProjectServiceGetProjectProcedure, connect_go.NewUnaryHandler(
		ProjectServiceGetProjectProcedure,
		svc.GetProject,
		opts...,
	))
	mux.Handle(ProjectServiceGetProjectsProcedure, connect_go.NewUnaryHandler(
		ProjectServiceGetProjectsProcedure,
		svc.GetProjects,
		opts...,
	))
	mux.Handle(ProjectServiceCreateProjectProcedure, connect_go.NewUnaryHandler(
		ProjectServiceCreateProjectProcedure,
		svc.CreateProject,
		opts...,
	))
	mux.Handle(ProjectServiceDeleteProjectProcedure, connect_go.NewUnaryHandler(
		ProjectServiceDeleteProjectProcedure,
		svc.DeleteProject,
		opts...,
	))
	mux.Handle(ProjectServiceGetResourcesProcedure, connect_go.NewUnaryHandler(
		ProjectServiceGetResourcesProcedure,
		svc.GetResources,
		opts...,
	))
	mux.Handle(ProjectServiceSaveProjectProcedure, connect_go.NewUnaryHandler(
		ProjectServiceSaveProjectProcedure,
		svc.SaveProject,
		opts...,
	))
	mux.Handle(ProjectServiceCreateResourceProcedure, connect_go.NewUnaryHandler(
		ProjectServiceCreateResourceProcedure,
		svc.CreateResource,
		opts...,
	))
	mux.Handle(ProjectServiceRunWorklowProcedure, connect_go.NewUnaryHandler(
		ProjectServiceRunWorklowProcedure,
		svc.RunWorklow,
		opts...,
	))
	mux.Handle(ProjectServiceRunNodeProcedure, connect_go.NewUnaryHandler(
		ProjectServiceRunNodeProcedure,
		svc.RunNode,
		opts...,
	))
	return "/project.ProjectService/", mux
}

// UnimplementedProjectServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedProjectServiceHandler struct{}

func (UnimplementedProjectServiceHandler) GetProject(context.Context, *connect_go.Request[gen.GetProjectRequest]) (*connect_go.Response[gen.GetProjectResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.GetProject is not implemented"))
}

func (UnimplementedProjectServiceHandler) GetProjects(context.Context, *connect_go.Request[gen.GetProjectsRequest]) (*connect_go.Response[gen.GetProjectsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.GetProjects is not implemented"))
}

func (UnimplementedProjectServiceHandler) CreateProject(context.Context, *connect_go.Request[gen.CreateProjectRequest]) (*connect_go.Response[gen.CreateProjectResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.CreateProject is not implemented"))
}

func (UnimplementedProjectServiceHandler) DeleteProject(context.Context, *connect_go.Request[gen.DeleteProjectRequest]) (*connect_go.Response[gen.DeleteProjectResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.DeleteProject is not implemented"))
}

func (UnimplementedProjectServiceHandler) GetResources(context.Context, *connect_go.Request[gen.GetResourcesRequest]) (*connect_go.Response[gen.GetResourcesResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.GetResources is not implemented"))
}

func (UnimplementedProjectServiceHandler) SaveProject(context.Context, *connect_go.Request[gen.SaveProjectRequest]) (*connect_go.Response[gen.SaveProjectResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.SaveProject is not implemented"))
}

func (UnimplementedProjectServiceHandler) CreateResource(context.Context, *connect_go.Request[gen.CreateResourceRequest]) (*connect_go.Response[gen.CreateResourceResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.CreateResource is not implemented"))
}

func (UnimplementedProjectServiceHandler) RunWorklow(context.Context, *connect_go.Request[gen.RunWorkflowRequest]) (*connect_go.Response[gen.RunOutput], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.RunWorklow is not implemented"))
}

func (UnimplementedProjectServiceHandler) RunNode(context.Context, *connect_go.Request[gen.RunNodeRequest]) (*connect_go.Response[gen.RunOutput], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.RunNode is not implemented"))
}
