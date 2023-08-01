package code

import (
	"context"
	"github.com/protoflow-labs/protoflow/gen/code"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"testing"
)

func TestFunctionNode(t *testing.T) {
	io := func(ctx context.Context, input graph.IO) (graph.IO, error) {
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
		return graph.IO{
			Observable: rxgo.FromChannel(output),
		}, nil
	}
	fn := NewFunctionNode(base.NewNode("test"), &code.Function{}, WithFunction(io))
	o, err := fn.Wire(context.Background(), graph.IO{
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
