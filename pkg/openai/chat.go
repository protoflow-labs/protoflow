package openai

import (
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/reactivex/rxgo/v2"
	"github.com/sashabaranov/go-openai"
	"sync"
)

var ChatProviderSet = wire.NewSet(
	NewChat,
	ProviderSet,
)

type ChatState struct {
	chatChan chan rxgo.Item
	context  []openai.ChatCompletionMessage
}

type ChatServer struct {
	chatLock sync.Mutex
	chats    map[string]*ChatState
	qaClient QAClient
}

func NewChat(qaClient QAClient) *ChatServer {
	return &ChatServer{
		chats:    make(map[string]*ChatState),
		qaClient: qaClient,
	}
}

func (c *ChatServer) getChatState(chat *gen.Chat) *ChatState {
	chatState, ok := c.chats[chat.Id]
	if !ok {
		c.chatLock.Lock()
		defer c.chatLock.Unlock()
		chatChan := make(chan rxgo.Item)
		chatState = &ChatState{
			chatChan: chatChan,
		}
		c.chats[chat.Id] = chatState
	}
	return chatState
}

func (c *ChatServer) Send(chatMsg *gen.SendChatRequest) (rxgo.Observable, error) {
	chatState := c.getChatState(chatMsg.Chat)
	chatState.context = append(chatState.context, openai.ChatCompletionMessage{
		Role:    "user",
		Content: chatMsg.Message,
	})
	return c.qaClient.StreamResponse(chatState.context)
}
