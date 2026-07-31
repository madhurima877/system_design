package repository

import (
	"context"

	"github.com/madhurima877/chat-system/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChatRepository struct {
	Collection *mongo.Collection
}

func NewChatRepository(client *mongo.Client) *ChatRepository {
	collection := client.Database("chat_app").Collection("chat")
	return &ChatRepository{
		Collection: collection,
	}
}
func (r *ChatRepository) SaveMessage(message *models.Message) error {
	_, err := r.Collection.InsertOne(context.Background(), message)
	if err != nil {
		return err
	}

	return nil
}

func (r *ChatRepository) GetMessages(senderID, receiverID string) ([]models.Message, error) {
	filter := bson.M{
		"$or": []bson.M{
			{
				"sender_id":   senderID,
				"receiver_id": receiverID,
			},
			{
				"sender_id":   receiverID,
				"receiver_id": senderID,
			},
		},
	}
	cursor, err := r.Collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var messages []models.Message
	err = cursor.All(context.Background(), &messages)
	if err != nil {
		return nil, err
	}
	return messages, nil

}
