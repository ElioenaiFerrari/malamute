package sms

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type SMSController struct {
	smsService *SMSService
}

func NewSMSController(smsService *SMSService) *SMSController {
	return &SMSController{
		smsService: smsService,
	}
}

func (smsController *SMSController) SendMessage(c echo.Context) error {
	var params map[string]string
	wg := &sync.WaitGroup{}

	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if params["to"] == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing `to` param")
	}

	if params["body"] == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing `body` param")
	}

	wg.Add(1)
	go smsController.smsService.SendMessage(wg, params["to"], params["body"])
	wg.Wait()

	return c.NoContent(http.StatusNoContent)
}
