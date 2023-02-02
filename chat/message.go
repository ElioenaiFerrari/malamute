package chat

import (
	"time"

	"github.com/watson-developer-cloud/go-sdk/v3/assistantv2"
)

type Message struct {
	Context   *assistantv2.MessageContextStateless `json:"context" bson:"context" validate:"-"`
	CreatedAt time.Time                            `json:"created_at" bson:"created_at" validate:"required,string"`
	From      string                               `json:"from" bson:"from" validate:"required,string,oneof='user' 'assistant'"`
	Platform  string                               `json:"platform" bson:"platform" validate:"required,oneof='web' 'whatsapp' 'sms'"`
	Status    string                               `json:"status" bson:"status" validate:"required,string"`
	Text      string                               `json:"text" bson:"text" validate:"required,string"`
}
