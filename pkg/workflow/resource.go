package workflow

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"github.com/rs/zerolog/log"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
	"gocloud.dev/docstore"
	"gocloud.dev/docstore/memdocstore"
	_ "gocloud.dev/docstore/memdocstore"
)

const (
	GRPCResourceType      = "grpc"
	DocstoreResourceType  = "docstore"
	BlobstoreResourceType = "blobstore"
	LanguageServiceType   = "language"
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
	Name() string
}

type GRPCResource struct {
	*gen.GRPCService
}

func (r *GRPCResource) Name() string {
	return GRPCResourceType
}

func (r *GRPCResource) Init() (func(), error) {
	// TODO breadchris this is a hack to get the grpc server running, this is not ideal
	if !strings.HasPrefix(r.Host, "http://") {
		r.Host = "http://" + r.Host
	}
	if err := ensureRunning(r.Host); err != nil {
		// TODO breadchris ignore errors for now
		// return nil, errors.Wrapf(err, "unable to get the %s grpc server running", r.Name())
		return nil, nil
	}
	return nil, nil
}

type DocstoreResource struct {
	*gen.Docstore
}

func (r *DocstoreResource) Name() string {
	return DocstoreResourceType
}

func (r *DocstoreResource) Init() (func(), error) {
	return nil, nil
}

func (r *DocstoreResource) WithCollection(name string) (*docstore.Collection, func(), error) {
	var (
		coll     *docstore.Collection
		err      error
		protoDir string
	)
	if strings.HasPrefix(r.Url, "mem://") {
		// TODO breadchris replace this with cache.Cache.GetFolder
		protoDir, err = util.ProtoflowHomeDir()
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not get protoflow home dir")
		}

		filename := path.Join(protoDir, name+".json")

		// TODO breadchris "id" is
		coll, err = memdocstore.OpenCollection("id", &memdocstore.Options{
			Filename: filename,
		})
		if err != nil {
			// remove file if it exists
			if os.IsNotExist(err) {
				return nil, nil, errors.Wrapf(err, "could not open memory docstore collection: %s", name)
			}
			err = os.Remove(filename)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "could not remove memory docstore collection: %s", name)
			}
		}
	} else {
		coll, err = docstore.OpenCollection(context.Background(), r.Url+"/"+name)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "could not open docstore collection: %s", name)
		}
	}

	return coll, func() {
		if coll == nil {
			log.Debug().Msg("docstore collection is nil")
			return
		}
		err = coll.Close()
		if err != nil {
			log.Error().Msgf("error closing docstore collection: %+v", err)
		}
	}, nil
}

type BlobstoreResource struct {
	*gen.Blobstore
}

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

type LanguageServiceResource struct {
	*gen.LanguageService
	GRPC *GRPCResource
}

func (r *LanguageServiceResource) Name() string {
	return LanguageServiceType
}

func (r *LanguageServiceResource) Init() (func(), error) {
	r.GRPC = &GRPCResource{
		GRPCService: r.LanguageService.Grpc,
	}
	return r.GRPC.Init()
}

func ensureRunning(host string) error {
	maxRetries := 1
	retryInterval := 2 * time.Second

	u, err := url.Parse(host)
	if err != nil {
		return errors.Wrapf(err, "unable to parse url %s", host)
	}

	log.Debug().Str("host", host).Msg("waiting for host to come online")
	for i := 1; i <= maxRetries; i++ {
		conn, err := net.DialTimeout("tcp", u.Host, time.Second)
		if err == nil {
			conn.Close()
			log.Debug().Str("host", host).Msg("host is not listening")
			return nil
		} else {
			log.Debug().Err(err).Int("attempt", i).Int("max", maxRetries).Msg("error connecting to host")
			time.Sleep(retryInterval)
		}
	}
	return errors.New("host did not come online in time")
}
