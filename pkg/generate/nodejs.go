package generate

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/grpc/manager"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/protoflow-labs/protoflow/templates"
	"github.com/samber/lo"
	"google.golang.org/protobuf/reflect/protoreflect"
	"path"
	"strings"
)

type LanguageManager interface {
	Generate(r *resource.LanguageServiceResource, nodes []node.Node, nodeInfoLookup map[string]*node.Info) error
}

type NodeJSManager struct {
	codeRoot bucket.Bucket
}

var _ LanguageManager = &NodeJSManager{}

func NewNodeJSManager(c bucket.Bucket) (*NodeJSManager, error) {
	codeRoot, err := c.WithDir("nodejs")
	if err != nil {
		return nil, errors.Wrapf(err, "error creating bucket for %s", "nodejs")
	}
	return &NodeJSManager{
		codeRoot: codeRoot,
	}, nil
}

type Method struct {
	Name       string
	InputType  string
	OutputType string
}

type FunctionTemplate struct {
	Node *node.FunctionNode
}

type ServiceTemplate struct {
	Runtime    string
	Methods    []Method
	DescLookup map[string]protoreflect.MessageDescriptor
	EnumLookup map[string]protoreflect.EnumDescriptor
	Descs      []string
}

func (s *NodeJSManager) Generate(r *resource.LanguageServiceResource, nodes []node.Node, nodeInfoLookup map[string]*node.Info) error {
	var err error

	tmpl, err := s.generateServiceTemplate(r, nodes, nodeInfoLookup)
	if err != nil {
		return errors.Wrapf(err, "error scaffolding functions")
	}

	err = s.generateServiceProto(tmpl)
	if err != nil {
		return errors.Wrapf(err, "error generating service protos")
	}

	err = s.generateServices(tmpl)
	if err != nil {
		return errors.Wrapf(err, "error generating services")
	}
	return nil
}

func (s *NodeJSManager) GenerateFunctionImpl(r *resource.LanguageServiceResource, nodes []node.Node) ([]Method, error) {
	// TODO breadchris generate the function implementation
	// default impl is just a console.log and return
	return nil, nil
}

func (s *NodeJSManager) generateServiceTemplate(r *resource.LanguageServiceResource, nodes []node.Node, nodeInfoLookup map[string]*node.Info) (*ServiceTemplate, error) {
	tmpl := &ServiceTemplate{
		Runtime:    strings.ToLower(r.Runtime.String()),
		Methods:    []Method{},
		DescLookup: map[string]protoreflect.MessageDescriptor{},
		EnumLookup: map[string]protoreflect.EnumDescriptor{},
		Descs:      []string{},
	}
	for _, resNode := range nodes {
		switch n := resNode.(type) {
		case *node.FunctionNode:
			// create function directory
			funcDir := path.Join("functions", n.NormalizedName())
			funcDirPath, err := s.codeRoot.GetFolder(funcDir)
			if err != nil {
				return nil, errors.Wrapf(err, "error creating function directory %s", funcDir)
			}

			info, ok := nodeInfoLookup[n.ID()]
			if !ok {
				return nil, errors.Errorf("error getting node info for %s", n.ID())
			}

			tmpl.DescLookup = lo.Assign(tmpl.DescLookup, info.Method.DescLookup)
			tmpl.EnumLookup = lo.Assign(tmpl.EnumLookup, info.Method.EnumLookup)

			method := Method{
				Name:       n.NormalizedName(),
				InputType:  string(info.Method.Input.Name()),
				OutputType: string(info.Method.Output.Name()),
			}

			err = templates.TemplateFile("node/function.index.tmpl.js", path.Join(funcDirPath, "index.js"), method)
			if err != nil {
				return nil, err
			}

			tmpl.Methods = append(tmpl.Methods, method)
		}
	}

	// TODO breadchris refactor so that we build an actual proto file and add the types to print all at once

	// format proto descs that are needed for the input and output of the nodes
	var msgStrs []string
	for _, d := range tmpl.DescLookup {
		s, err := manager.PrintMessage(d)
		if err != nil {
			return nil, errors.Wrapf(err, "error printing message")
		}
		msgStrs = append(msgStrs, s)
	}
	tmpl.Descs = append(tmpl.Descs, strings.Join(msgStrs, "\n"))

	var enumStrs []string
	for _, e := range tmpl.EnumLookup {
		s, err := manager.PrintEnum(e)
		if err != nil {
			return nil, errors.Wrapf(err, "error printing enum")
		}
		enumStrs = append(enumStrs, s)
	}
	tmpl.Descs = append(tmpl.Descs, strings.Join(enumStrs, "\n"))
	return tmpl, nil
}

// TODO breadchris should generating service protos be done for all services at once? or in each service directory?
func (s *NodeJSManager) generateServiceProto(tmpl *ServiceTemplate) error {
	protosPath, err := s.codeRoot.GetFolder("protos")
	if err != nil {
		return errors.Wrapf(err, "error getting protos folder %s", path.Join(protosPath, "protos"))
	}

	runtime := "nodejs"
	protoPath := path.Join(protosPath, fmt.Sprintf("%s.proto", runtime))
	return templates.TemplateFile("service.tmpl.proto", protoPath, tmpl)
}

func (s *NodeJSManager) generateServices(tmpl *ServiceTemplate) error {
	projectPath, err := s.codeRoot.GetFolder("/")
	if err != nil {
		return errors.Wrapf(err, "error getting project folder")
	}

	return templates.TemplateDir("node/project", projectPath, tmpl)
}
