package node

import (
	"context"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"testing"
)

func TestFunctionNode(t *testing.T) {
	f := &gen.Node{
		Name: "test",
		Config: &gen.Node_Function{
			Function: &gen.Function{},
		},
	}
	io := func(ctx context.Context, input graph.Input) (graph.Output, error) {
		output := make(chan rxgo.Item)
		input.Observable.ForEach(func(item any) {
			log.Info().Interface("item", item).Msg("item")
			output <- rxgo.Of(item)
			close(output)
		}, func(err error) {
			log.Info().Err(err).Msg("err")
		}, func() {
			log.Info().Msg("complete")
		})
		return graph.Output{
			Observable: rxgo.FromChannel(output),
		}, nil
	}
	fn := NewFunctionNode(f, WithFunction(io))
	o, err := fn.Wire(context.Background(), graph.Input{
		Observable: rxgo.Defer([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
			next <- rxgo.Of("input")
		}}),
	})
	if err != nil {
		t.Fatal(err)
	}
	<-o.Observable.DoOnNext(func(item any) {
		log.Info().Interface("item", item).Msg("item")
	})
}
