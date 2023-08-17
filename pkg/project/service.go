package project

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/code"
	pdata "github.com/protoflow-labs/protoflow/gen/data"
	"github.com/protoflow-labs/protoflow/gen/genconnect"
	pgrpc "github.com/protoflow-labs/protoflow/gen/grpc"
	preason "github.com/protoflow-labs/protoflow/gen/reason"
	"github.com/protoflow-labs/protoflow/gen/storage"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/data"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/reason"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	openaiclient "github.com/protoflow-labs/protoflow/pkg/openai"
	"github.com/protoflow-labs/protoflow/pkg/store"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
)

type Service struct {
	store          store.Project
	manager        workflow.Manager
	cache          bucket.Bucket
	chat           *openaiclient.ChatServer
	defaultProject *gen.Project
}

var ProviderSet = wire.NewSet(
	store.ProviderSet,
	workflow.ProviderSet,
	openaiclient.ChatProviderSet,
	NewService,
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
	manager workflow.Manager,
	cache bucket.Bucket,
	chat *openaiclient.ChatServer,
	defaultProject *gen.Project,
) (*Service, error) {
	return &Service{
		store:          store,
		manager:        manager,
		cache:          cache,
		chat:           chat,
		defaultProject: defaultProject,
	}, nil
}

func nodesFromFiles(u string) ([]*gen.Node, error) {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse url %s", u)
	}
	// TODO breadchris support recursive enumeration
	files, err := os.ReadDir(parsedUrl.Path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read dir %s", parsedUrl.Path)
	}

	var nodes []*gen.Node
	for _, file := range files {
		// TODO breadchris need to collapse this instantiation
		nodes = append(nodes, &gen.Node{
			Name: file.Name(),
			Type: &gen.Node_Storage{
				Storage: &storage.Storage{
					Type: &storage.Storage_File{
						File: &storage.File{
							Path: file.Name(),
						},
					},
				},
			},
		})
	}
	return nodes, nil
}

// TODO breadchris this will be something that needs to be specified when someone is calling the API
func enumerateProvidersFromNodes(nodes []*gen.Node) ([]*gen.EnumeratedProvider, error) {
	var providers []*gen.EnumeratedProvider
	for _, node := range nodes {
		info := &gen.ProviderInfo{
			State: gen.ProviderState_READY,
			Error: "",
		}

		var (
			providedNodes []*gen.Node
			err           error
		)
		switch t := node.Type.(type) {
		case *gen.Node_Storage:
			switch u := t.Storage.Type.(type) {
			case *storage.Storage_Folder:
				providedNodes, err = nodesFromFiles(u.Folder.Url)
			}
		case *gen.Node_Grpc:
			switch u := t.Grpc.Type.(type) {
			case *pgrpc.GRPC_Server:
				providedNodes, err = grpc.EnumerateResourceBlocks(u.Server, false)
			}
		case *gen.Node_Code:
			switch u := t.Code.Type.(type) {
			case *code.Code_Server:
				providedNodes, err = grpc.EnumerateResourceBlocks(u.Server.Grpc, false)
			}
		case *gen.Node_Data:
			switch t.Data.Type.(type) {
			case *pdata.Data_Input:
				providedNodes = []*gen.Node{data.NewProto("input", data.NewInputProto())}
			}
		case *gen.Node_Reason:
			switch t.Reason.Type.(type) {
			case *preason.Reason_Engine:
				providedNodes = []*gen.Node{reason.NewProto("prompt", reason.NewPromptProto())}
			}
		default:
			continue
		}
		if len(providedNodes) == 0 {
			log.Warn().Msgf("no nodes provided by %s", node.Name)
			continue
		}
		if err != nil {
			info.State = gen.ProviderState_ERROR
			info.Error = err.Error()
		}
		providers = append(providers, &gen.EnumeratedProvider{
			Provider: node,
			Nodes:    providedNodes,
			Info:     info,
		})
	}
	return providers, nil
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

func (s *Service) GetNodeInfo(ctx context.Context, c *connect.Request[gen.GetNodeInfoRequest]) (*connect.Response[gen.GetNodeInfoResponse], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	w, err := workflow.Default().
		WithProtoProject(graph.ConvertProto(project)).
		Build()
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
