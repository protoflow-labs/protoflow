package bucket

import (
	"os"
	"os/user"
	"path"

	"github.com/google/wire"
	"github.com/pkg/errors"
)

type Bucket interface {
	GetFile(name string) (string, error)
	GetFolder(name string) (string, error)
	WithDir(dir string) (Bucket, error)
}

type LocalBucket struct {
	dir    string
	config Config
}

var _ Bucket = &LocalBucket{}

var ProviderSet = wire.NewSet(
	NewConfig,
	NewUserCache,
	wire.Bind(new(Bucket), new(*LocalBucket)),
)

func (c *LocalBucket) GetFile(name string) (string, error) {
	return path.Join(c.dir, name), nil
}

func (c *LocalBucket) GetFolder(name string) (string, error) {
	p := path.Join(c.dir, name)
	return p, ensureDirExists(p)
}

func (c *LocalBucket) WithDir(dir string) (Bucket, error) {
	p := path.Join(c.dir, dir)
	err := ensureDirExists(p)
	if err != nil {
		return nil, err
	}
	return FromDir(p)
}

func NewUserCache(c Config) (*LocalBucket, error) {
	folder, err := createLocalDir(c.Name)
	if err != nil {
		return nil, err
	}

	return &LocalBucket{
		dir:    folder,
		config: c,
	}, nil
}

func FromDir(dir string) (*LocalBucket, error) {
	err := ensureDirExists(dir)
	if err != nil {
		return nil, err
	}
	return &LocalBucket{
		dir: dir,
	}, nil
}

func ensureDirExists(p string) error {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		if err := os.MkdirAll(p, 0700); err != nil {
			return errors.Wrapf(err, "could not create all dirs for path: %v", p)
		}
	}
	return nil
}

func createLocalDir(dirName string) (string, error) {
	// Get the current user
	u, err := user.Current()
	if err != nil {
		return "", errors.Wrapf(err, "could not get current user")
	}

	p := path.Join(u.HomeDir, dirName)
	return p, ensureDirExists(p)
}
