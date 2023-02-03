package web

import (
	"context"
	"time"

	"github.com/ElioenaiFerrari/malamute/assistant"
	"github.com/ElioenaiFerrari/malamute/chat"
	"github.com/ElioenaiFerrari/malamute/env"
	"github.com/IBM/go-sdk-core/core"
	"github.com/bytedance/sonic"
	"github.com/olahol/melody"
	"github.com/watson-developer-cloud/go-sdk/v3/assistantv2"
)

var e env.Environment = env.New()

type WebController struct {
	assistantService *assistantv2.AssistantV2
	chatService      *chat.ChatService
}

func NewWebController(
	assistantService *assistantv2.AssistantV2,
	chatService *chat.ChatService,
) *WebController {
	return &WebController{
		assistantService: assistantService,
		chatService:      chatService,
	}
}

func (webController *WebController) SendMessage(s *melody.Session, b []byte) {
	userMessageTimestamp := time.Now()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var params map[string]string

	if err := sonic.Unmarshal(b, &params); err != nil {
		s.Write([]byte(err.Error()))
		return
	}

	if params["text"] == "" {
		s.Write([]byte("missing `text` param"))
		return
	}

	from := "+5527999152059"

	userChat, _ := webController.chatService.Find(ctx, "id", from)

	if userChat == nil {
		userChat = &chat.Chat{
			LastMessage: &chat.Message{
				Context: &assistantv2.MessageContextStateless{},
			},
		}
	}

	assistantMessage, _, err := webController.assistantService.MessageStatelessWithContext(ctx, &assistantv2.MessageStatelessOptions{
		AssistantID: &e.Assistant.ID,
		Input: &assistantv2.MessageInputStateless{
			MessageType: core.StringPtr("text"),
			Text:        core.StringPtr(params["text"]),
		},
		Context: userChat.LastMessage.Context,
	})

	if err := assistant.TakeAction(assistantMessage.Context); err != nil {
		s.Write([]byte(err.Error()))
		return
	}

	if err != nil {
		s.Write([]byte(err.Error()))
		return
	}

	assistantMessageTimestamp := time.Now()

	generic := assistantMessage.Output.Generic[0]
	var parsedMessage assistantv2.RuntimeResponseGeneric

	assistantMessageB, _ := sonic.Marshal(generic)
	if err := sonic.Unmarshal(assistantMessageB, &parsedMessage); err != nil {
		s.Write([]byte(err.Error()))
		return
	}

	messages := []chat.Message{
		{
			Text:      params["text"],
			From:      chat.IssuerUser,
			Context:   nil,
			CreatedAt: userMessageTimestamp,
			Status:    chat.MessageStatusRead,
			Platform:  chat.PlatformWeb,
		},
		{
			Text:      *parsedMessage.Text,
			From:      chat.IssuerAssistant,
			Context:   assistantMessage.Context,
			CreatedAt: assistantMessageTimestamp,
			Status:    chat.MessageStatusSent,
			Platform:  chat.PlatformWeb,
		},
	}

	if _, err := webController.chatService.PushMessages(ctx, from, messages); err != nil {
		s.Write([]byte(err.Error()))
		return
	}

	b, err = sonic.Marshal(assistantMessage)
	if err != nil {
		s.Write([]byte(err.Error()))
		return
	}

	s.Write(b)
}
