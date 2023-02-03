package whatsapp

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ElioenaiFerrari/malamute/chat"
	"github.com/IBM/go-sdk-core/core"
	"github.com/bytedance/sonic"
	"github.com/labstack/echo/v4"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/watson-developer-cloud/go-sdk/v3/assistantv2"
)

type WhatsappController struct {
	whatsappService  *WhatsappService
	assistantService *assistantv2.AssistantV2
	chatService      *chat.ChatService
}

func NewWhatsappController(
	whatsappService *WhatsappService,
	assistantService *assistantv2.AssistantV2,
	chatService *chat.ChatService,
) *WhatsappController {
	return &WhatsappController{
		whatsappService:  whatsappService,
		assistantService: assistantService,
		chatService:      chatService,
	}
}

func (whatsappController *WhatsappController) SendMessage(c echo.Context) error {
	userMessageTimestamp := time.Now()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	from := c.FormValue("From")
	to := c.FormValue("To")
	body := c.FormValue("Body")

	userChat, _ := whatsappController.chatService.Find(ctx, "id", RawPhone(from))

	if userChat == nil {
		userChat = &chat.Chat{
			LastMessage: &chat.Message{
				Context: &assistantv2.MessageContextStateless{},
			},
		}
	}

	assistantMessage, _, err := whatsappController.assistantService.MessageStatelessWithContext(ctx, &assistantv2.MessageStatelessOptions{
		AssistantID: &e.Assistant.ID,
		Input: &assistantv2.MessageInputStateless{
			MessageType: core.StringPtr("text"),
			Text:        &body,
		},
		Context: userChat.LastMessage.Context,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	assistantMessageTimestamp := time.Now()

	generic := assistantMessage.Output.Generic[0]
	var parsedMessage assistantv2.RuntimeResponseGeneric

	assistantMessageB, _ := sonic.Marshal(generic)
	if err := sonic.Unmarshal(assistantMessageB, &parsedMessage); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ch := make(chan *openapi.ApiV2010Message)
	go whatsappController.whatsappService.SendMessage(ctx, ch, to, from, *parsedMessage.Text)
	defer func() {
		<-ch
	}()

	messages := []chat.Message{
		{
			Context:   nil,
			CreatedAt: userMessageTimestamp,
			From:      chat.IssuerUser,
			Platform:  chat.PlatformWhatsapp,
			Status:    chat.MessageStatusRead,
			Text:      body,
		},
		{
			Context:   assistantMessage.Context,
			CreatedAt: assistantMessageTimestamp,
			From:      chat.IssuerAssistant,
			Platform:  chat.PlatformWhatsapp,
			Status:    chat.MessageStatusSent,
			Text:      *parsedMessage.Text,
		},
	}

	if _, err := whatsappController.chatService.PushMessages(ctx, RawPhone(from), messages); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (whatsappController *WhatsappController) Callback(c echo.Context) error {
	sid := c.FormValue("MessageSid")
	status := c.FormValue("MessageStatus")
	to := c.FormValue("To")

	log.Printf("whatsapp::callback message %s %s to %s", sid, status, to)
	return c.NoContent(http.StatusNoContent)
}
