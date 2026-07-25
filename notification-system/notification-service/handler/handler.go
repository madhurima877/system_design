package handler

import (
	"context"
	"log"
	pb "notification-system/proto/notification"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (h *NotificationHandler) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	userID := req.UserId
	message := req.Message
	log.Println(userID, message)
	return &pb.SendNotificationResponse{
		Status: "accepted",
	}, nil
}
