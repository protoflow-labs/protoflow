package project

import (
	"context"
	"encoding/json"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	store "github.com/protoflow-labs/protoflow/pkg/store"
	"github.com/rs/zerolog/log"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/genconnect"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
)

type Service struct {
	store   store.Project
	manager workflow.Manager
	cache   bucket.Bucket
}

var ProviderSet = wire.NewSet(
	store.ProviderSet,
	workflow.ProviderSet,
	NewService,
	wire.Bind(new(genconnect.ProjectServiceHandler), new(*Service)),
)

var _ genconnect.ProjectServiceHandler = (*Service)(nil)

func NewService(
	store store.Project,
	manager workflow.Manager,
	cache bucket.Bucket,
) (*Service, error) {
	return &Service{
		store:   store,
		manager: manager,
		cache:   cache,
	}, nil
}

func hydrateBlocksForResources(projectResources []*gen.Resource) ([]*gen.EnumeratedResource, error) {
	var resources []*gen.EnumeratedResource
	for _, resource := range projectResources {
		info := &gen.ResourceInfo{
			State: gen.ResourceState_READY,
			Error: "",
		}
		nodes, err := grpc.EnumerateResourceBlocks(resource)
		if err != nil {
			info.State = gen.ResourceState_ERROR
			info.Error = err.Error()
		}
		resources = append(resources, &gen.EnumeratedResource{
			Resource: resource,
			Nodes:    nodes,
			Info:     info,
		})
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

	r := c.Msg.Resource
	r.Id = uuid.New().String()

	project.Resources = append(project.Resources, r)
	_, err = s.store.SaveProject(project)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.CreateResourceResponse{
		ResourceId: r.Id,
	}), nil
}

func (s *Service) RunWorklow(ctx context.Context, c *connect.Request[gen.RunWorkflowRequest]) (*connect.Response[gen.RunOutput], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	w, err := workflow.FromProject(project)
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

	out, err := json.Marshal(res)
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
	w, err := workflow.FromProject(project)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get workflow from project")
	}
	n, err := w.GetNode(c.Msg.NodeId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %s", c.Msg.NodeId)
	}
	nodeInfo, err := w.GetNodeInfo(n)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node info for node %s", c.Msg.NodeId)
	}
	if nodeInfo == nil {
		return nil, errors.Errorf("node %s has no info", c.Msg.NodeId)
	}
	typeInfo, err := nodeInfo.Method.Proto()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get proto for node %s", c.Msg.NodeId)
	}
	proto, err := nodeInfo.BuildProto()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build proto for node %s", c.Msg.NodeId)
	}
	return connect.NewResponse(&gen.GetNodeInfoResponse{
		MethodProto: proto,
		TypeInfo:    typeInfo,
	}), nil
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
	// TODO breadchris this folder should be configurable
	bucketDir, err := s.cache.GetFolder(".protoflow")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get bucket dir")
	}

	project := getDefaultProject(req.Msg.Name, bucketDir)

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
		return nil, errors.Wrapf(err, "failed to get project %s", req.Msg.ProjectId)
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

func (s *Service) GetWorkflowRuns(ctx context.Context, c *connect.Request[gen.GetWorkflowRunsRequest]) (*connect.Response[gen.GetWorkflowRunsResponse], error) {
	runs, err := s.store.GetWorkflowRunsForProject(c.Msg.ProjectId)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&gen.GetWorkflowRunsResponse{Runs: runs}), nil
}
