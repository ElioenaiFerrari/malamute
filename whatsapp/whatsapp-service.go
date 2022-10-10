package whatsapp

import (
	"log"
	"sync"

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

func (whatsappService *WhatsappService) SendMessage(wg *sync.WaitGroup, from, to, body string) error {
	defer wg.Done()

	message, err := whatsappService.client.Api.CreateMessage(&openapi.CreateMessageParams{
		From:           &from,
		To:             &to,
		Body:           &body,
		StatusCallback: &e.Twilio.WhatsappCallbackURL,
	})

	if err != nil {
		return err
	}

	log.Printf("whatsapp::message message: %s %s", *message.Sid, *message.Status)

	return nil
}
