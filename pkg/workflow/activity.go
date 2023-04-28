package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	grpcanal "github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/util"
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
	log.Debug().Msgf("executing rest: %v", node.Method)
	res, err := util.InvokeMethodOnUrl(node.Method, node.Path, input.Params)
	if err != nil {
		return Result{}, errors.Wrapf(err, "error invoking method: %s", node.Method)
	}
	return Result{
		Data: res,
	}, nil
}

func (a *Activity) ExecuteFunctionNode(ctx context.Context, node *FunctionNode, input Input) (Result, error) {
	log.Debug().Msgf("executing function node: %v", node.Function.Runtime)
	g, ok := input.Resources[LanguageServiceType].(*LanguageServiceResource)
	if !ok {
		return Result{}, fmt.Errorf("error getting GRPC resource: %s.%s", node.Function.Runtime, node.Name)
	}

	ser, err := json.Marshal(input.Params)
	if err != nil {
		return Result{}, errors.Wrapf(err, "error marshalling params: %s", node.Name)
	}

	inputData := gen.Data{
		Value: string(ser),
	}

	// TODO breadchris how is the method name formatted?
	methodName := fmt.Sprintf("protoflow.%sService.%s", node.Function.Runtime, util.ToTitleCase(node.Name))
	data, err := grpcanal.CallMethod(g.Conn, inputData, methodName)
	if err != nil {
		return Result{}, errors.Wrapf(err, "error calling method: %s", node.Name)
	}
	return Result{
		Data: data,
	}, nil
}
