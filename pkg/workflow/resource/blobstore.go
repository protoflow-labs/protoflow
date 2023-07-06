package resource

import (
	"context"
	"fmt"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/rs/zerolog/log"
	"gocloud.dev/blob"
)

type FileStoreResource struct {
	*BaseResource
	*gen.FileStore
}

var _ Resource = &FileStoreResource{}

func (r *FileStoreResource) Init() (func(), error) {
	return nil, nil
}

func (r *FileStoreResource) WithPath(path string) (*blob.Bucket, func(), error) {
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
