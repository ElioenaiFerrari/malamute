package whatsapp

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ElioenaiFerrari/malamute/chat"
	"github.com/labstack/echo/v4"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type WhatsappController struct {
	whatsappService *WhatsappService
	chatService     *chat.ChatService
}

func NewWhatsappController(
	whatsappService *WhatsappService,
	chatService *chat.ChatService,
) *WhatsappController {
	return &WhatsappController{
		whatsappService: whatsappService,
		chatService:     chatService,
	}
}

func (wc *WhatsappController) SendMessage(c echo.Context) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	from := c.FormValue("From")
	to := c.FormValue("To")
	body := c.FormValue("Body")

	assistantMessage, err := wc.chatService.SendMessage(ctx, chat.PlatformWhatsapp, RawPhone(from), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ch := make(chan *openapi.ApiV2010Message)
	go wc.whatsappService.SendMessage(ctx, ch, to, from, assistantMessage.Text)
	defer func() {
		<-ch
	}()

	return c.NoContent(http.StatusNoContent)
}

func (wc *WhatsappController) Callback(c echo.Context) error {
	sid := c.FormValue("MessageSid")
	status := c.FormValue("MessageStatus")
	to := c.FormValue("To")

	log.Printf("whatsapp::callback message %s %s to %s", sid, status, to)
	return c.NoContent(http.StatusNoContent)
}
