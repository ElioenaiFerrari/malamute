package env

import (
	"log"

	"github.com/Netflix/go-env"
)

type Environment struct {
	Assistant struct {
		APIKey          string `env:"ASSISTANT_API_KEY,required=true"`
		FallbackMessage string `env:"ASSISTANT_FALLBACK_MESSAGE,required=true"`
		ID              string `env:"ASSISTANT_ID,required=true"`
		InitialMessage  string `env:"ASSISTANT_INITIAL_MESSAGE,required=true"`
		URL             string `env:"ASSISTANT_URL,required=true"`
		Version         string `env:"ASSISTANT_VERSION,required=true"`
	}
	DB struct {
		URL string `env:"DB_URL,required=true"`
	}
	SMTP struct {
		User string `env:"SMTP_USER,required=true"`
		Pass string `env:"SMTP_PASS,required=true"`
	}
	Twilio struct {
		AccountSID          string `env:"TWILIO_ACCOUNT_SID,required=true"`
		ApprovedPhone       string `env:"TWILIO_APPROVED_PHONE,required=true"`
		AuthToken           string `env:"TWILIO_AUTH_TOKEN,required=true"`
		WhatsappCallbackURL string `env:"TWILIO_WHATSAPP_CALLBACK_URL,required=true"`
	}
}

func New() Environment {
	var environment Environment

	if _, err := env.UnmarshalFromEnviron(&environment); err != nil {
		log.Fatal(err)
	}

	return environment
}
