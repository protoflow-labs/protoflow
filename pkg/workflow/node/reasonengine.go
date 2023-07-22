package node

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

type PromptNode struct {
	BaseNode
	Prompt *gen.Prompt
}

var _ graph.Node = &PromptNode{}

func NewPromptNode(node *gen.Node) *PromptNode {
	return &PromptNode{
		BaseNode: NewBaseNode(node),
		Prompt:   node.GetPrompt(),
	}
}
func (n *PromptNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	r, ok := input.Resource.(*resource.ReasoningEngineResource)
	if !ok {
		return graph.Output{}, fmt.Errorf("error getting reasoning engine resource: %s", n.Name)
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
