package node

import (
	"github.com/google/uuid"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/code"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/data"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/grpc"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/http"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/reason"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/storage"
)

// TODO breadchris make this something that can be modularized
func New(node *gen.Node) graph.Node {
	if node.Id == "" {
		node.Id = uuid.NewString()
	}
	b := base.NodeFromProto(node)
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
