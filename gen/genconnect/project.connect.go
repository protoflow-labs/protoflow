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
	// ProjectServiceGetProjectTypesProcedure is the fully-qualified name of the ProjectService's
	// GetProjectTypes RPC.
	ProjectServiceGetProjectTypesProcedure = "/project.ProjectService/GetProjectTypes"
	// ProjectServiceSendChatProcedure is the fully-qualified name of the ProjectService's SendChat RPC.
	ProjectServiceSendChatProcedure = "/project.ProjectService/SendChat"
	// ProjectServiceExportProjectProcedure is the fully-qualified name of the ProjectService's
	// ExportProject RPC.
	ProjectServiceExportProjectProcedure = "/project.ProjectService/ExportProject"
	// ProjectServiceLoadProjectProcedure is the fully-qualified name of the ProjectService's
	// LoadProject RPC.
	ProjectServiceLoadProjectProcedure = "/project.ProjectService/LoadProject"
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
	// ProjectServiceEnumerateProvidersProcedure is the fully-qualified name of the ProjectService's
	// EnumerateProviders RPC.
	ProjectServiceEnumerateProvidersProcedure = "/project.ProjectService/EnumerateProviders"
	// ProjectServiceGetNodeInfoProcedure is the fully-qualified name of the ProjectService's
	// GetNodeInfo RPC.
	ProjectServiceGetNodeInfoProcedure = "/project.ProjectService/GetNodeInfo"
	// ProjectServiceSaveProjectProcedure is the fully-qualified name of the ProjectService's
	// SaveProject RPC.
	ProjectServiceSaveProjectProcedure = "/project.ProjectService/SaveProject"
	// ProjectServiceRunWorkflowProcedure is the fully-qualified name of the ProjectService's
	// RunWorkflow RPC.
	ProjectServiceRunWorkflowProcedure = "/project.ProjectService/RunWorkflow"
	// ProjectServiceStopWorkflowProcedure is the fully-qualified name of the ProjectService's
	// StopWorkflow RPC.
	ProjectServiceStopWorkflowProcedure = "/project.ProjectService/StopWorkflow"
	// ProjectServiceGetWorkflowRunsProcedure is the fully-qualified name of the ProjectService's
	// GetWorkflowRuns RPC.
	ProjectServiceGetWorkflowRunsProcedure = "/project.ProjectService/GetWorkflowRuns"
)

