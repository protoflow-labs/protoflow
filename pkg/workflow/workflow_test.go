package workflow

import (
	"context"
	"github.com/protoflow-labs/protoflow/pkg/node"
	"github.com/protoflow-labs/protoflow/pkg/node/code"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"testing"

	"github.com/protoflow-labs/protoflow/gen"
)

func TestRun(t *testing.T) {
	// TODO breadchris start server to listen for localhost:8080?

	r := node.NewProto(&gen.Resource{
		Type: &gen.Resource_LanguageService{
			LanguageService: &gen.LanguageService{},
		},
	})

	n1 := code.NewFunctionNode(
		code.NewFunctionProto("test 1", r.Id),
		code.WithFunction(code.InMemoryObserver("test 1")),
	)
	n2 := code.NewFunctionNode(
		code.NewFunctionProto("test 2", r.Id),
		code.WithFunction(code.InMemoryObserver("test 2")),
	)

	a, err := Default().
		WithBuiltNodes(n1, n2).
		WithBuiltEdges(graph.Edge{
			From: n1,
			To:   n2,
		}).
		Build()
	if err != nil {
		t.Fatal(err)
	}

	input := rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		next <- rxgo.Of("input")
	}})
	obs, err := a.WireNodes(context.Background(), n1.ID(), input)
	if err != nil {
		t.Fatal(err)
	}
	<-obs.ForEach(func(item any) {
		log.Info().Interface("item", item).Msg("trace")
	}, func(err error) {
		log.Error().Err(err).Msg("err")
	}, func() {
		log.Info().Msg("complete")
	})
}
