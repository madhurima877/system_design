package grpc

import (
	"context"
	"time"

	"github.com/madhurima877/chat-system/internal/service"
	"github.com/madhurima877/chat-system/proto/chat"
)

type ChatServer struct {
	chat.UnimplementedChatServiceServer
	service *service.ChatService
}

func NewChatServer(service *service.ChatService) *ChatServer {
	return &ChatServer{service: service}

}
func (s *ChatServer) SendMessage(ctx context.Context, req *chat.SendMessageRequest) (*chat.SendMessageResponse, error) {

	err := s.service.SendMessage(
		req.SenderId,
		req.ReceiverId,
		req.Content,
	)
	if err != nil {
		return nil, err
	}
	return &chat.SendMessageResponse{
		Id:     "1",
		Status: "sent",
	}, nil
}

func (s *ChatServer) GetMessages(ctx context.Context, req *chat.GetMessagesRequest) (*chat.GetMessagesResponse, error) {
	messages, err := s.service.GetMessages(
		req.SenderId,
		req.ReceiverId,
	)
	if err != nil {
		return nil, err
	}
	var msgs []*chat.Message
	for _, data := range messages {
		msg := chat.Message{
			SenderId:   data.SenderID,
			ReceiverId: data.ReceiverID,
			Content:    data.Content,
			CreatedAt:  data.CreatedAt.Format(time.RFC3339),
		}
		msgs = append(msgs, &msg)

	}
	return &chat.GetMessagesResponse{
		Messages: msgs,
	}, nil
}
