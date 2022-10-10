package sms

import (
	"log"
	"sync"

	"github.com/ElioenaiFerrari/malamute/env"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSService struct {
	client *twilio.RestClient
}

var e env.Environment = env.New()

func NewSMSService(client *twilio.RestClient) *SMSService {
	return &SMSService{
		client: client,
	}
}

func (smsService *SMSService) SendMessage(wg *sync.WaitGroup, to, body string) error {
	defer wg.Done()

	message, err := smsService.client.Api.CreateMessage(&openapi.CreateMessageParams{
		From: &e.Twilio.ApprovedPhone,
		To:   &to,
		Body: &body,
	})

	if err != nil {
		return err
	}

	log.Printf("sms::message %smsService %smsService", *message.Sid, *message.Status)

	return nil
}
