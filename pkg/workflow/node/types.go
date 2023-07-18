package node

import (
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/util"
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

type RESTNode struct {
	BaseNode
	*gen.REST
}

var _ Node = &RESTNode{}

type CollectionNode struct {
	BaseNode
	Collection *gen.Collection
}

var _ Node = &CollectionNode{}

type BucketNode struct {
	BaseNode
	*gen.Bucket
}

var _ Node = &BucketNode{}

type QueryNode struct {
	BaseNode
	Query *gen.Query
}

var _ Node = &QueryNode{}

type PromptNode struct {
	BaseNode
	Prompt *gen.Prompt
}

var _ Node = &PromptNode{}

type ConfigNode struct {
	BaseNode
	Config *gen.Config
}

var _ Node = &ConfigNode{}

type SecretNode struct {
	BaseNode
	Secret *gen.Secret
}

var _ Node = &SecretNode{}

type TemplateNode struct {
	BaseNode
	Template *gen.Template
}

var _ Node = &TemplateNode{}

type RouteNode struct {
	BaseNode
	Route *gen.Route
}

var _ Node = &RouteNode{}

type FileNode struct {
	BaseNode
	File *gen.File
}

var _ Node = &FileNode{}
