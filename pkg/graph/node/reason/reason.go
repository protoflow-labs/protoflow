package reason

import (
	"context"
	"encoding/json"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/reason"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/data"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	openaiclient "github.com/protoflow-labs/protoflow/pkg/openai"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
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

func NewPromptProto() *reason.Reason {
	return &reason.Reason{
		Type: &reason.Reason_Prompt{
			Prompt: &reason.Prompt{},
		},
	}
}

func (n *PromptNode) Type() (*graph.Info, error) {
	reqMsg := builder.NewMessage("Request")
	reqMsg = reqMsg.AddField(builder.NewField("message", builder.FieldTypeString()))
	req := builder.RpcTypeMessage(reqMsg, true)

	resMsg := builder.NewMessage("Response")
	resMsg = resMsg.AddField(builder.NewField("result", builder.FieldTypeString()))
	res := builder.RpcTypeMessage(resMsg, false)

	s := builder.NewService("Service")
	b := builder.NewMethod(n.NormalizedName(), req, res)
	s.AddMethod(b)

	m, err := b.Build()
	if err != nil {
		return nil, err
	}

	mthd, err := grpc.NewMethodDescriptor(m.UnwrapMethod())
	if err != nil {
		return nil, err
	}
	return &graph.Info{
		Method: mthd,
	}, nil
}

func (n *PromptNode) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	p, err := n.Provider()
	if err != nil {
		return graph.IO{}, err
	}
	r, ok := p.(*Engine)
	if !ok {
		return graph.IO{}, errors.New("error getting reason engine resource")
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

		s, err := r.QAClient.Ask(c, int(n.MinTokenCount))
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
	return graph.IO{
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

func (n *Engine) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	// TODO breadchris replace with some type of dependency injection capability
	var (
		configProvider config.Provider
		err            error
	)
	staticConfig := map[string]interface{}{
		"openai": openaiclient.NewDefaultConfig(),
	}
	p, err := n.Provider()
	if err != nil {
		return graph.IO{}, err
	}
	t, ok := p.(*data.ConfigNode)
	if !ok {
		return graph.IO{}, errors.New("error getting config node resource")
	}

	// TODO breadchris how do we handle resources that need to be initialized before others?
	configProvider, err = t.NewConfigProvider(config.Static(staticConfig))
	if err != nil {
		return graph.IO{}, errors.Wrapf(err, "failed to build config provider")
	}

	if configProvider == nil {
		return graph.IO{}, errors.New("config provider not found")
	}
	c, err := openaiclient.Wire(configProvider)
	if err != nil {
		return graph.IO{}, errors.Wrapf(err, "failed to initialize openai client")
	}
	n.QAClient = c
	return graph.IO{
		Observable: input.Observable,
	}, nil
}

func (n *Engine) Provide() ([]*gen.Node, error) {
	return []*gen.Node{NewProto("prompt", NewPromptProto())}, nil
}
