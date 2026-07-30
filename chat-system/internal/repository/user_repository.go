package repository

import (
	"context"

	"github.com/madhurima877/chat-system/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) *UserRepository {
	db := client.Database("chat_app")
	collection := db.Collection("users")
	return &UserRepository{collection: collection}
}

func (m *UserRepository) CreateUser(user *models.User) error {
	_, err := m.collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	return nil

}

func (m *UserRepository) FindByEmail(email string) (*models.User, error) {
	filter := bson.M{
		"email": email,
	}
	var user models.User
	result := m.collection.FindOne(context.Background(), filter)
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

