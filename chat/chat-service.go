package chat

import (
	"context"
	"os"
	"time"

	"github.com/ElioenaiFerrari/malamute/assistant"
	"github.com/ElioenaiFerrari/malamute/env"
	"github.com/IBM/go-sdk-core/core"
	"github.com/bytedance/sonic"
	"github.com/watson-developer-cloud/go-sdk/v3/assistantv2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatService struct {
	db               *mongo.Database
	assistantService *assistantv2.AssistantV2
}

var e = env.New()

func NewChatService(db *mongo.Database, assistantService *assistantv2.AssistantV2) *ChatService {
	return &ChatService{
		db:               db,
		assistantService: assistantService,
	}
}

func (s *ChatService) Find(ctx context.Context, key string, value any) (*Chat, error) {
	var chat *Chat
	collection := s.db.Collection("chats")
	result := collection.FindOne(ctx, bson.D{
		primitive.E{Key: key, Value: value},
	})

	if err := result.Decode(&chat); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatService) GetInitialMessage() *Message {
	return &Message{
		From:      IssuerAssistant,
		CreatedAt: time.Now(),
		Context:   nil,
		Platform:  PlatformWeb,
		Status:    "sent",
		Text:      os.Getenv("ASSISTANT_INITIAL_MESSAGE"),
	}
}

func (*ChatService) takeAction(ch chan string, assistantMessage *assistantv2.MessageResponseStateless) {
	var parsedMessage assistantv2.RuntimeResponseGeneric
	generic := assistantMessage.Output.Generic
	if len(generic) == 0 {
		parsedMessage.Text = &e.Assistant.FallbackMessage
	} else {
		assistantMessageB, err := sonic.Marshal(generic[0])
		if err != nil {
			ch <- ""
			return
		}

		if err := sonic.Unmarshal(assistantMessageB, &parsedMessage); err != nil {
			ch <- ""
			return
		}

		ch <- *parsedMessage.Text
	}

	templateMessage, err := assistant.UpdateMessageByAction(parsedMessage.UserDefined)
	if err != nil {
		ch <- *parsedMessage.Text
		return
	}

	ch <- templateMessage
}

func (s *ChatService) update(
	ch chan *Message,
	actionCh chan string,
	assistantMessage *assistantv2.MessageResponseStateless,
	userChat *Chat,
	userMessageTimestamp time.Time,
	ctx context.Context,
	id,
	text,
	platform string,
) {
	var c *Chat
	assistantText := <-actionCh
	messages := []Message{
		{
			Text:      text,
			From:      IssuerUser,
			Context:   nil,
			CreatedAt: userMessageTimestamp,
			Status:    MessageStatusRead,
			Platform:  platform,
		},
		{
			Text:      assistantText,
			From:      IssuerAssistant,
			Context:   assistantMessage.Context,
			CreatedAt: time.Now(),
			Status:    MessageStatusSent,
			Platform:  platform,
		},
	}

	lastMessage := messages[len(messages)-1]
	if assistantText != "" {
		lastMessage.Text = assistantText
	}
	collection := s.db.Collection("chats")
	if len(userChat.Messages) == 0 {
		c = &Chat{
			ID:          id,
			Messages:    messages,
			LastMessage: &lastMessage,
		}
		if _, err := collection.InsertOne(ctx, &c); err != nil {
			ch <- nil
			return
		}

		result := collection.FindOneAndUpdate(
			ctx,
			bson.D{primitive.E{Key: "id", Value: id}},
			bson.M{"$push": bson.M{"messages": bson.M{"$each": messages}}, "$set": bson.M{"last_message": &lastMessage}},
		)

		if err := result.Decode(&c); err != nil {
			ch <- nil
			return
		}
	}

	ch <- &lastMessage
}

func (s *ChatService) SendMessage(ctx context.Context, platform, id, text string) (*Message, error) {
	userMessageTimestamp := time.Now()
	userChat, _ := s.Find(ctx, "id", id)
	if userChat == nil {
		userChat = &Chat{
			ID: id,
			LastMessage: &Message{
				Context: &assistantv2.MessageContextStateless{},
			},
		}
	}

	assistantMessage, _, err := s.assistantService.MessageStatelessWithContext(ctx, &assistantv2.MessageStatelessOptions{
		AssistantID: &e.Assistant.ID,
		Input: &assistantv2.MessageInputStateless{
			MessageType: core.StringPtr("text"),
			Text:        core.StringPtr(text),
		},
		Context: userChat.LastMessage.Context,
	})
	if err != nil {
		return nil, err
	}

	actionCh := make(chan string)
	go s.takeAction(actionCh, assistantMessage)

	messageCh := make(chan *Message)
	go s.update(messageCh, actionCh, assistantMessage, userChat, userMessageTimestamp, ctx, id, text, platform)

	return <-messageCh, nil
}
