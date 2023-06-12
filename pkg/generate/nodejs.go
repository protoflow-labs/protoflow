package generate

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/protoflow-labs/protoflow/templates"
	"path"
	"strings"
)

type LanguageManager interface {
	Generate(r *resource.LanguageServiceResource, nodes []node.Node) error
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
	Name string
}

type FunctionTemplate struct {
	Node *node.FunctionNode
}

type ServiceTemplate struct {
	Runtime string
	Methods []Method
}

func (s *NodeJSManager) Generate(r *resource.LanguageServiceResource, nodes []node.Node) error {
	var err error

	methods, err := s.scaffoldFunctions(r, nodes)
	if err != nil {
		return errors.Wrapf(err, "error scaffolding functions")
	}

	tmpl := ServiceTemplate{
		Runtime: strings.ToLower(r.Runtime.String()),
		Methods: methods,
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

func (s *NodeJSManager) scaffoldFunctions(r *resource.LanguageServiceResource, nodes []node.Node) ([]Method, error) {
	var methods []Method
	for _, resNode := range nodes {
		switch node := resNode.(type) {
		case *node.FunctionNode:
			// create function directory
			funcDir := path.Join("functions", node.NormalizedName())
			funcDirPath, err := s.codeRoot.GetFolder(funcDir)
			if err != nil {
				return nil, errors.Wrapf(err, "error creating function directory %s", funcDir)
			}

			method := Method{
				Name: node.NormalizedName(),
			}

			err = templates.TemplateFile("node/function.index.tmpl.js", path.Join(funcDirPath, "index.js"), method)
			if err != nil {
				return nil, err
			}
			methods = append(methods, method)
		}
	}
	return methods, nil
}

// TODO breadchris should generating service protos be done for all services at once? or in each service directory?
func (s *NodeJSManager) generateServiceProto(tmpl ServiceTemplate) error {
	protosPath, err := s.codeRoot.GetFolder("protos")
	if err != nil {
		return errors.Wrapf(err, "error getting protos folder %s", path.Join(protosPath, "protos"))
	}

	runtime := "nodejs"
	protoPath := path.Join(protosPath, fmt.Sprintf("%s.proto", runtime))

	return templates.TemplateFile("service.tmpl.proto", protoPath, tmpl)
}

func (s *NodeJSManager) generateServices(tmpl ServiceTemplate) error {
	projectPath, err := s.codeRoot.GetFolder("/")
	if err != nil {
		return errors.Wrapf(err, "error getting project folder")
	}

	return templates.TemplateDir("node/project", projectPath, tmpl)
}
