package whatsapp

import (
	"context"

	"github.com/ElioenaiFerrari/malamute/env"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var e env.Environment = env.New()

type WhatsappService struct {
	client *twilio.RestClient
}

func NewWhatsappService(client *twilio.RestClient) *WhatsappService {
	return &WhatsappService{
		client: client,
	}
}

func (ws *WhatsappService) SendMessage(ctx context.Context, ch chan *openapi.ApiV2010Message, from, to, body string) {
	message, err := ws.client.Api.CreateMessage(&openapi.CreateMessageParams{
		From:           &from,
		To:             &to,
		Body:           &body,
		StatusCallback: &e.Twilio.WhatsappCallbackURL,
	})

	if err != nil {
		ch <- nil
		return
	}

	ch <- message
}
