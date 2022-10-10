package whatsapp

import (
	"log"
	"net/http"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/watson-developer-cloud/go-sdk/v2/assistantv2"
)

type WhatsappController struct {
	whatsappService  *WhatsappService
	assistantService *assistantv2.AssistantV2
}

func NewWhatsappController(whatsappService *WhatsappService, assistantService *assistantv2.AssistantV2) *WhatsappController {
	return &WhatsappController{
		whatsappService:  whatsappService,
		assistantService: assistantService,
	}
}

func (whatsappController *WhatsappController) SendMessage(ctx *fiber.Ctx) error {
	wg := &sync.WaitGroup{}
	from := ctx.FormValue("From")
	to := ctx.FormValue("To")
	body := ctx.FormValue("Body")
	messageType := "text"

	message, _, err := whatsappController.assistantService.MessageStateless(&assistantv2.MessageStatelessOptions{
		AssistantID: &e.Assistant.ID,
		Input: &assistantv2.MessageInputStateless{
			MessageType: &messageType,
			Text:        &body,
		},
	})

	generic := message.Output.Generic[0]
	var parsedMessage assistantv2.RuntimeResponseGeneric

	messageB, _ := sonic.Marshal(generic)
	sonic.Unmarshal(messageB, &parsedMessage)

	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	wg.Add(1)
	go whatsappController.whatsappService.SendMessage(wg, to, from, *parsedMessage.Text)
	wg.Wait()

	return ctx.SendStatus(http.StatusNoContent)
}

func (whatsappController *WhatsappController) Callback(ctx *fiber.Ctx) error {
	sid := ctx.FormValue("MessageSid")
	status := ctx.FormValue("MessageStatus")

	log.Printf("whatsapp::callback message: %s %s", sid, status)
	return ctx.SendStatus(http.StatusNoContent)
}
