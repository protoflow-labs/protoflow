package workflow

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func (a *Activity) ExecuteGRPCNode(ctx workflow.Context, node *GRPCNode, input Input) (Result, error) {
	log.Debug().Msgf("executing node: %s", node.Service)

	g, err := getResource[GRPCResource](input.Resources)
	if err != nil {
		return Result{}, errors.Wrap(err, "error getting GRPC resource")
	}

	conn, err := grpc.Dial(g.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return Result{}, errors.Wrapf(err, "unable to connect to python server at %s", g.Host)
	}
	var reply interface{}

	err = conn.Invoke(context.Background(), node.Method, input.Params, reply)
	return Result{
		Data: reply,
	}, err
}

func (a *Activity) ExecuteRestNode(ctx workflow.Context, node *RESTNode, input Input) (Result, error) {
	log.Debug().Msgf("executing input: %v", node.Method)
	return Result{}, nil
}
