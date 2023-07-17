package project

import (
	"context"
	"encoding/json"
	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/genconnect"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	openaiclient "github.com/protoflow-labs/protoflow/pkg/openai"
	"github.com/protoflow-labs/protoflow/pkg/store"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
)

type Service struct {
	store   store.Project
	manager workflow.Manager
	cache   bucket.Bucket
	chat    *openaiclient.ChatServer
}

var ProviderSet = wire.NewSet(
	store.ProviderSet,
	workflow.ProviderSet,
	openaiclient.ChatProviderSet,
	NewService,
	wire.Bind(new(genconnect.ProjectServiceHandler), new(*Service)),
)

var _ genconnect.ProjectServiceHandler = (*Service)(nil)

func NewService(
	store store.Project,
	manager workflow.Manager,
	cache bucket.Bucket,
	chat *openaiclient.ChatServer,
) (*Service, error) {
	return &Service{
		store:   store,
		manager: manager,
		cache:   cache,
		chat:    chat,
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

func (s *Service) SendChat(ctx context.Context, c *connect.Request[gen.SendChatRequest], c2 *connect.ServerStream[gen.SendChatResponse]) error {
	obs, err := s.chat.Send(c.Msg)
	if err != nil {
		return errors.Wrapf(err, "failed to create ai chat")
	}
	msgChan := obs.Observe()
	for {
		select {
		case item := <-msgChan:
			if item.Error() {
				return errors.Wrapf(item.E, "failed to get message")
			}
			if item.V == nil {
				return nil
			}
			msg, ok := item.V.(string)
			if !ok {
				return errors.Errorf("invalid message type: %T", item.V)
			}
			if err := c2.Send(&gen.SendChatResponse{Message: msg}); err != nil {
				return errors.Wrapf(err, "failed to send message")
			}
		case <-ctx.Done():
			return nil
		}
	}
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

func (s *Service) UpdateResource(ctx context.Context, c *connect.Request[gen.UpdateResourceRequest]) (*connect.Response[gen.UpdateResourceResponse], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	var newResources []*gen.Resource
	for _, resource := range project.Resources {
		if resource.Id == c.Msg.Resource.Id {
			newResources = append(newResources, c.Msg.Resource)
			continue
		}
		newResources = append(newResources, resource)
	}
	project.Resources = newResources
	_, err = s.store.SaveProject(project)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&gen.UpdateResourceResponse{}), nil
}

func (s *Service) RunWorkflow(ctx context.Context, c *connect.Request[gen.RunWorkflowRequest], c2 *connect.ServerStream[gen.NodeExecution]) error {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	w, err := workflow.FromProject(project)
	if err != nil {
		return err
	}

	// TODO breadchris temporary for when the input is not set
	if c.Msg.Input == "" {
		c.Msg.Input = "{}"
	}

	// TODO breadchris this is a _little_ sketchy, we would like to be able to use the correct type, which might just be some data!
	var workflowInput map[string]interface{}
	err = json.Unmarshal([]byte(c.Msg.Input), &workflowInput)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal workflow input")
	}

	log.Debug().
		Str("workflow", w.ID).
		Str("node", c.Msg.NodeId).
		Msg("workflow starting")

	inputChan := make(chan rxgo.Item)
	o := rxgo.FromChannel(inputChan, rxgo.WithPublishStrategy())

	obs, err := s.manager.ExecuteWorkflow(ctx, w, c.Msg.NodeId, o)
	if err != nil {
		return err
	}

	// TODO breadchris support streaming input
	inputChan <- rx.NewItem(workflowInput)
	close(inputChan)

	var (
		obsErr error
	)
	<-obs.ForEach(func(item any) {
		log.Debug().Interface("item", item).Msg("workflow item")
		out, err := json.Marshal(item)
		if err != nil {
			obsErr = errors.Wrapf(err, "failed to marshal result data")
			return
		}

		// TODO breadchris node executions should be passed to the observable with the node id, input, and output
		err = c2.Send(&gen.NodeExecution{
			Output: string(out),
		})
		if err != nil {
			obsErr = errors.Wrapf(err, "failed to send node execution")
			return
		}
	}, func(err error) {
		obsErr = err
	}, func() {
		log.Debug().
			Str("workflow", w.ID).
			Str("node", c.Msg.NodeId).
			Msg("workflow finished")
	})
	return obsErr
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
		log.Warn().Str("node", c.Msg.NodeId).Msg("node has no info")
		return connect.NewResponse(&gen.GetNodeInfoResponse{}), nil
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
