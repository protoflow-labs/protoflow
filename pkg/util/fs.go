package util

import (
	"os"
	"path"
)

func ProtoflowHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homeDir, ".protoflow"), nil
}
