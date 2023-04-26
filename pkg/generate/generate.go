package generate

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/pkg/errors"

	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/templates"
)

type Generator struct {
	project              *gen.Project
	projectDir           string
	functionNodes        []*gen.Node
	serviceProtoTemplate *template.Template
}

func NewFromProject(project *gen.Project) (*Generator, error) {
	if project == nil {
		return nil, errors.New("project must be provided")
	}

	serviceProtoTemplate, err := template.
		ParseFS(templates.Templates, "service.template.proto")
	if err != nil {
		return nil, err
	}

	return &Generator{
		project:              project,
		serviceProtoTemplate: serviceProtoTemplate,
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
		os.MkdirAll(g.projectDir+"/functions/"+node.GetName(), 0700)
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
		protoFile, err := os.Create(fmt.Sprintf("%s/%s.service.proto", protosDir, runtime))
		if err != nil {
			return err
		}

		err = g.serviceProtoTemplate.Execute(protoFile, map[string]interface{}{
			"Runtime": runtime,
			"Methods": methods,
		})
		if err != nil {
			return err
		}

		protoFile.Close()
	}

	return nil
}

func (g *Generator) GenerateServices() error {
	cmd := "npx"
	args := []string{"buf", "generate", "proto"}

	_, err := exec.Command(cmd, args...).Output()

	return err
}
