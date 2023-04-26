package workflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	grpcanal "github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type Activity struct{}

func getResource[T any](resources map[string]any) (*T, error) {
	var instance *T
	for _, r := range resources {
		log.Debug().Msgf("resource: %s", r)
		switch r.(type) {
		case *T:
			instance = r.(*T)
		}
	}
	if instance == nil {
		return nil, errors.New("resource not found")
	}
	return instance, nil
}

type ProtoType struct{}

func (a *Activity) ExecuteGRPCNode(ctx workflow.Context, node *GRPCNode, input Input) (Result, error) {
	log.Info().Msgf("executing node: %s", node.Service)

	g, ok := input.Resources["grpc"].(*GRPCResource)
	if !ok {
		return Result{}, fmt.Errorf("error getting GRPC resource: %s.%s", node.Service, node.Method)
	}

	requestFunc := func(m proto.Message) error {
		req, err := json.Marshal(input.Params)
		if err != nil {
			return err
		}

		if err := jsonpb.Unmarshal(bytes.NewReader(req), m); err != nil {
			return status.Errorf(codes.InvalidArgument, err.Error())
		}
		return io.EOF
	}
	var headers []string
	methods, err := grpcanal.AllMethodsViaReflection(context.Background(), g.Conn)
	if err != nil {
		return Result{}, errors.Wrap(err, "error getting all methods")
	}

	// TODO breadchris is this how you get the fully qualified name?
	methodName := node.Service + "." + node.Method
	if node.Package != "" {
		methodName = node.Package + "." + methodName
	}

	// TODO breadchris do we have to do this to make this call?
	for _, m := range methods {
		if m.GetFullyQualifiedName() == methodName {
			descSource, err := grpcurl.DescriptorSourceFromFileDescriptors(m.GetFile())
			if err != nil {
				return Result{}, errors.Wrap(err, "error getting descriptor source")
			}
			result := grpcanal.RpcResult{
				DescSource: descSource,
			}
			if err := grpcurl.InvokeRPC(context.Background(), descSource, g.Conn, methodName, headers, &result, requestFunc); err != nil {
				return Result{}, err
			}

			if len(result.Responses) == 0 {
				return Result{}, fmt.Errorf("no responses received")
			}

			resp := result.Responses[0]
			var data interface{}
			err = json.Unmarshal(resp.Data, &data)
			if err != nil {
				return Result{}, errors.Wrapf(err, "error unmarshalling response: %s", string(resp.Data))
			}

			return Result{
				Data: data,
			}, err
		}
	}
	return Result{}, fmt.Errorf("method not found: %s", node.Method)
}

func (a *Activity) ExecuteRestNode(ctx workflow.Context, node *RESTNode, input Input) (Result, error) {
	log.Debug().Msgf("executing input: %v", node.Method)
	return Result{}, nil
}
