package project

import (
	"context"
	"encoding/json"
	"github.com/protoflow-labs/protoflow/pkg/grpc/manager"
	"github.com/rs/zerolog/log"
	"html/template"

	"github.com/protoflow-labs/protoflow/pkg/cache"
	"github.com/protoflow-labs/protoflow/pkg/grpc"

	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen/genconnect"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
	"github.com/protoflow-labs/protoflow/templates"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/gen"
)

type Service struct {
	store              Store
	manager            workflow.Manager
	blockProtoTemplate *template.Template
	cache              cache.Cache
}

var ProviderSet = wire.NewSet(
	StoreProviderSet,
	workflow.ProviderSet,
	NewService,
	wire.Bind(new(genconnect.ProjectServiceHandler), new(*Service)),
)

var _ genconnect.ProjectServiceHandler = (*Service)(nil)

func NewService(
	store Store,
	manager workflow.Manager,
	cache cache.Cache,
) (*Service, error) {
	blockProtoTemplate, err := template.New("block").ParseFS(templates.Templates, "*.template.proto")
	if err != nil {
		return nil, err
	}

	return &Service{
		store:              store,
		manager:            manager,
		blockProtoTemplate: blockProtoTemplate,
		cache:              cache,
	}, nil
}

func hydrateBlocksForResources(projectResources []*gen.Resource) ([]*gen.Resource, error) {
	var resources []*gen.Resource
	for _, resource := range projectResources {
		blocks, err := grpc.EnumerateResourceBlocks(resource)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get blocks for resource: %s", resource.Name)
		}
		resource.Blocks = blocks
		resources = append(resources, resource)
	}
	return resources, nil
}

func (s *Service) GetResources(ctx context.Context, c *connect.Request[gen.GetResourcesRequest]) (*connect.Response[gen.GetResourcesResponse], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	resources, err := hydrateBlocksForResources(project.Resources)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.GetResourcesResponse{
		Resources: resources,
	}), nil
}

func (s *Service) DeleteResource(ctx context.Context, c *connect.Request[gen.DeleteResourceRequest]) (*connect.Response[gen.DeleteResourceResponse], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	var newResources []*gen.Resource
	for _, resource := range project.Resources {
		if resource.Id == c.Msg.ResourceId {
			continue
		}
		newResources = append(newResources, resource)
	}
	project.Resources = newResources
	_, err = s.store.SaveProject(project)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&gen.DeleteResourceResponse{}), nil
}

func (s *Service) CreateResource(ctx context.Context, c *connect.Request[gen.CreateResourceRequest]) (*connect.Response[gen.CreateResourceResponse], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	resource := c.Msg.Resource
	resource.Id = uuid.New().String()

	project.Resources = append(project.Resources, resource)
	_, err = s.store.SaveProject(project)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.CreateResourceResponse{
		ResourceId: resource.Id,
	}), nil
}

