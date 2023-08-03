package storage

import (
	"github.com/protoflow-labs/protoflow/gen/storage"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
)

func New(b *base.Node, node *storage.Storage) graph.Node {
	switch t := node.Type.(type) {
	case *storage.Storage_Document:
		return NewDocumentNode(b, t.Document)
		//case *storage.Storage_File:
		//	return NewFileNode(b, t.File)
		//case *storage.Storage_Folder:
		//	return NewFolderNode(b, t.Folder)
	default:
		return nil
	}
}

func NewDocumentNode(b *base.Node, node *storage.Document) graph.Node {
	switch t := node.Type.(type) {
	case *storage.Document_Store:
		return NewStore(b, t.Store)
	case *storage.Document_Collection:
		return NewCollection(b, t.Collection)
	case *storage.Document_Query:
		return NewQuery(b, t.Query)
	default:
		return nil
	}
}
