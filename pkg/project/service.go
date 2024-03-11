package project

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/genconnect"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/graph/edge"
	ngrpc "github.com/protoflow-labs/protoflow/pkg/graph/node/grpc"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/llm"
	"github.com/protoflow-labs/protoflow/pkg/llm/schemas"
	"github.com/protoflow-labs/protoflow/pkg/protobuf"
	"github.com/protoflow-labs/protoflow/pkg/store"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
	"github.com/rs/zerolog/log"
)

type Service struct {
	store          store.Project
	cache          bucket.Bucket
	defaultProject *gen.Project
	manager        *workflow.ManagerBuilder
	// TODO breadchris rename this to something that is more relevant
	workflowManager *workflow.WorkflowManager
	llm             *llm.Agent
}

var ProviderSet = wire.NewSet(
	store.ProviderSet,
	NewService,
	workflow.ProviderSet,
	llm.ProviderSet,
	wire.Bind(new(genconnect.ProjectServiceHandler), new(*Service)),
)

var _ genconnect.ProjectServiceHandler = (*Service)(nil)

func NewDefaultProject(cache bucket.Bucket) (*gen.Project, error) {
	// TODO breadchris this folder should be configurable
	bucketDir, err := cache.GetFolder("filestore")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get bucket dir")
	}
	return getDefaultProject("local", bucketDir), nil
}

func NewService(
	store store.Project,
	cache bucket.Bucket,
	defaultProject *gen.Project,
	manager *workflow.ManagerBuilder,
	workflowManager *workflow.WorkflowManager,
	llm *llm.Agent,
) (*Service, error) {
	return &Service{
		store:           store,
		cache:           cache,
		defaultProject:  defaultProject,
		manager:         manager,
		workflowManager: workflowManager,
		llm:             llm,
	}, nil
}

func (s *Service) EnumerateProviders(ctx context.Context, c *connect.Request[gen.GetProvidersRequest]) (*connect.Response[gen.GetProvidersResponse], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	w, err := FromProto(project)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	return connect.NewResponse(&gen.GetProvidersResponse{
		Providers: w.EnumerateProviders(),
	}), nil
}

func getProjectTypes() (*gen.ProjectTypes, error) {
	// TODO breadchris when types are bound to a project, this should be specific to a project
	// return the rules for different layers
	n := &gen.Node{}
	nd, err := grpc.SerializeType(n)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to wrap message")
	}
	e := &gen.Edge{}
	ed, err := grpc.SerializeType(e)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to wrap message")
	}

	// TODO breadchris cleanup this code, see blocks.go:76
	tr := grpc.NewTypeResolver()
	tr = tr.ResolveLookup(n)
	tr = tr.ResolveLookup(e)

	sr := tr.Serialize()

	return &gen.ProjectTypes{
		NodeType:   nd.AsDescriptorProto(),
		EdgeType:   ed.AsDescriptorProto(),
		DescLookup: sr.DescLookup,
		EnumLookup: sr.EnumLookup,
	}, nil
}

func maybeDefaultPath(path string) string {
	if path == "" {
		return "protoflow.bin"
	}
	return path
}

func (s *Service) ExportProject(ctx context.Context, c *connect.Request[gen.ExportProjectRequest]) (*connect.Response[gen.ExportProjectResponse], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	return connect.NewResponse(
		&gen.ExportProjectResponse{},
	), SaveToFile(project, maybeDefaultPath(c.Msg.Path))
}

func (s *Service) LoadProject(ctx context.Context, c *connect.Request[gen.LoadProjectRequest]) (*connect.Response[gen.LoadProjectResponse], error) {
	project, err := LoadFromFile(maybeDefaultPath(c.Msg.Path))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load project from %s", c.Msg.Path)
	}
	_, err = s.store.SaveProject(project)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to save project %s", project.Id)
	}
	return connect.NewResponse(
		&gen.LoadProjectResponse{
			Project: project,
		},
	), nil
}

