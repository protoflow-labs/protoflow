package bucket

import (
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"path"
)

type Builder struct {
	*blob.Bucket
	config Config

	path string
}

func (s *Builder) Dir(name string) *Builder {
	ns := *s
	ns.path = path.Join(s.config.Path, name)
	return &ns
}

func (s *Builder) Build() (string, error) {
	return s.path, ensureDirExists(s.path)
}

func (s *Builder) File(name string) (string, error) {
	return path.Join(s.path, name), ensureDirExists(s.path)
}

func NewBuilder(config Config) (*Builder, error) {
	var (
		err error
	)
	if config.Path == "" {
		config.Path, err = createLocalDir(config.LocalName)
		if err != nil {
			return nil, err
		}
		err = ensureDirExists(path.Join(config.Path, "bucket"))
		if err != nil {
			return nil, err
		}
	}
	bucket, err := fileblob.OpenBucket(config.Path, &fileblob.Options{
		CreateDir: true,
	})
	if err != nil {
		return nil, err
	}
	return &Builder{
		Bucket: bucket,
		config: config,
		path:   config.Path,
	}, nil
}

func NewTestBuilder() *Builder {
	return &Builder{
		Bucket: nil,
		config: Config{
			LocalName: "test",
			Path:      "",
		},
		path: "",
	}
}
