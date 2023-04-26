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
		g := r.GetGrpcService()
		return &GRPCResource{
			GRPCService: g,
		}, nil
	case *gen.Resource_Docstore:
		d := r.GetDocstore()
		return &DocstoreResource{
			Docstore: d,
		}, nil
	case *gen.Resource_Blobstore:
		b := r.GetBlobstore()
		return &BlobstoreResource{
			Blobstore: b,
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
	Collection *docstore.Collection
}

func (r *DocstoreResource) Init() (func(), error) {
	coll, err := docstore.OpenCollection(context.Background(), r.Url)
	if err != nil {
		return nil, fmt.Errorf("could not open collection: %v", err)
	}
	r.Collection = coll
	return func() {
		err = coll.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing docstore collection")
		}
	}, nil
}

type BlobstoreResource struct {
	*gen.Blobstore
	Bucket *blob.Bucket
}

func (r *BlobstoreResource) Init() (func(), error) {
	bucket, err := blob.OpenBucket(context.Background(), r.Url)
	if err != nil {
		return nil, fmt.Errorf("could not open bucket: %v", err)
	}
	r.Bucket = bucket
	return func() {
		err = bucket.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing blobstore bucket")
		}
	}, nil
}
