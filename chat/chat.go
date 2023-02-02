package chat

type Chat struct {
	ID          string    `json:"id" bson:"id" validate:"required"`
	Messages    []Message `json:"messages" bson:"messages" validate:"-"`
	LastMessage *Message  `json:"last_message" bson:"last_message" validate:"-"`
}
