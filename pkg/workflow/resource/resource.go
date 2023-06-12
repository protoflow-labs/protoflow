package resource

import (
	"fmt"
	"github.com/protoflow-labs/protoflow/gen"
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

type Resource interface {
	Init() (func(), error)
	Name() string
	ID() string
}

type BaseResource struct {
	id string
}

func (r *BaseResource) ID() string {
	return r.id
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
	case *gen.Resource_Docstore:
		return &DocstoreResource{
			Docstore: t.Docstore,
		}, nil
	case *gen.Resource_Blobstore:
		return &BlobstoreResource{
			Blobstore: t.Blobstore,
		}, nil
	default:
		return nil, fmt.Errorf("no resource found with type: %s", t)
	}
}
