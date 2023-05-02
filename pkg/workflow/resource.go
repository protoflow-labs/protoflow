package workflow

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/cache"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"github.com/rs/zerolog/log"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
	"gocloud.dev/docstore"
	"gocloud.dev/docstore/memdocstore"
	_ "gocloud.dev/docstore/memdocstore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	Conn *grpc.ClientConn
}

func (r *GRPCResource) Name() string {
	return GRPCResourceType
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
	Conn  *grpc.ClientConn
	Cache cache.Cache
	cmd   *exec.Cmd
}

func (r *LanguageServiceResource) Name() string {
	return LanguageServiceType
}

func (r *LanguageServiceResource) Init() (func(), error) {
	if err := r.ensureRunning(); err != nil {
		return nil, errors.Wrapf(err, "Unable to get the %s language server running", r.Name())
	}
	conn, err := grpc.Dial(r.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to grpc server at %s", r.Host)
	}
	cleanup := func() {
		// TODO breadchris configure this from the frontend? configure?
		//if !r.CloseOnCleanup {
		//	return
		//}
		if r.cmd != nil && r.cmd.Process != nil {
			syscall.Kill(-r.cmd.Process.Pid, syscall.SIGKILL)
		}
		err = conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing grpc connection")
		}
	}
	r.Conn = conn
	return cleanup, nil
}

func (r *LanguageServiceResource) ensureRunning() error {
	originalWd, err := os.Getwd()
	if err != nil {
		return err
	}

	// TODO - Figure out how to pass project name in cleanly
	projectFolder, _ := r.Cache.GetFolder("projects/local")
	if err := os.Chdir(projectFolder); err != nil {
		return err
	}

	r.cmd = exec.Command("npm", "run", "dev")
	r.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stdout, err := r.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err = r.cmd.Start(); err != nil {
		return err
	}

	if err := os.Chdir(originalWd); err != nil {
		return err
	}

	err = r.waitForPort()
	if err != nil {
		out, err := io.ReadAll(stdout)
		if err != nil {
			return errors.Wrapf(err, "Unable to get the %s language server running, failed to read stdout", r.Name())
		}
		return errors.Wrapf(err, "Unable to get the %s language server running: %s", r.Name(), string(out))
	}
	return nil
}

func (r *LanguageServiceResource) waitForPort() error {
	maxRetries := 20
	retryInterval := 2 * time.Second

	fmt.Printf("Waiting for %s to start listening...\n", r.LanguageService.Host)
	for i := 1; i <= maxRetries; i++ {
		conn, err := net.DialTimeout("tcp", r.LanguageService.Host, time.Second)
		if err == nil {
			conn.Close()
			fmt.Printf("%s is now listening!\n", r.LanguageService.Host)
			return nil
		} else {
			fmt.Printf("%s is not yet listening (attempt %d/%d)\n", r.LanguageService.Host, i, maxRetries)
			time.Sleep(retryInterval)
		}
	}
	return fmt.Errorf("port did not come online in time")
}
