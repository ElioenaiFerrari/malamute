package web

import (
	"context"
	"time"

	"github.com/ElioenaiFerrari/malamute/chat"
	"github.com/ElioenaiFerrari/malamute/env"
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

func (wc *WebController) InitialMessage(s *melody.Session) {
	// ctx := context.Background()
	// ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	// defer cancel()

	initialMessage := wc.chatService.GetInitialMessage()
	b, err := sonic.Marshal(initialMessage)
	if err != nil {
		s.Write([]byte(err.Error()))
		return
	}

	s.Write(b)
}

func (wc *WebController) SendMessage(s *melody.Session, b []byte) {
	go func() {
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

		assistantMessage, err := wc.chatService.SendMessage(ctx, chat.PlatformWeb, from, params["text"])
		if err != nil {
			s.Write([]byte(err.Error()))
			return
		}

		b, err = sonic.Marshal(assistantMessage)
		if err != nil {
			s.Write([]byte(err.Error()))
			return
		}

		s.Write(b)
	}()
}
