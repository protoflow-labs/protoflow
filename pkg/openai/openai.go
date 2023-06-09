package openai

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"io"
	"strings"
	"time"
)

type ModelDetails struct {
	MaxTokens     int
	TokensPerMsg  int
	TokensPerName int
}

var (
	models = map[string]ModelDetails{
		openai.GPT3Dot5Turbo: {
			MaxTokens:     4096,
			TokensPerMsg:  4,
			TokensPerName: -1,
		},
		openai.GPT4: {
			MaxTokens:     8000,
			TokensPerMsg:  3,
			TokensPerName: 1,
		},
	}
)

var ProviderSet = wire.NewSet(
	NewConfig,
	NewOpenAIQAClient,
	wire.Bind(new(QAClient), new(*OpenAIQAClient)),
)

type QAClient interface {
	AskWithContext(chatCtx []openai.ChatCompletionMessage, stream bool) (rxgo.Observable, error)
}

func NewOpenAIQAClient(c Config) (*OpenAIQAClient, error) {
	if _, ok := models[c.Model]; !ok {
		var availableModels []string
		for model := range models {
			availableModels = append(availableModels, model)
		}
		return nil, fmt.Errorf("model not supported: %s, available models: %s", c.Model, strings.Join(availableModels, ", "))
	}

	client := openai.NewClient(c.APIKey)
	return &OpenAIQAClient{
		client:       client,
		config:       c,
		model:        c.Model,
		modelDetails: models[c.Model],
	}, nil
}

type OpenAIQAClient struct {
	client       *openai.Client
	config       Config
	model        string
	modelDetails ModelDetails
}

var _ QAClient = &OpenAIQAClient{}

func (c *OpenAIQAClient) AskWithContext(chatCtx []openai.ChatCompletionMessage, stream bool) (rxgo.Observable, error) {
	respTokenCount, err := validateChatCtx(chatCtx, c.modelDetails)
	if err != nil {
		return nil, err
	}

	log.Debug().
		Int("respTokenCount", respTokenCount).
		Bool("stream", stream).
		Msg("Sending request to OpenAI")

	req := openai.ChatCompletionRequest{
		Model:       c.model,
		Temperature: float32(0),
		MaxTokens:   respTokenCount - 200,
		Stream:      stream,
		Messages:    chatCtx,
		// TODO breadchris need to validate the context size here
		Functions: []openai.FunctionDefinition{{
			Name:        "test",
			Description: "This is a test function",
			Parameters: &jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"count": {
						Type:        jsonschema.Number,
						Description: "total number of words in sentence",
					},
					"words": {
						Type:        jsonschema.Array,
						Description: "list of words in sentence",
						Items: &jsonschema.Definition{
							Type: jsonschema.String,
						},
					},
					"enumTest": {
						Type: jsonschema.String,
						Enum: []string{"hello", "world"},
					},
				},
			},
		}},
	}

	// context timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)

	chatStream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		cancel()
		return nil, errors.Wrapf(err, "failed to send request to OpenAI")
	}

	//resp, err := c.client.CreateChatCompletion(ctx, req)
	//if err != nil {
	//	cancel()
	//	return nil, errors.Wrapf(err, "failed to send request to OpenAI")
	//}

	msgChan := make(chan rxgo.Item)
	go func() {
		defer cancel()
		defer chatStream.Close()
		for {
			response, err := chatStream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				log.Error().Err(err).Msg("failed to get response from OpenAI")
				msgChan <- rx.NewError(err)
				break
			}

			d := response.Choices[0].Delta
			if d.FunctionCall != nil {
				log.Debug().
					Str("name", d.FunctionCall.Name).
					Str("arguments", d.FunctionCall.Arguments).
					Msg("function call received")
				continue
			}

			msgChan <- rx.NewItem(d.Content)
		}
		close(msgChan)
	}()
	return rxgo.FromEventSource(msgChan), nil
}
