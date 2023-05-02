package generate

import (
	"fmt"
	"path"

	"github.com/protoflow-labs/protoflow/pkg/cache"
	"github.com/protoflow-labs/protoflow/pkg/util"

	"github.com/pkg/errors"

	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/templates"
)

type Generator interface {
	Generate(project *gen.Project) error
}

type Generate struct {
	cache cache.Cache
}

var _ Generator = &Generate{}

func NewGenerate(cache cache.Cache) *Generate {
	return &Generate{
		cache: cache,
	}
}

func (g *Generate) Generate(project *gen.Project) error {
	var err error

	// normalize all node names
	for _, node := range project.GetGraph().GetNodes() {
		node.Name = util.ToTitleCase(node.Name)
	}

	projectDir := path.Join("projects", project.GetName())

	functionNodes, err := g.ScaffoldFunctions(project, projectDir)
	if err != nil {
		return errors.Wrapf(err, "error scaffolding functions for %s", project.GetName())
	}

	err = g.GenerateServiceProtos(projectDir, functionNodes)
	if err != nil {
		return errors.Wrapf(err, "error generating service protos for %s", project.GetName())
	}

	err = g.GenerateServices(projectDir, functionNodes)
	if err != nil {
		return errors.Wrapf(err, "error generating services for %s", project.GetName())
	}

	return nil
}

func (g *Generate) ScaffoldFunctions(project *gen.Project, projectDir string) ([]*gen.Node, error) {
	nodes := project.GetGraph().GetNodes()
	var funcNodes []*gen.Node
	for _, node := range nodes {
		if node.GetFunction() == nil {
			continue
		}

		funcNodes = append(funcNodes, node)

		// create function directory
		funcDir := path.Join(projectDir, "functions", node.GetName())
		funcDirPath, err := g.cache.GetFolder(funcDir)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating function directory %s", funcDir)
		}

		if node.GetFunction().Runtime == "node" {
			err := templates.TemplateFile("node/function.index.template.js", funcDirPath+"/index.js", map[string]interface{}{
				"Node": node,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return funcNodes, nil
}

func (g *Generate) GenerateServiceProtos(projectDir string, functionNodes []*gen.Node) error {
	functions := make(map[string][]string, 0)

	protosPath, err := g.cache.GetFolder(path.Join(projectDir, "protos"))
	if err != nil {
		return errors.Wrapf(err, "error getting protos folder %s", path.Join(protosPath, "protos"))
	}

	for _, node := range functionNodes {
		functions[node.GetFunction().Runtime] = append(functions[node.GetFunction().Runtime], node.GetName())
	}

	for runtime, methods := range functions {
		protoPath := fmt.Sprintf("%s/%s.service.proto", protosPath, runtime)
		err := templates.TemplateFile("service.template.proto", protoPath, map[string]interface{}{
			"Runtime": runtime,
			"Methods": methods,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generate) GenerateServices(projectDir string, functionNodes []*gen.Node) error {
	projectPath, err := g.cache.GetFolder(projectDir)
	if err != nil {
		return errors.Wrapf(err, "error getting project folder %s", projectDir)
	}
	return templates.TemplateDir("node/project", projectPath, map[string]interface{}{
		"FunctionNodes": functionNodes,
	})
}
