package generate

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/code"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/templates"
	"path"
	"strings"
)

type LanguageManager interface {
	GenerateGRPCService(r *code.Server) error
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
	Node *code.FunctionNode
}

type ServiceTemplate struct {
	Runtime string
	Methods []Method
}

func (s *NodeJSManager) GenerateGRPCService(r *code.Server) error {
	var err error

	tmpl, err := s.generateServiceTemplate(r)
	if err != nil {
		return errors.Wrapf(err, "error scaffolding functions")
	}

	// TODO breadchris this needs to run once for the project
	//err = s.generateServiceProto(tmpl)
	//if err != nil {
	//	return errors.Wrapf(err, "error generating service protos")
	//}

	err = s.generateServices(tmpl)
	if err != nil {
		return errors.Wrapf(err, "error generating services")
	}
	return nil
}

func (s *NodeJSManager) UpdateNodeType(n graph.Node, nodeInfo *graph.Info) error {
	funcDirPath, err := s.codeRoot.GetFolder("protos")
	if err != nil {
		return errors.Wrapf(err, "error getting protos dir")
	}
	fileDescs, err := grpc.ParseProtoDir(funcDirPath)
	if err != nil {
		return errors.Wrapf(err, "error parsing protos")
	}

	updates, err := generateProtoType(fileDescs, n, nodeInfo)
	if err != nil {
		return errors.Wrapf(err, "error generating proto type")
	}
	for _, update := range updates {
		err = writeProtoFile(path.Join(funcDirPath, update.FilePath), update.FileDesc)
		if err != nil {
			return errors.Wrapf(err, "error writing proto file")
		}
	}
	return nil
}

func (s *NodeJSManager) GenerateFunctionImpl(r *code.Server, n graph.Node) error {
	switch n.(type) {
	case *code.FunctionNode:
		// create function directory
		funcDir := path.Join("functions", n.NormalizedName())
		funcDirPath, err := s.codeRoot.GetFolder(funcDir)
		if err != nil {
			return errors.Wrapf(err, "error creating function directory %s", funcDir)
		}

		method := Method{
			Name: n.NormalizedName(),
		}

		err = templates.TemplateFile("node/function.index.tmpl.js", path.Join(funcDirPath, "index.js"), method)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (s *NodeJSManager) generateServiceTemplate(r *code.Server) (*ServiceTemplate, error) {
	tmpl := &ServiceTemplate{
		Runtime: strings.ToLower(r.Runtime.String()),
		Methods: []Method{},
	}
	for _, resNode := range r.Dependents() {
		switch n := resNode.(type) {
		case *code.FunctionNode:

			method := Method{
				Name: n.NormalizedName(),
			}
			tmpl.Methods = append(tmpl.Methods, method)
		}
	}
	return tmpl, nil
}

// TODO breadchris should generating service protos be done for all services at once? or in each service directory?
// TODO breadchris this should only be run once when a project is initialized
func (s *NodeJSManager) generateServiceProto(tmpl *ServiceTemplate) error {
	protosPath, err := s.codeRoot.GetFolder("protos")
	if err != nil {
		return errors.Wrapf(err, "error getting protos folder %s", path.Join(protosPath, "protos"))
	}

	protoPath := path.Join(protosPath, fmt.Sprintf("%s.proto", tmpl.Runtime))
	return templates.TemplateFile("service.tmpl.proto", protoPath, tmpl)
}

func (s *NodeJSManager) generateServices(tmpl *ServiceTemplate) error {
	projectPath, err := s.codeRoot.GetFolder("/")
	if err != nil {
		return errors.Wrapf(err, "error getting project folder")
	}

	return templates.TemplateDir("node/project", projectPath, tmpl)
}
