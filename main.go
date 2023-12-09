package main

//go:generate npx buf generate proto
//go:generate go run github.com/google/wire/cmd/wire ./...
//go:generate protoc --jsonschema_out=pkg/llm/schemas --proto_path=proto proto/ai.proto

import (
	"encoding/gob"
	"github.com/protoflow-labs/protoflow/pkg/cli"
	"os"

	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/rs/zerolog/log"
)

func main() {
	// TODO breadchris gob doesn't know how to serialize map[string]interface{}, register it, should this be a specific type?
	gob.Register(map[string]interface{}{})

	app, err := cli.Wire(bucket.NewDefaultConfig())
	if err != nil {
		log.Error().Msgf("%+v\n", err)
		return
	}

	if err := app.Run(os.Args); err != nil {
		log.Error().Msgf("%+v\n", err)
	}
}
