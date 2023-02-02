package email

import (
	"bytes"
	"fmt"
	"html/template"
	"sync"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	client *gomail.Dialer
}

func NewEmailService(client *gomail.Dialer) *EmailService {
	return &EmailService{
		client: client,
	}
}

func (emailService *EmailService) SendEmail(wg *sync.WaitGroup, view string, ctx map[string]string) error {
	defer wg.Done()

	view = fmt.Sprintf("template/%s.html", view)
	t, _ := template.ParseFiles(view)

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, ctx); err != nil {
		return err
	}

	result := tpl.String()
	msg := gomail.NewMessage()
	msg.SetHeader("Subject", ctx["subject"])
	msg.SetHeader("From", "Dachshund <no-reply@dachshund.io>")
	msg.SetHeader("To", fmt.Sprintf("%s <%s>", ctx["to_name"], ctx["to_email"]))
	msg.SetBody("text/html", result)

	return emailService.client.DialAndSend(msg)
}
