package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ElioenaiFerrari/malamute/assistant"
	"github.com/ElioenaiFerrari/malamute/chat"
	"github.com/ElioenaiFerrari/malamute/email"
	"github.com/ElioenaiFerrari/malamute/env"
	"github.com/ElioenaiFerrari/malamute/sms"
	"github.com/ElioenaiFerrari/malamute/web"
	"github.com/ElioenaiFerrari/malamute/whatsapp"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/olahol/melody"
	"github.com/twilio/twilio-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/gomail.v2"
)

var e env.Environment = env.New()

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI(e.DB.URL))
	if err != nil {
		log.Fatal(err)
	}
	client.Connect(context.Background())

	db := client.Database("malamute")
	t := twilio.NewRestClientWithParams(twilio.ClientParams{
		AccountSid: e.Twilio.AccountSID,
		Password:   e.Twilio.AuthToken,
	})
	smtpClient := gomail.NewDialer("smtp.gmail.com", 587, e.SMTP.User, e.SMTP.Pass)
	chatService := chat.NewChatService(db)
	whatsappService := whatsapp.NewWhatsappService(t)
	assistantService := assistant.NewAssistantService()
	whatsappController := whatsapp.NewWhatsappController(whatsappService, assistantService, chatService)
	webController := web.NewWebController(assistantService, chatService)
	smsService := sms.NewSMSService(t)
	smsController := sms.NewSMSController(smsService)
	emailService := email.NewEmailService(smtpClient)
	emailController := email.NewEmailController(emailService)

	app := echo.New()
	v1 := app.Group("/api/v1")
	whatsappV1 := v1.Group("/whatsapp")
	smsV1 := v1.Group("/sms")
	emailV1 := v1.Group("/email")

	app.Use(middleware.Logger())
	app.Use(middleware.CORS())
	app.Use(middleware.Gzip())

	websocket := melody.New()

	websocket.HandleMessage(webController.SendMessage)

	v1.GET("/ws", func(c echo.Context) error {
		websocket.HandleRequest(c.Response().Writer, c.Request())
		return nil
	})

	whatsappV1.POST("/messages", whatsappController.SendMessage)
	whatsappV1.POST("/callback", whatsappController.Callback)
	smsV1.POST("/messages", smsController.SendMessage)
	emailV1.POST("/messages", emailController.SendMessage)

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-ch

		log.Println("shutting down...")
		os.Exit(1)
	}()

	app.Logger.Fatal(app.Start(":4000"))

}
