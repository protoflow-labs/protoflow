package workflow

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	grpcanal "github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/workflow"
)

type Activity struct{}

var (
	ErrResourceNotFound = errors.New("resource not found")
)

type ProtoType struct{}

func (a *Activity) ExecuteGRPCNode(ctx workflow.Context, node *GRPCNode, input Input) (Result, error) {
	log.Info().Msgf("executing node: %s", node.Service)

	g, ok := input.Resources[GRPCResourceType].(*GRPCResource)
	if !ok {
		return Result{}, fmt.Errorf("error getting GRPC resource: %s.%s", node.Service, node.Method)
	}

	// TODO breadchris is this how you get the fully qualified name?
	methodName := node.Service + "." + node.Method
	if node.Package != "" {
		methodName = node.Package + "." + methodName
	}
	data, err := grpcanal.CallMethod(g.Conn, &input, methodName)
	if err != nil {
		return Result{}, errors.Wrapf(err, "error calling method: %s", node.Method)
	}
	return Result{
		Data: data,
	}, fmt.Errorf("method not found: %s", node.Method)
}

func (a *Activity) ExecuteRestNode(ctx workflow.Context, node *RESTNode, input Input) (Result, error) {
	log.Debug().Msgf("TODO executing input: %v", node.Method)
	return Result{}, nil
}

func (a *Activity) ExecuteFunctionNode(ctx context.Context, node *FunctionNode, input Input) (Result, error) {
	log.Debug().Msgf("executing input: %v", node.Function.Runtime)
	g, ok := input.Resources[LanguageServiceType].(*LanguageServiceResource)
	if !ok {
		return Result{}, fmt.Errorf("error getting GRPC resource: %s.%s", node.Function.Runtime, node.Name)
	}

	// TODO breadchris how is the method name formatted?
	methodName := node.Function.Runtime + "." + node.Name
	data, err := grpcanal.CallMethod(g.Conn, &input, methodName)
	if err != nil {
		return Result{}, errors.Wrapf(err, "error calling method: %s", node.Name)
	}
	return Result{
		Data: data,
	}, nil
}
