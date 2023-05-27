package resource

import (
	"context"
	"fmt"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/rs/zerolog/log"
	"gocloud.dev/blob"
)

type BlobstoreResource struct {
	*BaseResource
	*gen.Blobstore
}

var _ Resource = &BlobstoreResource{}

func (r *BlobstoreResource) Name() string {
	return BlobstoreResourceType
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