// ProjectServiceClient is a client for the project.ProjectService service.
type ProjectServiceClient interface {
	GetProjectTypes(context.Context, *connect_go.Request[gen.GetProjectTypesRequest]) (*connect_go.Response[gen.ProjectTypes], error)
	// TODO breadchris unfortunately this is needed because of the buf fetch transport not supporting streaming
	// the suggestion is to build a custom transport that uses websockets https://github.com/bufbuild/connect-es/issues/366
	SendChat(context.Context, *connect_go.Request[gen.SendChatRequest]) (*connect_go.ServerStreamForClient[gen.SendChatResponse], error)
	ExportProject(context.Context, *connect_go.Request[gen.ExportProjectRequest]) (*connect_go.Response[gen.ExportProjectResponse], error)
	LoadProject(context.Context, *connect_go.Request[gen.LoadProjectRequest]) (*connect_go.Response[gen.LoadProjectResponse], error)
	GetProject(context.Context, *connect_go.Request[gen.GetProjectRequest]) (*connect_go.Response[gen.GetProjectResponse], error)
	GetProjects(context.Context, *connect_go.Request[gen.GetProjectsRequest]) (*connect_go.Response[gen.GetProjectsResponse], error)
	CreateProject(context.Context, *connect_go.Request[gen.CreateProjectRequest]) (*connect_go.Response[gen.CreateProjectResponse], error)
	DeleteProject(context.Context, *connect_go.Request[gen.DeleteProjectRequest]) (*connect_go.Response[gen.DeleteProjectResponse], error)
	EnumerateProviders(context.Context, *connect_go.Request[gen.GetProvidersRequest]) (*connect_go.Response[gen.GetProvidersResponse], error)
	GetNodeInfo(context.Context, *connect_go.Request[gen.GetNodeInfoRequest]) (*connect_go.Response[gen.GetNodeInfoResponse], error)
	SaveProject(context.Context, *connect_go.Request[gen.SaveProjectRequest]) (*connect_go.Response[gen.SaveProjectResponse], error)
	RunWorkflow(context.Context, *connect_go.Request[gen.RunWorkflowRequest]) (*connect_go.ServerStreamForClient[gen.NodeExecution], error)
	StopWorkflow(context.Context, *connect_go.Request[gen.StopWorkflowRequest]) (*connect_go.Response[gen.StopWorkflowResponse], error)
	GetWorkflowRuns(context.Context, *connect_go.Request[gen.GetWorkflowRunsRequest]) (*connect_go.Response[gen.GetWorkflowRunsResponse], error)
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
		getProjectTypes: connect_go.NewClient[gen.GetProjectTypesRequest, gen.ProjectTypes](
			httpClient,
			baseURL+ProjectServiceGetProjectTypesProcedure,
			opts...,
		),
		sendChat: connect_go.NewClient[gen.SendChatRequest, gen.SendChatResponse](
			httpClient,
			baseURL+ProjectServiceSendChatProcedure,
			opts...,
		),
		exportProject: connect_go.NewClient[gen.ExportProjectRequest, gen.ExportProjectResponse](
			httpClient,
			baseURL+ProjectServiceExportProjectProcedure,
			opts...,
		),
		loadProject: connect_go.NewClient[gen.LoadProjectRequest, gen.LoadProjectResponse](
			httpClient,
			baseURL+ProjectServiceLoadProjectProcedure,
			opts...,
		),
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
		enumerateProviders: connect_go.NewClient[gen.GetProvidersRequest, gen.GetProvidersResponse](
			httpClient,
			baseURL+ProjectServiceEnumerateProvidersProcedure,
			opts...,
		),
		getNodeInfo: connect_go.NewClient[gen.GetNodeInfoRequest, gen.GetNodeInfoResponse](
			httpClient,
			baseURL+ProjectServiceGetNodeInfoProcedure,
			opts...,
		),
		saveProject: connect_go.NewClient[gen.SaveProjectRequest, gen.SaveProjectResponse](
			httpClient,
			baseURL+ProjectServiceSaveProjectProcedure,
			opts...,
		),
		runWorkflow: connect_go.NewClient[gen.RunWorkflowRequest, gen.NodeExecution](
			httpClient,
			baseURL+ProjectServiceRunWorkflowProcedure,
			opts...,
		),
		stopWorkflow: connect_go.NewClient[gen.StopWorkflowRequest, gen.StopWorkflowResponse](
			httpClient,
			baseURL+ProjectServiceStopWorkflowProcedure,
			opts...,
		),
		getWorkflowRuns: connect_go.NewClient[gen.GetWorkflowRunsRequest, gen.GetWorkflowRunsResponse](
			httpClient,
			baseURL+ProjectServiceGetWorkflowRunsProcedure,
			opts...,
		),
	}
}

// projectServiceClient implements ProjectServiceClient.
type projectServiceClient struct {
	getProjectTypes    *connect_go.Client[gen.GetProjectTypesRequest, gen.ProjectTypes]
	sendChat           *connect_go.Client[gen.SendChatRequest, gen.SendChatResponse]
	exportProject      *connect_go.Client[gen.ExportProjectRequest, gen.ExportProjectResponse]
	loadProject        *connect_go.Client[gen.LoadProjectRequest, gen.LoadProjectResponse]
	getProject         *connect_go.Client[gen.GetProjectRequest, gen.GetProjectResponse]
	getProjects        *connect_go.Client[gen.GetProjectsRequest, gen.GetProjectsResponse]
	createProject      *connect_go.Client[gen.CreateProjectRequest, gen.CreateProjectResponse]
	deleteProject      *connect_go.Client[gen.DeleteProjectRequest, gen.DeleteProjectResponse]
	enumerateProviders *connect_go.Client[gen.GetProvidersRequest, gen.GetProvidersResponse]
	getNodeInfo        *connect_go.Client[gen.GetNodeInfoRequest, gen.GetNodeInfoResponse]
	saveProject        *connect_go.Client[gen.SaveProjectRequest, gen.SaveProjectResponse]
	runWorkflow        *connect_go.Client[gen.RunWorkflowRequest, gen.NodeExecution]
	stopWorkflow       *connect_go.Client[gen.StopWorkflowRequest, gen.StopWorkflowResponse]
	getWorkflowRuns    *connect_go.Client[gen.GetWorkflowRunsRequest, gen.GetWorkflowRunsResponse]
}

