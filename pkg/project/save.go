package project

import (
	"github.com/protoflow-labs/protoflow/gen"
	"google.golang.org/protobuf/proto"
	"os"
)

func SaveToFile(project *gen.Project, path string) error {
	b, err := proto.Marshal(project)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func LoadFromFile(path string) (*gen.Project, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	p := &gen.Project{}
	err = proto.Unmarshal(b, p)
	return p, err
}
