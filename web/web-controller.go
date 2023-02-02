package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ElioenaiFerrari/malamute/chat"
	"github.com/ElioenaiFerrari/malamute/env"
	"github.com/IBM/go-sdk-core/core"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
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

func (webController *WebController) SendMessage(c *fiber.Ctx) error {
	userMessageTimestamp := time.Now()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	body := c.FormValue("body")

	userChat, _ := webController.chatService.Find(ctx, "id", "")

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
			Text:        &body,
		},
		Context: userChat.LastMessage.Context,
	})

	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	assistantMessageTimestamp := time.Now()

	generic := assistantMessage.Output.Generic[0]
	var parsedMessage assistantv2.RuntimeResponseGeneric

	assistantMessageB, _ := sonic.Marshal(generic)
	if err := sonic.Unmarshal(assistantMessageB, &parsedMessage); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	// ch := make(chan *openapi.ApiV2010Message)

	messages := []chat.Message{
		{
			Text:      body,
			From:      chat.IssuerUser,
			Context:   nil,
			CreatedAt: userMessageTimestamp,
			Status:    chat.MessageStatusRead,
		},
		{
			Text:      *parsedMessage.Text,
			From:      chat.IssuerAssistant,
			Context:   assistantMessage.Context,
			CreatedAt: assistantMessageTimestamp,
			Status:    chat.MessageStatusSent,
		},
	}

	if _, err := webController.chatService.PushMessages(ctx, "", messages); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(http.StatusNoContent)
}

func (webController *WebController) Callback(c *fiber.Ctx) error {
	sid := c.FormValue("MessageSid")
	status := c.FormValue("MessageStatus")
	to := c.FormValue("To")

	log.Printf("web::callback message %s %s to %s", sid, status, to)
	return c.SendStatus(http.StatusNoContent)
}
