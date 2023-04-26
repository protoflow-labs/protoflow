package workflow

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/rs/zerolog/log"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
	"gocloud.dev/docstore"
	_ "gocloud.dev/docstore/memdocstore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ResourceFromProto(r *gen.Resource) (Resource, error) {
	switch t := r.Type.(type) {
	case *gen.Resource_GrpcService:
		return &GRPCResource{
			GRPCService: r.GetGrpcService(),
		}, nil
	case *gen.Resource_Docstore:
		return &DocstoreResource{
			Docstore: r.GetDocstore(),
		}, nil
	case *gen.Resource_Blobstore:
		return &BlobstoreResource{
			Blobstore: r.GetBlobstore(),
		}, nil
	default:
		return nil, fmt.Errorf("no resource found with type: %s", t)
	}
}

type Resource interface {
	Init() (func(), error)
}

type GRPCResource struct {
	*gen.GRPCService
	Conn *grpc.ClientConn
}

func (r *GRPCResource) Init() (func(), error) {
	conn, err := grpc.Dial(r.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to grpc server at %s", r.Host)
	}
	cleanup := func() {
		err = conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing grpc connection")
		}
	}
	r.Conn = conn
	return cleanup, nil
}

type DocstoreResource struct {
	*gen.Docstore
}

func (r *DocstoreResource) Init() (func(), error) {
	return nil, nil
}

func (r *DocstoreResource) WithKeyField(keyfield string) (*docstore.Collection, func(), error) {
	coll, err := docstore.OpenCollection(context.Background(), r.Url+"/"+keyfield)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open collection: %v", err)
	}
	return coll, func() {
		err = coll.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing docstore collection")
		}
	}, nil
}

type BlobstoreResource struct {
	*gen.Blobstore
}

func (r *BlobstoreResource) Init() (func(), error) {
	return nil, nil
}

func (r *BlobstoreResource) WithPath(path string) (*blob.Bucket, func(), error) {
	// remove leading slash
	if path[0] == '/' {
		path = path[1:]
	}
	// TODO breadchris validation of this url working should be done on init
	bucket, err := blob.OpenBucket(context.Background(), r.Url+"?prefix="+path)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open bucket: %v", err)
	}
	return bucket, func() {
		err = bucket.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing blobstore bucket")
		}
	}, nil
}
