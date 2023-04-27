package generate

import (
	"fmt"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"os"

	"github.com/pkg/errors"

	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/templates"
)

type Generator struct {
	project       *gen.Project
	projectDir    string
	functionNodes []*gen.Node
}

func NewFromProject(project *gen.Project) (*Generator, error) {
	if project == nil {
		return nil, errors.New("project must be provided")
	}

	for _, node := range project.GetGraph().GetNodes() {
		node.Name = util.ToTitleCase(node.Name)
	}

	return &Generator{
		project: project,
	}, nil
}

func (g *Generator) Generate() error {
	var err error

	err = g.MakeProjectDir()
	if err != nil {
		return errors.Wrapf(err, "error creating project directory for %s", g.project.GetName())
	}

	err = g.ScaffoldFunctions()
	if err != nil {
		return errors.Wrapf(err, "error scaffolding functions for %s", g.project.GetName())
	}

	err = g.GenerateServiceProtos()
	if err != nil {
		return errors.Wrapf(err, "error generating service protos for %s", g.project.GetName())
	}

	err = g.GenerateServices()
	if err != nil {
		return errors.Wrapf(err, "error generating services for %s", g.project.GetName())
	}

	return nil
}

func (g *Generator) MakeProjectDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	g.projectDir = homeDir + "/.protoflow/projects/" + g.project.GetName()

	if err := os.MkdirAll(g.projectDir, 0700); err != nil {
		return err
	}

	return nil
}

func (g *Generator) ScaffoldFunctions() error {
	nodes := g.project.GetGraph().GetNodes()
	for _, node := range nodes {
		if node.GetFunction() == nil {
			continue
		}

		g.functionNodes = append(g.functionNodes, node)

		// create function directory
		funcDir := g.projectDir + "/functions/" + node.GetName()
		os.MkdirAll(funcDir, 0700)

		if node.GetFunction().Runtime == "node" {
			err := templates.TemplateFile("node/function.index.template.js", funcDir+"/index.js", map[string]interface{}{
				"Node": node,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Generator) GenerateServiceProtos() error {
	functions := make(map[string][]string, 0)
	protosDir := g.projectDir + "/protos"

	os.MkdirAll(protosDir, 0700)

	for _, node := range g.functionNodes {
		functions[node.GetFunction().Runtime] = append(functions[node.GetFunction().Runtime], node.GetName())
	}

	for runtime, methods := range functions {
		protoFilename := fmt.Sprintf("%s/%s.service.proto", protosDir, runtime)
		err := templates.TemplateFile("service.template.proto", protoFilename, map[string]interface{}{
			"Runtime": runtime,
			"Methods": methods,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) GenerateServices() error {
	err := templates.TemplateDir("node/project", g.projectDir, map[string]interface{}{
		"FunctionNodes": g.functionNodes,
	})
	return err
}
