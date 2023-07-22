package resource

import (
	"fmt"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
	_ "gocloud.dev/docstore/memdocstore"
)

type DeploymentInfo struct {
	ContainerURI string
	Ports        []int
	Volumes      []string
}

type BaseResource struct {
	id               string
	name             string
	dependencyLookup map[string]graph.Resource
	nodes            []graph.Node
}

func (r *BaseResource) ID() string {
	return r.id
}

func (r *BaseResource) Name() string {
	return r.name
}

func (r *BaseResource) Init() (func(), error) {
	return func() {}, nil
}

func (r *BaseResource) Info(n graph.Node) (*graph.Info, error) {
	return nil, nil
}

func (r *BaseResource) AddNode(n graph.Node) {
	r.nodes = append(r.nodes, n)
}

func (r *BaseResource) Nodes() []graph.Node {
	return r.nodes
}

func (r *BaseResource) ResolveDependencies(dp graph.DependencyProvider) error {
	for id := range r.dependencyLookup {
		if _, ok := dp[id]; !ok {
			return fmt.Errorf("dependency not found: %s", id)
		}
		r.dependencyLookup[id] = dp[id]
	}
	return nil
}

func FromProto(r *gen.Resource) (graph.Resource, error) {
	base := &BaseResource{
		id:   r.Id,
		name: r.Name,
	}
	// TODO breadchris is this too sketch?
	base.dependencyLookup = make(map[string]graph.Resource)
	for _, dep := range r.Dependencies {
		base.dependencyLookup[dep] = nil
	}
	switch t := r.Type.(type) {
	case *gen.Resource_LanguageService:
		return &LanguageServiceResource{
			GRPCResource: &GRPCResource{
				BaseResource: base,
				GRPCService:  t.LanguageService.Grpc,
			},
			LanguageService: t.LanguageService,
		}, nil
	case *gen.Resource_GrpcService:
		return &GRPCResource{
			BaseResource: base,
			GRPCService:  t.GrpcService,
		}, nil
	case *gen.Resource_DocStore:
		return &DocstoreResource{
			BaseResource: base,
			DocStore:     t.DocStore,
		}, nil
	case *gen.Resource_FileStore:
		return &FileStoreResource{
			BaseResource: base,
			FileStore:    t.FileStore,
		}, nil
	case *gen.Resource_ReasoningEngine:
		return &ReasoningEngineResource{
			BaseResource:    base,
			ReasoningEngine: t.ReasoningEngine,
		}, nil
	case *gen.Resource_ConfigProvider:
		return &ConfigProviderResource{
			BaseResource:   base,
			ConfigProvider: t.ConfigProvider,
		}, nil
	case *gen.Resource_SecretStore:
		return &SecretStoreResource{
			BaseResource: base,
			SecretStore:  t.SecretStore,
		}, nil
	case *gen.Resource_TemplateService:
		return &TemplateServiceResource{
			BaseResource:    base,
			TemplateService: t.TemplateService,
		}, nil
	case *gen.Resource_HttpRouter:
		return &HTTPRouterResource{
			BaseResource: base,
			HTTPRouter:   t.HttpRouter,
		}, nil
	default:
		return nil, fmt.Errorf("no resource found with type: %s", t)
	}
}
