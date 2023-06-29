package rx

import (
	"context"
	"github.com/reactivex/rxgo/v2"
)

type ItemSink chan<- rxgo.Item

func NewItem(value any) rxgo.Item {
	return rxgo.Item{
		V: value,
	}
}

func NewError(err error) rxgo.Item {
	return rxgo.Item{
		E: err,
	}
}

func FromValues(values ...any) rxgo.Observable {
	return rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		for _, value := range values {
			next <- rxgo.Of(value)
		}
	}})
}