func resultToAny(res *workflow.Result) (*anypb.Any, error) {
	data, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	output, err := anypb.New(&gen.Result{
		Data: data,
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (s *Service) RunWorklow(ctx context.Context, c *connect.Request[gen.RunWorkflowRequest]) (*connect.Response[gen.RunOutput], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	projectResources := workflow.ResourceMap{}

	w, err := workflow.FromProject(project, projectResources)
	if err != nil {
		return nil, err
	}

	// TODO breadchris temporary for when the input is not set
	if c.Msg.Input == "" {
		c.Msg.Input = "{}"
	}

	// TODO breadchris this is a _little_ sketchy, we would like to be able to use the correct type, which might just be some data!
	var workflowInput map[string]interface{}
	err = json.Unmarshal([]byte(c.Msg.Input), &workflowInput)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal workflow input")
	}

	log.Debug().Str("workflow", w.ID).Str("node", c.Msg.NodeId).Msg("workflow starting")
	res, err := s.manager.ExecuteWorkflowSync(ctx, w, c.Msg.NodeId, workflowInput)
	if err != nil {
		return nil, err
	}
	log.Debug().Str("workflow", w.ID).Str("node", c.Msg.NodeId).Msg("workflow finished")

	out, err := json.Marshal(res.Data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal result data")
	}

	return connect.NewResponse(&gen.RunOutput{
		Output: string(out),
	}), nil
}

func (s *Service) RunNode(ctx context.Context, c *connect.Request[gen.RunNodeRequest]) (*connect.Response[gen.RunOutput], error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetNodeInfo(ctx context.Context, c *connect.Request[gen.GetNodeInfoRequest]) (*connect.Response[gen.GetNodeInfoResponse], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	var node *gen.Node
	for _, n := range project.Graph.Nodes {
		if n.Id == c.Msg.NodeId {
			node = n
			break
		}
	}
	if node == nil {
		return nil, errors.Errorf("node %s not found", c.Msg.NodeId)
	}

	var nodeResources []*gen.Resource
	for _, r := range project.Resources {
		for _, id := range node.ResourceIds {
			if r.Id == id {
				nodeResources = append(nodeResources, r)
			}
		}
	}

	getNodeGRPCInfo := func(node *gen.GRPC) (*gen.GetNodeInfoResponse, error) {
		var grpcResource *gen.GRPCService
		for _, r := range nodeResources {
			switch r.Type.(type) {
			case *gen.Resource_LanguageService:
				grpcResource = r.GetLanguageService().Grpc
			case *gen.Resource_GrpcService:
				grpcResource = r.GetGrpcService()
			}
		}

		if grpcResource == nil {
			return nil, errors.Errorf("node %s does not have a grpc resource", c.Msg.NodeId)
		}

		// TODO breadchris I think a grpc resource should have a host that has a protocol
		m := manager.NewReflectionManager("http://" + grpcResource.Host)
		cleanup, err := m.Init()
		if err != nil {
			return nil, errors.Wrapf(err, "error initializing reflection manager")
		}
		defer cleanup()

		serviceName := node.Package + "." + node.Service
		method, err := m.ResolveMethod(serviceName, node.Method)
		if err != nil {
			return nil, errors.Wrapf(err, "error resolving method")
		}

		methodProto, err := manager.GetProtoForMethod(node.Package, node.Service, method)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting proto for method")
		}

		return &gen.GetNodeInfoResponse{
			MethodProto: methodProto,
		}, nil
	}

	var resp *gen.GetNodeInfoResponse
	// TODO breadchris we are probably going to want to reuse this logic, so there should probably be a NodeBuilder?
	switch node.Config.(type) {
	case *gen.Node_Function:
		f := node.GetFunction()
		resp, err = getNodeGRPCInfo(f.Grpc)
		if err != nil {
			return nil, err
		}
	case *gen.Node_Grpc:
		resp, err = getNodeGRPCInfo(node.GetGrpc())
		if err != nil {
			return nil, err
		}
	}
	return connect.NewResponse(resp), nil
}

func (s *Service) GetProject(context.Context, *connect.Request[gen.GetProjectRequest]) (*connect.Response[gen.GetProjectResponse], error) {
	proj, err := s.store.GetProject("local")
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.GetProjectResponse{Project: proj}), nil

}

func (s *Service) GetProjects(ctx context.Context, req *connect.Request[gen.GetProjectsRequest]) (*connect.Response[gen.GetProjectsResponse], error) {
	projects, err := s.store.ListProjects()
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.GetProjectsResponse{Projects: projects}), nil
}

func (s *Service) CreateProject(ctx context.Context, req *connect.Request[gen.CreateProjectRequest]) (*connect.Response[gen.CreateProjectResponse], error) {
	bucketDir, err := s.cache.GetFolder(".protoflow")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get bucket dir")
	}

	project := getDefaultProject(req.Msg.Name, bucketDir)

	for _, resource := range project.Resources {
		blocks, err := grpc.EnumerateResourceBlocks(resource)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create resource: %s", resource.Name)
		}
		resource.Blocks = blocks
	}

	_, err = s.store.CreateProject(&project)

	if err != nil {
		return connect.NewResponse(&gen.CreateProjectResponse{Project: nil}), nil
	}

	return connect.NewResponse(&gen.CreateProjectResponse{Project: &project}), nil
}

func (s *Service) DeleteProject(context.Context, *connect.Request[gen.DeleteProjectRequest]) (*connect.Response[gen.DeleteProjectResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("project.ProjectService.DeleteProject is not implemented"))
}

func (s *Service) SaveProject(ctx context.Context, req *connect.Request[gen.SaveProjectRequest]) (*connect.Response[gen.SaveProjectResponse], error) {
	project, err := s.store.GetProject(req.Msg.ProjectId)
	if err != nil {
		return nil, err
	}

	project.Graph = req.Msg.Graph

	if len(req.Msg.Resources) > 0 {
		project.Resources = req.Msg.Resources
	}

	_, err = s.store.SaveProject(project)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to save project %s", project.Id)
	}

	return connect.NewResponse(&gen.SaveProjectResponse{Project: project}), nil
}
