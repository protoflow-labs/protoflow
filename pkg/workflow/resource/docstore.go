package resource

import (
	"context"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"github.com/rs/zerolog/log"
	"gocloud.dev/docstore"
	"gocloud.dev/docstore/memdocstore"
	"os"
	"path"
	"strings"
)

type DocstoreResource struct {
	*BaseResource
	*gen.DocStore
}

var _ Resource = &DocstoreResource{}

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
		// TODO breadchris replace this with bucket.Cache.GetFolder
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
