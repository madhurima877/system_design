package service

import (
	"time"

	"github.com/madhurima877/chat-system/internal/models"
	"github.com/madhurima877/chat-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatService struct {
	repo *repository.ChatRepository
}

func NewChatService(repo *repository.ChatRepository) *ChatService {
	return &ChatService{
		repo: repo,
	}
}
func (s *ChatService) SendMessage(senderID, receiverID, content string) error {
	data := models.Message{
		ID:         primitive.NewObjectID(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  time.Now(),
	}
	return s.repo.SaveMessage(&data)

}
func (s *ChatService) GetMessages(senderID, receiverID string) ([]models.Message, error) {
	return s.repo.GetMessages(senderID, receiverID)
}
