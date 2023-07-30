package node

import (
	"github.com/google/uuid"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/node/base"
	"github.com/protoflow-labs/protoflow/pkg/node/code"
	"github.com/protoflow-labs/protoflow/pkg/node/data"
	"github.com/protoflow-labs/protoflow/pkg/node/grpc"
	"github.com/protoflow-labs/protoflow/pkg/node/http"
	"github.com/protoflow-labs/protoflow/pkg/node/reason"
	"github.com/protoflow-labs/protoflow/pkg/node/storage"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
)

// TODO breadchris make this something that can be modularized
func New(node *gen.Node) graph.Node {
	if node.Id == "" {
		node.Id = uuid.NewString()
	}
	b := base.NewNode(node)
	switch t := node.Type.(type) {
	case *gen.Node_Data:
		return data.New(b, t.Data)
	case *gen.Node_Reason:
		return reason.New(b, t.Reason)
	case *gen.Node_Grpc:
		return grpc.New(b, t.Grpc)
	case *gen.Node_Http:
		return http.New(b, t.Http)
	case *gen.Node_Storage:
		return storage.New(b, t.Storage)
	case *gen.Node_Code:
		return code.New(b, t.Code)
	default:
		return nil
	}
}