func (s *Service) GetProjectTypes(ctx context.Context, c *connect.Request[gen.GetProjectTypesRequest]) (*connect.Response[gen.ProjectTypes], error) {
	projectTypes, err := getProjectTypes()
	return &connect.Response[gen.ProjectTypes]{
		Msg: projectTypes,
	}, err
}

func (s *Service) AddMethod(ctx context.Context, c *connect.Request[gen.AddMethodRequest]) (*connect.Response[gen.AddMethodResponse], error) {
	// TODO breadchris this should be configurable
	dir := "./proto"

	err := protobuf.AddMethod(dir, c.Msg.File, c.Msg.Package, c.Msg.Service, c.Msg.Method)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add method %s to proto %s", c.Msg.Method, c.Msg.File)
	}
	return connect.NewResponse(
		&gen.AddMethodResponse{},
	), nil
}

func (s *Service) RunGRPCMethod(ctx context.Context, c *connect.Request[gen.RunGRPCMethodRequest], c2 *connect.ServerStream[gen.NodeExecution]) error {
	pid := uuid.NewString()
	server := ngrpc.NewProto("server", ngrpc.NewServerProto(c.Msg.Host))
	method := ngrpc.NewProto("method", ngrpc.NewMethodProto(c.Msg.Package, c.Msg.Service, c.Msg.Method))
	p := &gen.Project{
		Id:   pid,
		Name: "run-grpc-method",
		Graph: &gen.Graph{
			Nodes: []*gen.Node{server, method},
			Edges: []*gen.Edge{edge.NewProvidesProto(server.Id, method.Id)},
		},
	}

	w, err := FromProto(p)
	if err != nil {
		return err
	}

	// TODO breadchris this is a _little_ sketchy, we would like to be able to use the correct type, which might just be some data!
	var workflowInput map[string]any
	err = json.Unmarshal([]byte(c.Msg.Input), &workflowInput)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal workflow input")
	}

	io, err := s.wireWorkflow(ctx, w, method.Id, workflowInput, nil, &gen.RunWorkflowRequest{})
	if err != nil {
		return errors.Wrapf(err, "failed to start workflow")
	}
	obs := io.Observable
	log.Debug().Msg("done wiring workflows")

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

		// TODO breadchris node executions should be passed to the observable with the node wID, input, and output
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
			Str("node", method.Id).
			Msg("workflow finished")
		if err != nil {
			obsErr = errors.Wrapf(err, "failed to stop workflow")
			return
		}
	})
	if obsErr != nil {
		log.Error().Err(obsErr).Msg("workflow error")
	}
	return obsErr
}

func (s *Service) GetGRPCServerInfo(ctx context.Context, c *connect.Request[gen.GetGRPCServerInfoRequest]) (*connect.Response[gen.GetGRPCServerInfoResponse], error) {
	services, err := grpc.GetGRPCTypeInfo(c.Msg.Host)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to enumerate resource blocks")
	}
	return connect.NewResponse(&gen.GetGRPCServerInfoResponse{
		Services: services,
	}), nil
}

func (s *Service) GetNodeInfo(ctx context.Context, c *connect.Request[gen.GetNodeInfoRequest]) (*connect.Response[gen.GetNodeInfoResponse], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	w, err := FromProto(project)
	if err != nil {
		return nil, err
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
	if nodeInfo.Method == nil {
		log.Warn().Str("node", c.Msg.NodeId).Msg("node has no method")
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

func (s *Service) GenerateAIStub(ctx context.Context, c *connect.Request[gen.GenerateAIStubRequest]) (*connect.Response[gen.GenerateCode], error) {
	t := &gen.GenerateCode{}

	jd, err := schemas.Schemas.ReadFile(fmt.Sprintf("%s.json", t.ProtoReflect().Descriptor().Name()))
	if err != nil {
		return nil, err
	}

	basePrompt := fmt.Sprintf("Generate a %s function that %s", c.Msg.Language, c.Msg.Description)

	err = s.llm.PromptToProto(ctx, basePrompt, t, jd)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to prompt to proto")
	}
	return connect.NewResponse(t), nil
}
