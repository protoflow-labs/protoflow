package cache

import (
	"github.com/google/wire"
	"github.com/pkg/errors"
	"os"
	"os/user"
	"path"
)

type Cache interface {
	GetFile(name string) (string, error)
	GetFolder(name string) (string, error)
	WithDir(dir string) (Cache, error)
}

type LocalCache struct {
	dir    string
	config Config
}

var _ Cache = &LocalCache{}

var ProviderSet = wire.NewSet(
	NewConfig,
	NewUserCache,
	wire.Bind(new(Cache), new(*LocalCache)),
)

func (c *LocalCache) GetFile(name string) (string, error) {
	return path.Join(c.dir, name), nil
}

func (c *LocalCache) GetFolder(name string) (string, error) {
	p := path.Join(c.dir, name)
	return p, ensureDirExists(p)
}

func (c *LocalCache) WithDir(dir string) (Cache, error) {
	p := path.Join(c.dir, dir)
	err := ensureDirExists(p)
	if err != nil {
		return nil, err
	}
	return FromDir(p)
}

func NewUserCache(c Config) (*LocalCache, error) {
	folder, err := createLocalDir(c.Name)
	if err != nil {
		return nil, err
	}

	return &LocalCache{
		dir:    folder,
		config: c,
	}, nil
}

func FromDir(dir string) (*LocalCache, error) {
	// check to see if dir exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "could not find dir: %v", dir)
	}
	return &LocalCache{
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
