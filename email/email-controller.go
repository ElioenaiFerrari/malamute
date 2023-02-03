package email

import (
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"
)

type EmailController struct {
	emailService *EmailService
}

func NewEmailController(emailService *EmailService) *EmailController {
	return &EmailController{
		emailService: emailService,
	}
}

func (emailController *EmailController) SendMessage(c echo.Context) error {
	var params map[string]string
	wg := &sync.WaitGroup{}

	if err := c.Bind(&params); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	if params["subject"] == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing `subject` param")
	}

	if params["title"] == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing `title` param")
	}

	if params["description"] == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing `description` param")
	}

	if params["to_name"] == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing `to_name` param")
	}

	if params["to_email"] == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing `to_email` param")
	}

	wg.Add(1)
	go emailController.emailService.SendEmail(wg, "default", params)
	wg.Wait()

	return c.NoContent(http.StatusNoContent)
}
