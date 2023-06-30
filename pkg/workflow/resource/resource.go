package resource

import (
	"fmt"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
	_ "gocloud.dev/docstore/memdocstore"
)

const (
	GRPCResourceType      = "grpc"
	DocstoreResourceType  = "docstore"
	BlobstoreResourceType = "blobstore"
	LanguageServiceType   = "language"
)

type DeploymentInfo struct {
	ContainerURI string
	Ports        []int
	Volumes      []string
}

type Resource interface {
	Init() (func(), error)
	Name() string
	ID() string
	AddNode(n node.Node)
	Nodes() []node.Node
	Info(n node.Node) (*node.Info, error)
	//DeploymentInfo() (*DeploymentInfo, error)
}

type BaseResource struct {
	id    string
	nodes []node.Node
}

func (r *BaseResource) ID() string {
	return r.id
}

func (r *BaseResource) Name() string {
	return fmt.Sprintf("%s", r.ID())
}

func (r *BaseResource) Init() (func(), error) {
	return func() {}, nil
}

func (r *BaseResource) Info(n node.Node) (*node.Info, error) {
	return nil, nil
}

func (r *BaseResource) AddNode(n node.Node) {
	r.nodes = append(r.nodes, n)
}

func (r *BaseResource) Nodes() []node.Node {
	return r.nodes
}

func FromProto(r *gen.Resource) (Resource, error) {
	base := &BaseResource{
		id: r.Id,
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
	default:
		return nil, fmt.Errorf("no resource found with type: %s", t)
	}
}
