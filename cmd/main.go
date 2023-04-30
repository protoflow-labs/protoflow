package main

import (
	"encoding/gob"
	"os"

	"github.com/protoflow-labs/protoflow/pkg/cache"
	"github.com/protoflow-labs/protoflow/pkg/cli"
	"github.com/rs/zerolog/log"
)

func main() {
	// TODO breadchris gob doesn't know how to serialize map[string]interface{}, register it, should this be a specific type?
	gob.Register(map[string]interface{}{})

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
