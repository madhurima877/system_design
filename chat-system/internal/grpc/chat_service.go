package grpc

import (
	"context"

	"github.com/madhurima877/chat-system/proto/chat"
)

type ChatServer struct {
	chat.UnimplementedChatServiceServer
}

func NewChatServer() *ChatServer {
	return &ChatServer{}

}
func (s *ChatServer) SendMessage(ctx context.Context, req *chat.SendMessageRequest) (*chat.SendMessageResponse, error) {
	return &chat.SendMessageResponse{
		Id:     "1",
		Status: "sent",
	}, nil
}
