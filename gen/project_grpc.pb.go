// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: project.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ProjectServiceClient is the client API for ProjectService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProjectServiceClient interface {
	NewNode(ctx context.Context, in *NewNodeRequest, opts ...grpc.CallOption) (*NewNodeResponse, error)
	GetProjectTypes(ctx context.Context, in *GetProjectTypesRequest, opts ...grpc.CallOption) (*ProjectTypes, error)
	ExportProject(ctx context.Context, in *ExportProjectRequest, opts ...grpc.CallOption) (*ExportProjectResponse, error)
	LoadProject(ctx context.Context, in *LoadProjectRequest, opts ...grpc.CallOption) (*LoadProjectResponse, error)
	GetProject(ctx context.Context, in *GetProjectRequest, opts ...grpc.CallOption) (*GetProjectResponse, error)
	GetProjects(ctx context.Context, in *GetProjectsRequest, opts ...grpc.CallOption) (*GetProjectsResponse, error)
	CreateProject(ctx context.Context, in *CreateProjectRequest, opts ...grpc.CallOption) (*CreateProjectResponse, error)
	DeleteProject(ctx context.Context, in *DeleteProjectRequest, opts ...grpc.CallOption) (*DeleteProjectResponse, error)
	EnumerateProviders(ctx context.Context, in *GetProvidersRequest, opts ...grpc.CallOption) (*GetProvidersResponse, error)
	GetNodeInfo(ctx context.Context, in *GetNodeInfoRequest, opts ...grpc.CallOption) (*GetNodeInfoResponse, error)
	SaveProject(ctx context.Context, in *SaveProjectRequest, opts ...grpc.CallOption) (*SaveProjectResponse, error)
	RunWorkflow(ctx context.Context, in *RunWorkflowRequest, opts ...grpc.CallOption) (ProjectService_RunWorkflowClient, error)
	StopWorkflow(ctx context.Context, in *StopWorkflowRequest, opts ...grpc.CallOption) (*StopWorkflowResponse, error)
	GetWorkflowRuns(ctx context.Context, in *GetWorkflowRunsRequest, opts ...grpc.CallOption) (*GetWorkflowRunsResponse, error)
	GetRunningWorkflows(ctx context.Context, in *GetRunningWorkflowsRequest, opts ...grpc.CallOption) (*GetRunningWorkflowResponse, error)
}

type projectServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProjectServiceClient(cc grpc.ClientConnInterface) ProjectServiceClient {
	return &projectServiceClient{cc}
}

