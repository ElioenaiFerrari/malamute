package assistant

import (
	"log"

	"github.com/ElioenaiFerrari/malamute/env"
	"github.com/IBM/go-sdk-core/core"
	"github.com/watson-developer-cloud/go-sdk/v3/assistantv2"
)

var e env.Environment = env.New()

func NewAssistantService() *assistantv2.AssistantV2 {
	assistant, err := assistantv2.NewAssistantV2(&assistantv2.AssistantV2Options{
		URL:     e.Assistant.URL,
		Version: &e.Assistant.Version,
		Authenticator: &core.IamAuthenticator{
			ApiKey:                 e.Assistant.APIKey,
			DisableSSLVerification: true,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	return assistant
}

func TakeAction(context *assistantv2.MessageContextStateless) error {
	userVars := context.Skills["actions skill"].UserDefined
	switch userVars["action"] {
	}
	return nil
}
