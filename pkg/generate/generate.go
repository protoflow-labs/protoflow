package generate

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/project"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/rs/zerolog/log"
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

// TODO breadchris should this be derived from somewhere?
const protoflowDir = "protoflow"

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
		projectDir = path.Join(cwd, protoflowDir)
	}
	c, err := bucket.FromDir(projectDir)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating bucket from %s", projectDir)
	}
	return &Generate{
		bucket: c,
	}, nil
}

func (s *Generate) GenerateImplementation(project *project.Project, n graph.Node) error {
	r, ok := project.Workflow.Resources[n.ResourceID()]
	if !ok || r == nil {
		return errors.Errorf("resource %s not found", n.ResourceID())
	}
	switch r := r.(type) {
	case *resource.LanguageServiceResource:
		if r.Runtime == gen.Runtime_NODEJS {
			jsManager, err := NewNodeJSManager(s.bucket)
			if err != nil {
				return errors.Wrap(err, "error creating nodejs manager")
			}

			err = jsManager.GenerateFunctionImpl(r, n)
			if err != nil {
				log.Error().Err(err).Msg("error generating function implementation")
			}
			err = jsManager.GenerateGRPCService(r)
			if err != nil {
				log.Error().Err(err).Msg("error generating service files")
			}
		}
	}
	return nil
}

func (s *Generate) InferNodeType(project *project.Project, n graph.Node) error {
	info, err := project.Workflow.GetNodeInfo(n)
	if err != nil {
		return errors.Wrapf(err, "error getting node info")
	}

	r, ok := project.Workflow.Resources[n.ResourceID()]
	if !ok || r == nil {
		return errors.Errorf("resource %s not found", n.ResourceID())
	}
	switch r := r.(type) {
	case *resource.LanguageServiceResource:
		switch r.Runtime {
		case gen.Runtime_NODEJS:
			jsManager, err := NewNodeJSManager(s.bucket)
			if err != nil {
				return errors.Wrap(err, "error creating nodejs manager")
			}

			err = jsManager.UpdateNodeType(n, info)
			if err != nil {
				log.Error().Err(err).Msg("error updating node type")
			}
		}
	}
	return nil
}

func (s *Generate) Generate(project *project.Project) error {
	for _, r := range project.Workflow.Resources {
		if r == nil {
			log.Error().Msg("resource is nil")
			continue
		}
		switch r := r.(type) {
		case *resource.LanguageServiceResource:
			if r.Runtime == gen.Runtime_NODEJS {
				jsManager, err := NewNodeJSManager(s.bucket)
				if err != nil {
					return errors.Wrap(err, "error creating nodejs manager")
				}

				for _, n := range r.Nodes() {
					info, err := project.Workflow.GetNodeInfo(n)
					if err != nil {
						return errors.Wrapf(err, "error getting node info")
					}
					err = jsManager.UpdateNodeType(n, info)
					if err != nil {
						log.Error().Err(err).Msg("error updating node type")
					}
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
