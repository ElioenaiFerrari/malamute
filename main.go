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
	"github.com/ElioenaiFerrari/malamute/whatsapp"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
	smsService := sms.NewSMSService(t)
	smsController := sms.NewSMSController(smsService)
	emailService := email.NewEmailService(smtpClient)
	emailController := email.NewEmailController(emailService)

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})
	v1 := app.Group("/api/v1")

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{AllowMethods: "POST,GET,OPTIONS"}))
	app.Use(logger.New())
	app.Use(cache.New())

	v1.Post("/whatsapp/messages", whatsappController.SendMessage)
	v1.Post("/whatsapp/callback", whatsappController.Callback)
	v1.Post("/sms/messages", smsController.SendMessage)
	v1.Post("/email/messages", emailController.SendEmail)

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-ch

		log.Println("shutting down...")
		os.Exit(1)
	}()

	log.Fatal(app.Listen(":4000"))

}
