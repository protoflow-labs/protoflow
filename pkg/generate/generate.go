package generate

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/project"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
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
	nodeInfoLookup := map[string]*node.Info{}
	for _, n := range project.Workflow.NodeLookup {
		info, err := project.Workflow.GetNodeInfo(n)
		if err != nil {
			return errors.Wrapf(err, "error getting node info")
		}
		nodeInfoLookup[n.ID()] = info
	}
	for _, r := range project.Workflow.Resources {
		if r == nil {
			log.Error().Msg("resource is nil")
			continue
		}
		var nodes []node.Node
		for _, n := range project.Workflow.NodeLookup {
			if n.ResourceID() == r.ID() {
				nodes = append(nodes, n)
			}
		}
		switch r := r.(type) {
		case *resource.LanguageServiceResource:
			if r.Runtime == gen.Runtime_NODEJS {
				jsManager, err := NewNodeJSManager(s.bucket)
				if err != nil {
					return errors.Wrap(err, "error creating nodejs manager")
				}
				if err := jsManager.Generate(r, nodes, nodeInfoLookup); err != nil {
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