// GetProjectTypes calls project.ProjectService.GetProjectTypes.
func (c *projectServiceClient) GetProjectTypes(ctx context.Context, req *connect_go.Request[gen.GetProjectTypesRequest]) (*connect_go.Response[gen.ProjectTypes], error) {
	return c.getProjectTypes.CallUnary(ctx, req)
}

// SendChat calls project.ProjectService.SendChat.
func (c *projectServiceClient) SendChat(ctx context.Context, req *connect_go.Request[gen.SendChatRequest]) (*connect_go.ServerStreamForClient[gen.SendChatResponse], error) {
	return c.sendChat.CallServerStream(ctx, req)
}

// ExportProject calls project.ProjectService.ExportProject.
func (c *projectServiceClient) ExportProject(ctx context.Context, req *connect_go.Request[gen.ExportProjectRequest]) (*connect_go.Response[gen.ExportProjectResponse], error) {
	return c.exportProject.CallUnary(ctx, req)
}

// LoadProject calls project.ProjectService.LoadProject.
func (c *projectServiceClient) LoadProject(ctx context.Context, req *connect_go.Request[gen.LoadProjectRequest]) (*connect_go.Response[gen.LoadProjectResponse], error) {
	return c.loadProject.CallUnary(ctx, req)
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

// EnumerateProviders calls project.ProjectService.EnumerateProviders.
func (c *projectServiceClient) EnumerateProviders(ctx context.Context, req *connect_go.Request[gen.GetProvidersRequest]) (*connect_go.Response[gen.GetProvidersResponse], error) {
	return c.enumerateProviders.CallUnary(ctx, req)
}

// GetNodeInfo calls project.ProjectService.GetNodeInfo.
func (c *projectServiceClient) GetNodeInfo(ctx context.Context, req *connect_go.Request[gen.GetNodeInfoRequest]) (*connect_go.Response[gen.GetNodeInfoResponse], error) {
	return c.getNodeInfo.CallUnary(ctx, req)
}

// SaveProject calls project.ProjectService.SaveProject.
func (c *projectServiceClient) SaveProject(ctx context.Context, req *connect_go.Request[gen.SaveProjectRequest]) (*connect_go.Response[gen.SaveProjectResponse], error) {
	return c.saveProject.CallUnary(ctx, req)
}

// RunWorkflow calls project.ProjectService.RunWorkflow.
func (c *projectServiceClient) RunWorkflow(ctx context.Context, req *connect_go.Request[gen.RunWorkflowRequest]) (*connect_go.ServerStreamForClient[gen.NodeExecution], error) {
	return c.runWorkflow.CallServerStream(ctx, req)
}

// StopWorkflow calls project.ProjectService.StopWorkflow.
func (c *projectServiceClient) StopWorkflow(ctx context.Context, req *connect_go.Request[gen.StopWorkflowRequest]) (*connect_go.Response[gen.StopWorkflowResponse], error) {
	return c.stopWorkflow.CallUnary(ctx, req)
}

// GetWorkflowRuns calls project.ProjectService.GetWorkflowRuns.
func (c *projectServiceClient) GetWorkflowRuns(ctx context.Context, req *connect_go.Request[gen.GetWorkflowRunsRequest]) (*connect_go.Response[gen.GetWorkflowRunsResponse], error) {
	return c.getWorkflowRuns.CallUnary(ctx, req)
}

// ProjectServiceHandler is an implementation of the project.ProjectService service.
type ProjectServiceHandler interface {
	GetProjectTypes(context.Context, *connect_go.Request[gen.GetProjectTypesRequest]) (*connect_go.Response[gen.ProjectTypes], error)
	// TODO breadchris unfortunately this is needed because of the buf fetch transport not supporting streaming
	// the suggestion is to build a custom transport that uses websockets https://github.com/bufbuild/connect-es/issues/366
	SendChat(context.Context, *connect_go.Request[gen.SendChatRequest], *connect_go.ServerStream[gen.SendChatResponse]) error
	ExportProject(context.Context, *connect_go.Request[gen.ExportProjectRequest]) (*connect_go.Response[gen.ExportProjectResponse], error)
	LoadProject(context.Context, *connect_go.Request[gen.LoadProjectRequest]) (*connect_go.Response[gen.LoadProjectResponse], error)
	GetProject(context.Context, *connect_go.Request[gen.GetProjectRequest]) (*connect_go.Response[gen.GetProjectResponse], error)
	GetProjects(context.Context, *connect_go.Request[gen.GetProjectsRequest]) (*connect_go.Response[gen.GetProjectsResponse], error)
	CreateProject(context.Context, *connect_go.Request[gen.CreateProjectRequest]) (*connect_go.Response[gen.CreateProjectResponse], error)
	DeleteProject(context.Context, *connect_go.Request[gen.DeleteProjectRequest]) (*connect_go.Response[gen.DeleteProjectResponse], error)
	EnumerateProviders(context.Context, *connect_go.Request[gen.GetProvidersRequest]) (*connect_go.Response[gen.GetProvidersResponse], error)
	GetNodeInfo(context.Context, *connect_go.Request[gen.GetNodeInfoRequest]) (*connect_go.Response[gen.GetNodeInfoResponse], error)
	SaveProject(context.Context, *connect_go.Request[gen.SaveProjectRequest]) (*connect_go.Response[gen.SaveProjectResponse], error)
	RunWorkflow(context.Context, *connect_go.Request[gen.RunWorkflowRequest], *connect_go.ServerStream[gen.NodeExecution]) error
	StopWorkflow(context.Context, *connect_go.Request[gen.StopWorkflowRequest]) (*connect_go.Response[gen.StopWorkflowResponse], error)
	GetWorkflowRuns(context.Context, *connect_go.Request[gen.GetWorkflowRunsRequest]) (*connect_go.Response[gen.GetWorkflowRunsResponse], error)
}

// NewProjectServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewProjectServiceHandler(svc ProjectServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	projectServiceGetProjectTypesHandler := connect_go.NewUnaryHandler(
		ProjectServiceGetProjectTypesProcedure,
		svc.GetProjectTypes,
		opts...,
	)
	projectServiceSendChatHandler := connect_go.NewServerStreamHandler(
		ProjectServiceSendChatProcedure,
		svc.SendChat,
		opts...,
	)
	projectServiceExportProjectHandler := connect_go.NewUnaryHandler(
		ProjectServiceExportProjectProcedure,
		svc.ExportProject,
		opts...,
	)
	projectServiceLoadProjectHandler := connect_go.NewUnaryHandler(
		ProjectServiceLoadProjectProcedure,
		svc.LoadProject,
		opts...,
	)
	projectServiceGetProjectHandler := connect_go.NewUnaryHandler(
		ProjectServiceGetProjectProcedure,
		svc.GetProject,
		opts...,
	)
	projectServiceGetProjectsHandler := connect_go.NewUnaryHandler(
		ProjectServiceGetProjectsProcedure,
		svc.GetProjects,
		opts...,
	)
	projectServiceCreateProjectHandler := connect_go.NewUnaryHandler(
		ProjectServiceCreateProjectProcedure,
		svc.CreateProject,
		opts...,
	)
	projectServiceDeleteProjectHandler := connect_go.NewUnaryHandler(
		ProjectServiceDeleteProjectProcedure,
		svc.DeleteProject,
		opts...,
	)
	projectServiceEnumerateProvidersHandler := connect_go.NewUnaryHandler(
		ProjectServiceEnumerateProvidersProcedure,
		svc.EnumerateProviders,
		opts...,
	)
	projectServiceGetNodeInfoHandler := connect_go.NewUnaryHandler(
		ProjectServiceGetNodeInfoProcedure,
		svc.GetNodeInfo,
		opts...,
	)
	projectServiceSaveProjectHandler := connect_go.NewUnaryHandler(
		ProjectServiceSaveProjectProcedure,
		svc.SaveProject,
		opts...,
	)
	projectServiceRunWorkflowHandler := connect_go.NewServerStreamHandler(
		ProjectServiceRunWorkflowProcedure,
		svc.RunWorkflow,
		opts...,
	)
	projectServiceStopWorkflowHandler := connect_go.NewUnaryHandler(
		ProjectServiceStopWorkflowProcedure,
		svc.StopWorkflow,
		opts...,
	)
	projectServiceGetWorkflowRunsHandler := connect_go.NewUnaryHandler(
		ProjectServiceGetWorkflowRunsProcedure,
		svc.GetWorkflowRuns,
		opts...,
	)
	return "/project.ProjectService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ProjectServiceGetProjectTypesProcedure:
			projectServiceGetProjectTypesHandler.ServeHTTP(w, r)
		case ProjectServiceSendChatProcedure:
			projectServiceSendChatHandler.ServeHTTP(w, r)
		case ProjectServiceExportProjectProcedure:
			projectServiceExportProjectHandler.ServeHTTP(w, r)
		case ProjectServiceLoadProjectProcedure:
			projectServiceLoadProjectHandler.ServeHTTP(w, r)
		case ProjectServiceGetProjectProcedure:
			projectServiceGetProjectHandler.ServeHTTP(w, r)
		case ProjectServiceGetProjectsProcedure:
			projectServiceGetProjectsHandler.ServeHTTP(w, r)
		case ProjectServiceCreateProjectProcedure:
			projectServiceCreateProjectHandler.ServeHTTP(w, r)
		case ProjectServiceDeleteProjectProcedure:
			projectServiceDeleteProjectHandler.ServeHTTP(w, r)
		case ProjectServiceEnumerateProvidersProcedure:
			projectServiceEnumerateProvidersHandler.ServeHTTP(w, r)
		case ProjectServiceGetNodeInfoProcedure:
			projectServiceGetNodeInfoHandler.ServeHTTP(w, r)
		case ProjectServiceSaveProjectProcedure:
			projectServiceSaveProjectHandler.ServeHTTP(w, r)
		case ProjectServiceRunWorkflowProcedure:
			projectServiceRunWorkflowHandler.ServeHTTP(w, r)
		case ProjectServiceStopWorkflowProcedure:
			projectServiceStopWorkflowHandler.ServeHTTP(w, r)
		case ProjectServiceGetWorkflowRunsProcedure:
			projectServiceGetWorkflowRunsHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedProjectServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedProjectServiceHandler struct{}

