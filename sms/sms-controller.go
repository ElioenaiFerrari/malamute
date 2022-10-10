package sms

import (
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type SMSController struct {
	smsService *SMSService
}

func NewSMSController(smsService *SMSService) *SMSController {
	return &SMSController{
		smsService: smsService,
	}
}

func (smsController *SMSController) SendMessage(ctx *fiber.Ctx) error {
	var params map[string]string
	wg := &sync.WaitGroup{}

	if err := ctx.BodyParser(&params); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	if params["to"] == "" {
		return fiber.NewError(http.StatusBadRequest, "missing `to` param")
	}

	if params["body"] == "" {
		return fiber.NewError(http.StatusBadRequest, "missing `body` param")
	}

	wg.Add(1)
	go smsController.smsService.SendMessage(wg, params["to"], params["body"])
	wg.Wait()

	return ctx.SendStatus(http.StatusNoContent)
}
