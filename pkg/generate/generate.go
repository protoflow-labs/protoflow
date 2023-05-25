package generate

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/project"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
	"os"
	"path"
)

type Generator interface {
	Generate(project *project.Project) error
}

type Generate struct {
	bucket bucket.Bucket
}

var _ Generator = &Generate{}

func NewGenerate(config Config) (*Generate, error) {
	var projectDir string
	// if there is a project path defined, use this for where the bucket goes
	if config.ProjectPath != "" {
		projectDir = config.ProjectPath
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, errors.Wrap(err, "error getting current working directory")
		}
		projectDir = path.Join(cwd, "protoflow")
	}
	c, err := bucket.FromDir(projectDir)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating bucket from %s", projectDir)
	}
	return &Generate{
		bucket: c,
	}, nil
}

func (s *Generate) Generate(project *project.Project) error {
	for _, r := range project.Workflow.Resources {
		switch r := r.(type) {
		case *workflow.LanguageServiceResource:
			if r.Runtime == gen.Runtime_NODE {
				jsManager, err := NewNodeJSManager(s.bucket)
				if err != nil {
					return errors.Wrap(err, "error creating nodejs manager")
				}
				if err := jsManager.Generate(r); err != nil {
					return errors.Wrap(err, "error generating nodejs")
				}
			}
		}
	}

	d := NewDockerComposeGenerate(s.bucket)
	if err := d.Generate(project); err != nil {
		return errors.Wrap(err, "error generating docker-compose")
	}
	return nil
}
