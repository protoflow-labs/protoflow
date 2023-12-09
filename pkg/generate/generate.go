package generate

import (
	"github.com/pkg/errors"
	pcode "github.com/protoflow-labs/protoflow/gen/code"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/code"
	"github.com/protoflow-labs/protoflow/pkg/project"
	"github.com/protoflow-labs/protoflow/templates"
	"github.com/rs/zerolog/log"
	"io"
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

func copyFile() {

}

func (s *Generate) Init(project *project.Project) error {
	// TODO breadchris what does the wire frame look like?

	t := templates.Templates
	// TODO breadchris someone has done this already https://github.com/leaanthony/debme
	pPath := "proto"
	entries, err := t.ReadDir(pPath)
	if err != nil {
		return err
	}
	dPath, err := s.bucket.GetFolder(pPath)
	if err != nil {
		return err
	}
	for _, e := range entries {
		f, err := t.Open(path.Join(pPath, e.Name()))
		if err != nil {
			return err
		}
		df, err := os.Create(path.Join(dPath, e.Name()))
		_, err = io.Copy(df, f)
		if err != nil {
			return err
		}
	}
	bufFile := "buf.gen.yaml"
	tf, err := t.Open(bufFile)
	if err != nil {
		return err
	}
	b, err := s.bucket.GetFile(bufFile)
	if err != nil {
		return err
	}
	df, err := os.Create(b)
	if err != nil {
		return err
	}
	_, err = io.Copy(df, tf)
	if err != nil {
		return err
	}
	return nil
}

func (s *Generate) Service(project *project.Project, n graph.Node) error {
	return nil
}

func (s *Generate) Method(project *project.Project, n graph.Node) error {
	return nil
}

func (s *Generate) GenerateImplementation(project *project.Project, n graph.Node) error {
	r, err := n.Provider()
	if err != nil {
		return errors.Wrapf(err, "error getting node provider")
	}

	switch c := r.(type) {
	case *code.Server:
		if c.Runtime == pcode.Runtime_NODEJS {
			jsManager, err := NewNodeJSManager(s.bucket)
			if err != nil {
				return errors.Wrap(err, "error creating nodejs manager")
			}

			err = jsManager.GenerateFunctionImpl(c, n)
			if err != nil {
				log.Error().Err(err).Msg("error generating function implementation")
			}
			err = jsManager.GenerateGRPCService(c)
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

	r, err := n.Provider()
	if err != nil {
		return errors.Wrapf(err, "error getting node provider")
	}
	switch r := r.(type) {
	case *code.Server:
		switch r.Runtime {
		case pcode.Runtime_NODEJS:
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
	for _, r := range project.Workflow.NodeLookup {
		if r == nil {
			log.Error().Msg("resource is nil")
			continue
		}
		switch r := r.(type) {
		case *code.Server:
			if r.Runtime == pcode.Runtime_NODEJS {
				jsManager, err := NewNodeJSManager(s.bucket)
				if err != nil {
					return errors.Wrap(err, "error creating nodejs manager")
				}

				for _, n := range r.Dependents() {
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