func (c *projectServiceClient) NewNode(ctx context.Context, in *NewNodeRequest, opts ...grpc.CallOption) (*NewNodeResponse, error) {
	out := new(NewNodeResponse)
	err := c.cc.Invoke(ctx, "/project.ProjectService/NewNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) GetProjectTypes(ctx context.Context, in *GetProjectTypesRequest, opts ...grpc.CallOption) (*ProjectTypes, error) {
	out := new(ProjectTypes)
	err := c.cc.Invoke(ctx, "/project.ProjectService/GetProjectTypes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) ExportProject(ctx context.Context, in *ExportProjectRequest, opts ...grpc.CallOption) (*ExportProjectResponse, error) {
	out := new(ExportProjectResponse)
	err := c.cc.Invoke(ctx, "/project.ProjectService/ExportProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) LoadProject(ctx context.Context, in *LoadProjectRequest, opts ...grpc.CallOption) (*LoadProjectResponse, error) {
	out := new(LoadProjectResponse)
	err := c.cc.Invoke(ctx, "/project.ProjectService/LoadProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) GetProject(ctx context.Context, in *GetProjectRequest, opts ...grpc.CallOption) (*GetProjectResponse, error) {
	out := new(GetProjectResponse)
	err := c.cc.Invoke(ctx, "/project.ProjectService/GetProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) GetProjects(ctx context.Context, in *GetProjectsRequest, opts ...grpc.CallOption) (*GetProjectsResponse, error) {
	out := new(GetProjectsResponse)
	err := c.cc.Invoke(ctx, "/project.ProjectService/GetProjects", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) CreateProject(ctx context.Context, in *CreateProjectRequest, opts ...grpc.CallOption) (*CreateProjectResponse, error) {
	out := new(CreateProjectResponse)
	err := c.cc.Invoke(ctx, "/project.ProjectService/CreateProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) DeleteProject(ctx context.Context, in *DeleteProjectRequest, opts ...grpc.CallOption) (*DeleteProjectResponse, error) {
	out := new(DeleteProjectResponse)
	err := c.cc.Invoke(ctx, "/project.ProjectService/DeleteProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) EnumerateProviders(ctx context.Context, in *GetProvidersRequest, opts ...grpc.CallOption) (*GetProvidersResponse, error) {
	out := new(GetProvidersResponse)
	err := c.cc.Invoke(ctx, "/project.ProjectService/EnumerateProviders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) GetNodeInfo(ctx context.Context, in *GetNodeInfoRequest, opts ...grpc.CallOption) (*GetNodeInfoResponse, error) {
	out := new(GetNodeInfoResponse)
	err := c.cc.Invoke(ctx, "/project.ProjectService/GetNodeInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) SaveProject(ctx context.Context, in *SaveProjectRequest, opts ...grpc.CallOption) (*SaveProjectResponse, error) {
	out := new(SaveProjectResponse)
	err := c.cc.Invoke(ctx, "/project.ProjectService/SaveProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) RunWorkflow(ctx context.Context, in *RunWorkflowRequest, opts ...grpc.CallOption) (ProjectService_RunWorkflowClient, error) {
	stream, err := c.cc.NewStream(ctx, &ProjectService_ServiceDesc.Streams[0], "/project.ProjectService/RunWorkflow", opts...)
	if err != nil {
		return nil, err
	}
	x := &projectServiceRunWorkflowClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ProjectService_RunWorkflowClient interface {
	Recv() (*NodeExecution, error)
	grpc.ClientStream
}

type projectServiceRunWorkflowClient struct {
	grpc.ClientStream
}

func (x *projectServiceRunWorkflowClient) Recv() (*NodeExecution, error) {
	m := new(NodeExecution)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *projectServiceClient) StopWorkflow(ctx context.Context, in *StopWorkflowRequest, opts ...grpc.CallOption) (*StopWorkflowResponse, error) {
	out := new(StopWorkflowResponse)
	err := c.cc.Invoke(ctx, "/project.ProjectService/StopWorkflow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) GetWorkflowRuns(ctx context.Context, in *GetWorkflowRunsRequest, opts ...grpc.CallOption) (*GetWorkflowRunsResponse, error) {
	out := new(GetWorkflowRunsResponse)
	err := c.cc.Invoke(ctx, "/project.ProjectService/GetWorkflowRuns", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) GetRunningWorkflows(ctx context.Context, in *GetRunningWorkflowsRequest, opts ...grpc.CallOption) (*GetRunningWorkflowResponse, error) {
	out := new(GetRunningWorkflowResponse)
	err := c.cc.Invoke(ctx, "/project.ProjectService/GetRunningWorkflows", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProjectServiceServer is the server API for ProjectService service.
// All implementations should embed UnimplementedProjectServiceServer
// for forward compatibility
type ProjectServiceServer interface {
	NewNode(context.Context, *NewNodeRequest) (*NewNodeResponse, error)
	GetProjectTypes(context.Context, *GetProjectTypesRequest) (*ProjectTypes, error)
	ExportProject(context.Context, *ExportProjectRequest) (*ExportProjectResponse, error)
	LoadProject(context.Context, *LoadProjectRequest) (*LoadProjectResponse, error)
	GetProject(context.Context, *GetProjectRequest) (*GetProjectResponse, error)
	GetProjects(context.Context, *GetProjectsRequest) (*GetProjectsResponse, error)
	CreateProject(context.Context, *CreateProjectRequest) (*CreateProjectResponse, error)
	DeleteProject(context.Context, *DeleteProjectRequest) (*DeleteProjectResponse, error)
	EnumerateProviders(context.Context, *GetProvidersRequest) (*GetProvidersResponse, error)
	GetNodeInfo(context.Context, *GetNodeInfoRequest) (*GetNodeInfoResponse, error)
	SaveProject(context.Context, *SaveProjectRequest) (*SaveProjectResponse, error)
	RunWorkflow(*RunWorkflowRequest, ProjectService_RunWorkflowServer) error
	StopWorkflow(context.Context, *StopWorkflowRequest) (*StopWorkflowResponse, error)
	GetWorkflowRuns(context.Context, *GetWorkflowRunsRequest) (*GetWorkflowRunsResponse, error)
	GetRunningWorkflows(context.Context, *GetRunningWorkflowsRequest) (*GetRunningWorkflowResponse, error)
}

// UnimplementedProjectServiceServer should be embedded to have forward compatible implementations.
type UnimplementedProjectServiceServer struct {
}

func (UnimplementedProjectServiceServer) NewNode(context.Context, *NewNodeRequest) (*NewNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewNode not implemented")
}
func (UnimplementedProjectServiceServer) GetProjectTypes(context.Context, *GetProjectTypesRequest) (*ProjectTypes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProjectTypes not implemented")
}
func (UnimplementedProjectServiceServer) ExportProject(context.Context, *ExportProjectRequest) (*ExportProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExportProject not implemented")
}
func (UnimplementedProjectServiceServer) LoadProject(context.Context, *LoadProjectRequest) (*LoadProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadProject not implemented")
}
func (UnimplementedProjectServiceServer) GetProject(context.Context, *GetProjectRequest) (*GetProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProject not implemented")
}
func (UnimplementedProjectServiceServer) GetProjects(context.Context, *GetProjectsRequest) (*GetProjectsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProjects not implemented")
}
func (UnimplementedProjectServiceServer) CreateProject(context.Context, *CreateProjectRequest) (*CreateProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProject not implemented")
}
func (UnimplementedProjectServiceServer) DeleteProject(context.Context, *DeleteProjectRequest) (*DeleteProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProject not implemented")
}
func (UnimplementedProjectServiceServer) EnumerateProviders(context.Context, *GetProvidersRequest) (*GetProvidersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnumerateProviders not implemented")
}
func (UnimplementedProjectServiceServer) GetNodeInfo(context.Context, *GetNodeInfoRequest) (*GetNodeInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNodeInfo not implemented")
}
func (UnimplementedProjectServiceServer) SaveProject(context.Context, *SaveProjectRequest) (*SaveProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveProject not implemented")
}
func (UnimplementedProjectServiceServer) RunWorkflow(*RunWorkflowRequest, ProjectService_RunWorkflowServer) error {
	return status.Errorf(codes.Unimplemented, "method RunWorkflow not implemented")
}
func (UnimplementedProjectServiceServer) StopWorkflow(context.Context, *StopWorkflowRequest) (*StopWorkflowResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopWorkflow not implemented")
}
func (UnimplementedProjectServiceServer) GetWorkflowRuns(context.Context, *GetWorkflowRunsRequest) (*GetWorkflowRunsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWorkflowRuns not implemented")
}
func (UnimplementedProjectServiceServer) GetRunningWorkflows(context.Context, *GetRunningWorkflowsRequest) (*GetRunningWorkflowResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRunningWorkflows not implemented")
}

// UnsafeProjectServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProjectServiceServer will
// result in compilation errors.
type UnsafeProjectServiceServer interface {
	mustEmbedUnimplementedProjectServiceServer()
}

func RegisterProjectServiceServer(s grpc.ServiceRegistrar, srv ProjectServiceServer) {
	s.RegisterService(&ProjectService_ServiceDesc, srv)
}

func _ProjectService_NewNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).NewNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.ProjectService/NewNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).NewNode(ctx, req.(*NewNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_GetProjectTypes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProjectTypesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).GetProjectTypes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.ProjectService/GetProjectTypes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).GetProjectTypes(ctx, req.(*GetProjectTypesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_ExportProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).ExportProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.ProjectService/ExportProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).ExportProject(ctx, req.(*ExportProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_LoadProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).LoadProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.ProjectService/LoadProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).LoadProject(ctx, req.(*LoadProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_GetProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).GetProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.ProjectService/GetProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).GetProject(ctx, req.(*GetProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_GetProjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProjectsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).GetProjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.ProjectService/GetProjects",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).GetProjects(ctx, req.(*GetProjectsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_CreateProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).CreateProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.ProjectService/CreateProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).CreateProject(ctx, req.(*CreateProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_DeleteProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).DeleteProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.ProjectService/DeleteProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).DeleteProject(ctx, req.(*DeleteProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_EnumerateProviders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProvidersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).EnumerateProviders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.ProjectService/EnumerateProviders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).EnumerateProviders(ctx, req.(*GetProvidersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_GetNodeInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNodeInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).GetNodeInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.ProjectService/GetNodeInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).GetNodeInfo(ctx, req.(*GetNodeInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_SaveProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).SaveProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.ProjectService/SaveProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).SaveProject(ctx, req.(*SaveProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_RunWorkflow_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RunWorkflowRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ProjectServiceServer).RunWorkflow(m, &projectServiceRunWorkflowServer{stream})
}

type ProjectService_RunWorkflowServer interface {
	Send(*NodeExecution) error
	grpc.ServerStream
}

type projectServiceRunWorkflowServer struct {
	grpc.ServerStream
}

func (x *projectServiceRunWorkflowServer) Send(m *NodeExecution) error {
	return x.ServerStream.SendMsg(m)
}

func _ProjectService_StopWorkflow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopWorkflowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).StopWorkflow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.ProjectService/StopWorkflow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).StopWorkflow(ctx, req.(*StopWorkflowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_GetWorkflowRuns_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWorkflowRunsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).GetWorkflowRuns(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.ProjectService/GetWorkflowRuns",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).GetWorkflowRuns(ctx, req.(*GetWorkflowRunsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_GetRunningWorkflows_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRunningWorkflowsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).GetRunningWorkflows(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.ProjectService/GetRunningWorkflows",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).GetRunningWorkflows(ctx, req.(*GetRunningWorkflowsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProjectService_ServiceDesc is the grpc.ServiceDesc for ProjectService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProjectService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "project.ProjectService",
	HandlerType: (*ProjectServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewNode",
			Handler:    _ProjectService_NewNode_Handler,
		},
		{
			MethodName: "GetProjectTypes",
			Handler:    _ProjectService_GetProjectTypes_Handler,
		},
		{
			MethodName: "ExportProject",
			Handler:    _ProjectService_ExportProject_Handler,
		},
		{
			MethodName: "LoadProject",
			Handler:    _ProjectService_LoadProject_Handler,
		},
		{
			MethodName: "GetProject",
			Handler:    _ProjectService_GetProject_Handler,
		},
		{
			MethodName: "GetProjects",
			Handler:    _ProjectService_GetProjects_Handler,
		},
		{
			MethodName: "CreateProject",
			Handler:    _ProjectService_CreateProject_Handler,
		},
		{
			MethodName: "DeleteProject",
			Handler:    _ProjectService_DeleteProject_Handler,
		},
		{
			MethodName: "EnumerateProviders",
			Handler:    _ProjectService_EnumerateProviders_Handler,
		},
		{
			MethodName: "GetNodeInfo",
			Handler:    _ProjectService_GetNodeInfo_Handler,
		},
		{
			MethodName: "SaveProject",
			Handler:    _ProjectService_SaveProject_Handler,
		},
		{
			MethodName: "StopWorkflow",
			Handler:    _ProjectService_StopWorkflow_Handler,
		},
		{
			MethodName: "GetWorkflowRuns",
			Handler:    _ProjectService_GetWorkflowRuns_Handler,
		},
		{
			MethodName: "GetRunningWorkflows",
			Handler:    _ProjectService_GetRunningWorkflows_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RunWorkflow",
			Handler:       _ProjectService_RunWorkflow_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "project.proto",
}
