package email

import (
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type EmailController struct {
	emailService *EmailService
}

func NewEmailController(emailService *EmailService) *EmailController {
	return &EmailController{
		emailService: emailService,
	}
}

func (emailController *EmailController) SendEmail(c *fiber.Ctx) error {
	var params map[string]string
	wg := &sync.WaitGroup{}

	if err := c.BodyParser(&params); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	if params["subject"] == "" {
		return fiber.NewError(http.StatusBadRequest, "missing `subject` param")
	}

	if params["title"] == "" {
		return fiber.NewError(http.StatusBadRequest, "missing `title` param")
	}

	if params["description"] == "" {
		return fiber.NewError(http.StatusBadRequest, "missing `description` param")
	}

	if params["to_name"] == "" {
		return fiber.NewError(http.StatusBadRequest, "missing `to_name` param")
	}

	if params["to_email"] == "" {
		return fiber.NewError(http.StatusBadRequest, "missing `to_email` param")
	}

	wg.Add(1)
	go emailController.emailService.SendEmail(wg, "default", params)
	wg.Wait()

	return c.SendStatus(http.StatusNoContent)
}
