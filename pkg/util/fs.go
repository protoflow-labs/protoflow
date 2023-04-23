package util

import (
	"os"
)

func CreateTempDir() (string, func(), error) {
	dir, err := os.MkdirTemp("", "example")
	if err != nil {
		return "", nil, err
	}

	cleanup := func() {
		os.RemoveAll(dir)
	}

	return dir, cleanup, nil
}
