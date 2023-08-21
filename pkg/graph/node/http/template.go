package http

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/http"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/storage"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/reactivex/rxgo/v2"
	"html/template"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

type TemplateFSNode struct {
	*base.Node
	*http.TemplateFS
}

var _ graph.Node = &TemplateFSNode{}

func NewTemplateFSNode(b *base.Node, node *http.TemplateFS) *TemplateFSNode {
	return &TemplateFSNode{
		Node:       b,
		TemplateFS: node,
	}
}

func nodesFromFiles(u string) ([]*gen.Node, error) {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse url %s", u)
	}

	var nodes []*gen.Node
	err = filepath.WalkDir(parsedUrl.Path, func(p string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(parsedUrl.Path, p)
		if err != nil {
			return errors.Wrapf(err, "failed to get relative path")
		}
		// TODO breadchris templates don't accept any paths it seems. how do we want to structure this to avoid template name collisions?
		nodes = append(nodes, NewProto(rel, NewTemplateProto(path.Base(rel))))
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to walk dir %s", parsedUrl.Path)
	}
	return nodes, nil
}

func (n *TemplateFSNode) Provide() ([]*gen.Node, error) {
	p, err := n.Provider()
	if err != nil {
		return nil, errors.Wrapf(err, "error getting provider")
	}
	f, ok := p.(*storage.Folder)
	if !ok {
		return nil, errors.Wrapf(err, "error getting folder")
	}
	return nodesFromFiles(f.Url)
}

func (n *TemplateFSNode) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	return input, nil
}

func (n *TemplateFSNode) templates(name string) (*template.Template, error) {
	p, err := n.Provider()
	if err != nil {
		return nil, errors.Wrapf(err, "error getting provider")
	}
	f, ok := p.(*storage.Folder)
	if !ok {
		return nil, errors.Wrapf(err, "error getting folder")
	}
	parsedUrl, err := url.Parse(f.Url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse url %s", f.Url)
	}
	tmplPath := path.Join(parsedUrl.Path, n.Path)
	return template.New(name).ParseFS(os.DirFS(tmplPath), "**/*.html")
}

type TemplateNode struct {
	*base.Node
	*http.Template
}

var _ graph.Node = &TemplateNode{}

func NewTemplateNode(b *base.Node, node *http.Template) *TemplateNode {
	return &TemplateNode{
		Node:     b,
		Template: node,
	}
}

func NewTemplateProto(t string) *http.HTTP {
	return &http.HTTP{
		Type: &http.HTTP_Template{
			Template: &http.Template{
				Name: t,
			},
		},
	}
}

func (n *TemplateNode) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	output := make(chan rxgo.Item)

	p, err := n.Provider()
	if err != nil {
		return graph.IO{}, errors.Wrapf(err, "error getting provider")
	}
	t, ok := p.(*TemplateFSNode)
	if !ok {
		return graph.IO{}, errors.Wrapf(err, "error getting folder")
	}
	tmpls, err := t.templates(n.Template.Name)
	if err != nil {
		return graph.IO{}, errors.Wrapf(err, "error getting templates")
	}

	input.Observable.ForEach(func(item any) {
		b := &bytes.Buffer{}
		err = tmpls.Execute(b, item)
		if err != nil {
			output <- rx.NewError(err)
			return
		}
		resp := &http.Response{
			Headers: []*http.Header{},
			Body:    b.Bytes(),
		}
		output <- rx.NewItem(resp)
	}, func(err error) {
		output <- rx.NewError(err)
	}, func() {
		close(output)
	})

	return graph.IO{
		Observable: rxgo.FromChannel(output, rxgo.WithPublishStrategy()),
	}, nil
}
