package main

import (
	"os"

	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
	"github.com/protoflow-labs/protoflow/pkg/cli"
	"github.com/rs/zerolog/log"
)

func main() {
	cacheConfig := cache.Config{
		Name: ".protoflow",
	}

	app, err := cli.Wire(cacheConfig)
	if err != nil {
		log.Error().Msgf("%+v\n", err)
		return
	}

	if err := app.Run(os.Args); err != nil {
		log.Error().Msgf("%+v\n", err)
	}
}
