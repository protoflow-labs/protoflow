package rx

import (
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
)

func LogObserver(name string, o rxgo.Observable) {
	o.ForEach(func(item any) {
		log.Debug().Str("name", name).Interface("item", item).Msg("received item")
	}, func(err error) {
		log.Error().Str("name", name).Err(err).Msg("received error")
	}, func() {
		log.Debug().Str("name", name).Msg("received complete")
	})
}
