package workflow

import (
	"context"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"testing"

	"github.com/protoflow-labs/protoflow/gen"
)

func TestRun(t *testing.T) {
	// TODO breadchris start server to listen for localhost:8080?

	r := resource.NewProto(&gen.Resource{
		Type: &gen.Resource_LanguageService{
			LanguageService: &gen.LanguageService{},
		},
	})

	n1 := node.NewFunctionNode(
		node.NewFunctionProto("test 1", r.Id),
		node.WithFunction(node.InMemoryObserver("test 1")),
	)
	n2 := node.NewFunctionNode(
		node.NewFunctionProto("test 2", r.Id),
		node.WithFunction(node.InMemoryObserver("test 2")),
	)

	a, err := Default().
		WithResource(r).
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
