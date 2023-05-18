package generate

import (
	"fmt"
	"os"
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
	codeBucket cache.Cache
}

var _ Generator = &Generate{}

func NewGenerate(c cache.Cache) *Generate {
	return &Generate{
		codeBucket: c,
	}
}

func (g *Generate) Generate(project *gen.Project) error {
	var err error

	// TODO breadchris support multiple languages
	code, err := g.codeBucket.WithDir("nodejs")
	if err != nil {
		return errors.Wrapf(err, "error creating codeBucket for %s", "nodejs")
	}

	// normalize all node names
	for _, node := range project.GetGraph().GetNodes() {
		node.Name = util.ToTitleCase(node.Name)
	}

	functionNodes, err := g.scaffoldFunctions(code, project)
	if err != nil {
		return errors.Wrapf(err, "error scaffolding functions for %s", project.GetName())
	}

	err = g.generateServiceProtos(code, functionNodes)
	if err != nil {
		return errors.Wrapf(err, "error generating service protos for %s", project.GetName())
	}

	err = g.generateServices(code, functionNodes)
	if err != nil {
		return errors.Wrapf(err, "error generating services for %s", project.GetName())
	}

	return nil
}

func (g *Generate) scaffoldFunctions(code cache.Cache, project *gen.Project) ([]*gen.Node, error) {
	nodes := project.GetGraph().GetNodes()
	var funcNodes []*gen.Node
	for _, node := range nodes {
		if node.GetFunction() == nil {
			continue
		}

		funcNodes = append(funcNodes, node)

		// create function directory
		funcDir := path.Join("functions", node.GetName())
		funcDirPath, err := code.GetFolder(funcDir)
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

func (g *Generate) generateServiceProtos(code cache.Cache, functionNodes []*gen.Node) error {
	functions := make(map[string][]string, 0)

	protosPath, err := code.GetFolder("protos")
	if err != nil {
		return errors.Wrapf(err, "error getting protos folder %s", path.Join(protosPath, "protos"))
	}

	for _, node := range functionNodes {
		functions[node.GetFunction().Runtime] = append(functions[node.GetFunction().Runtime], node.GetName())
	}

	for runtime, methods := range functions {
		protoPath := fmt.Sprintf("%s/%s.service.proto", protosPath, runtime)

		// check if file exists with os
		if _, err := os.Stat(protoPath); err == nil {
			// TODO breadchris we want to be more intelligent here. If the proto file exists, we should check if the methods are the same
			// and compile a new protofile with changes
			// proto file exists, skip
			continue
		}

		err = templates.TemplateFile("service.template.proto", protoPath, map[string]interface{}{
			"Runtime": runtime,
			"Methods": methods,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generate) generateServices(code cache.Cache, functionNodes []*gen.Node) error {
	projectPath, err := code.GetFolder("/")
	if err != nil {
		return errors.Wrapf(err, "error getting project folder")
	}
	return templates.TemplateDir("node/project", projectPath, map[string]interface{}{
		"FunctionNodes": functionNodes,
	})
}

func (g *Generate) generateDeployment(code cache.Cache, project *gen.Project) error {
	// TODO breadchris build a docker-compose file to bring up the project in one command
	// https://github.com/compose-spec/compose-go/blob/32b078a49aebb9a4a67af1423ec508f70b8928f7/types/types.go#L29
	return nil
}
