package assistant

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

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

func parseMessageTemplate(vars map[string]interface{}) (string, error) {
	action := vars["action"]
	view := fmt.Sprintf("template/assistant/%s.html", action)
	t, err := template.ParseFiles(view)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, vars); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func UpdateMessageByAction(vars map[string]interface{}) (string, error) {
	switch vars["action"] {
	case "show_menu":
		return parseMessageTemplate(vars)
	default:
		return "", nil
	}
}
