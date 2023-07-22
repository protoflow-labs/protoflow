package node

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"strings"
)

type BaseNode struct {
	Name       string
	id         string
	resourceID string
}

func (n *BaseNode) NormalizedName() string {
	name := util.ToTitleCase(n.Name)
	if strings.Contains(name, ".") {
		name = strings.Split(name, ".")[1]
	}
	return name
}

func (n *BaseNode) ID() string {
	return n.id
}

func (n *BaseNode) ResourceID() string {
	return n.resourceID
}

func (n *BaseNode) Info(r graph.Resource) (*graph.Info, error) {
	return nil, nil
}

func (n *BaseNode) Represent() (string, error) {
	return "", errors.New("not implemented")
}

type RESTNode struct {
	BaseNode
	*gen.REST
}

var _ graph.Node = &RESTNode{}

type CollectionNode struct {
	BaseNode
	Collection *gen.Collection
}

var _ graph.Node = &CollectionNode{}

type BucketNode struct {
	BaseNode
	*gen.Bucket
}

var _ graph.Node = &BucketNode{}

type QueryNode struct {
	BaseNode
	Query *gen.Query
}

var _ graph.Node = &QueryNode{}

type PromptNode struct {
	BaseNode
	Prompt *gen.Prompt
}

var _ graph.Node = &PromptNode{}

type SecretNode struct {
	BaseNode
	Secret *gen.Secret
}

var _ graph.Node = &SecretNode{}

type TemplateNode struct {
	BaseNode
	Template *gen.Template
}

var _ graph.Node = &TemplateNode{}

type FileNode struct {
	BaseNode
	File *gen.File
}

var _ graph.Node = &FileNode{}
