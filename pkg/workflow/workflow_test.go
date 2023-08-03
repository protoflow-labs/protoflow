package workflow

import (
	"context"
	pcode "github.com/protoflow-labs/protoflow/gen/code"
	"github.com/protoflow-labs/protoflow/pkg/graph/edge"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	code2 "github.com/protoflow-labs/protoflow/pkg/graph/node/code"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"testing"
)

func TestRun(t *testing.T) {
	// TODO breadchris start server to listen for localhost:8080?

	r := code2.NewServer(base.NewNode("test 1"), code2.NewServerProto(pcode.Runtime_NODEJS).GetServer())

	n1 := code2.NewFunctionNode(
		base.NewNode("test 2"),
		code2.NewFunctionProto().GetFunction(),
		code2.WithFunction(code2.InMemoryObserver("test 1")),
	)
	n2 := code2.NewFunctionNode(
		base.NewNode("test 3"),
		code2.NewFunctionProto().GetFunction(),
		code2.WithFunction(code2.InMemoryObserver("test 2")),
	)

	a, err := Default().
		WithBuiltNodes(n1, n2).
		WithBuiltEdges(
			edge.New(edge.NewProvidesProto(r.ID(), n1.ID())),
			edge.New(edge.NewProvidesProto(r.ID(), n2.ID())),
			edge.New(edge.NewMapProto(n1.ID(), n2.ID())),
		).
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