func (UnimplementedProjectServiceHandler) GetProjectTypes(context.Context, *connect_go.Request[gen.GetProjectTypesRequest]) (*connect_go.Response[gen.ProjectTypes], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.GetProjectTypes is not implemented"))
}

func (UnimplementedProjectServiceHandler) SendChat(context.Context, *connect_go.Request[gen.SendChatRequest], *connect_go.ServerStream[gen.SendChatResponse]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.SendChat is not implemented"))
}

func (UnimplementedProjectServiceHandler) ExportProject(context.Context, *connect_go.Request[gen.ExportProjectRequest]) (*connect_go.Response[gen.ExportProjectResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.ExportProject is not implemented"))
}

func (UnimplementedProjectServiceHandler) LoadProject(context.Context, *connect_go.Request[gen.LoadProjectRequest]) (*connect_go.Response[gen.LoadProjectResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.LoadProject is not implemented"))
}

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

func (UnimplementedProjectServiceHandler) EnumerateProviders(context.Context, *connect_go.Request[gen.GetProvidersRequest]) (*connect_go.Response[gen.GetProvidersResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.EnumerateProviders is not implemented"))
}

func (UnimplementedProjectServiceHandler) GetNodeInfo(context.Context, *connect_go.Request[gen.GetNodeInfoRequest]) (*connect_go.Response[gen.GetNodeInfoResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.GetNodeInfo is not implemented"))
}

func (UnimplementedProjectServiceHandler) SaveProject(context.Context, *connect_go.Request[gen.SaveProjectRequest]) (*connect_go.Response[gen.SaveProjectResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.SaveProject is not implemented"))
}

func (UnimplementedProjectServiceHandler) RunWorkflow(context.Context, *connect_go.Request[gen.RunWorkflowRequest], *connect_go.ServerStream[gen.NodeExecution]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.RunWorkflow is not implemented"))
}

func (UnimplementedProjectServiceHandler) StopWorkflow(context.Context, *connect_go.Request[gen.StopWorkflowRequest]) (*connect_go.Response[gen.StopWorkflowResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.StopWorkflow is not implemented"))
}

func (UnimplementedProjectServiceHandler) GetWorkflowRuns(context.Context, *connect_go.Request[gen.GetWorkflowRunsRequest]) (*connect_go.Response[gen.GetWorkflowRunsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("project.ProjectService.GetWorkflowRuns is not implemented"))
}
