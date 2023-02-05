package chat

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatService struct {
	db *mongo.Database
}

func NewChatService(db *mongo.Database) *ChatService {
	return &ChatService{
		db: db,
	}
}

func (s *ChatService) Find(ctx context.Context, key string, value any) (*Chat, error) {
	var chat *Chat
	collection := s.db.Collection("chat")
	result := collection.FindOne(ctx, bson.D{
		primitive.E{Key: key, Value: value},
	})

	if err := result.Decode(&chat); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatService) PushMessages(ctx context.Context, id string, messages []Message) (*Chat, error) {
	var c *Chat
	lastMessage := messages[len(messages)-1]
	collection := s.db.Collection("chat")

	_, err := s.Find(ctx, "id", id)

	if err != nil {
		c = &Chat{
			ID:          id,
			Messages:    messages,
			LastMessage: &lastMessage,
		}

		if _, err := collection.InsertOne(ctx, &c); err != nil {
			return nil, err
		}

		return c, nil
	}

	result := collection.FindOneAndUpdate(
		ctx,
		bson.D{primitive.E{Key: "id", Value: id}},
		bson.M{"$push": bson.M{"messages": bson.M{"$each": messages}}, "$set": bson.M{"last_message": &lastMessage}},
	)

	if err := result.Decode(&c); err != nil {
		return nil, err
	}

	return c, nil
}
