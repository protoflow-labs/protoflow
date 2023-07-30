package reason

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen/reason"
	"github.com/protoflow-labs/protoflow/pkg/node/base"
	"github.com/protoflow-labs/protoflow/pkg/node/data"
	openaiclient "github.com/protoflow-labs/protoflow/pkg/openai"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/config"
)

type PromptNode struct {
	*base.Node
	*reason.Prompt
}

var _ graph.Node = &PromptNode{}

func NewPromptNode(b *base.Node, node *reason.Prompt) *PromptNode {
	return &PromptNode{
		Node:   b,
		Prompt: node,
	}
}

func (n *PromptNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	p, err := n.Provider()
	if err != nil {
		return graph.Output{}, err
	}
	r, ok := p.(*Engine)
	if !ok {
		return graph.Output{}, errors.New("error getting reason engine resource")
	}

	log.Info().
		Str("name", n.Name).
		Msg("setting up prompt node")

	outputStream := make(chan rxgo.Item)

	// TODO breadchris how should context be handled here?
	c := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: n.Prompt.Prompt,
		},
	}

	input.Observable.ForEach(func(item any) {
		log.Debug().
			Str("name", n.NormalizedName()).
			Interface("item", item).
			Msg("executing prompt node")

		var normalizedItem string
		switch t := item.(type) {
		case string:
			normalizedItem = t
		default:
			c, err := json.Marshal(item)
			if err != nil {
				outputStream <- rx.NewError(errors.Wrapf(err, "error marshalling input: %s", n.NormalizedName()))
				return
			}
			normalizedItem = string(c)
		}

		c = append(c, openai.ChatCompletionMessage{
			Role:    "user",
			Content: normalizedItem,
		})

		s, err := r.QAClient.Ask(c)
		if err != nil {
			outputStream <- rx.NewError(errors.Wrapf(err, "error executing prompt: %s", n.NormalizedName()))
		}

		// TODO breadchris react to a function call on s.FunctionCall

		// TODO breadchris this should be a static type. This is a brittle type that maps to workflow.go:133
		outputStream <- rx.NewItem(map[string]any{
			"result": s,
		})
	}, func(err error) {
		outputStream <- rx.NewError(err)
	}, func() {
		close(outputStream)
	})
	return graph.Output{
		Observable: rxgo.FromChannel(outputStream, rxgo.WithPublishStrategy()),
	}, nil
}

type Engine struct {
	*base.Node
	*reason.Engine
	QAClient openaiclient.QAClient
}

func NewEngineNode(b *base.Node, node *reason.Engine) *Engine {
	return &Engine{
		Node:   b,
		Engine: node,
	}
}

func (n *Engine) Init() (func(), error) {
	// TODO breadchris replace with some type of dependency injection capability
	var (
		configProvider config.Provider
		err            error
	)
	staticConfig := map[string]interface{}{
		"openai": openaiclient.NewDefaultConfig(),
	}
	for _, n := range n.Predecessors() {
		switch t := n.(type) {
		case *data.ConfigNode:
			// TODO breadchris how do we handle resources that need to be initialized before others?
			configProvider, err = t.NewConfigProvider(config.Static(staticConfig))
			if err != nil {
				return nil, errors.Wrapf(err, "failed to build config provider")
			}
		}
	}
	if configProvider == nil {
		return nil, errors.New("config provider not found")
	}
	c, err := openaiclient.Wire(configProvider)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initialize openai client")
	}
	n.QAClient = c
	return nil, nil
}

func (n *Engine) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	return graph.Output{
		Observable: input.Observable,
	}, nil
}
