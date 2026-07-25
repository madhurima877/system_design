package handler

import (
	"context"
	"log"
	"notification-system/notification-service/kafka"
	pb "notification-system/proto/notification"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
	producer *kafka.Producer
}

func NewNotificationHandler(producer *kafka.Producer) *NotificationHandler {
	return &NotificationHandler{producer: producer}
}

func (h *NotificationHandler) SendNotification(
	ctx context.Context,
	req *pb.SendNotificationRequest,
) (*pb.SendNotificationResponse, error) {

	userID := req.UserId
	message := req.Message

	log.Println(userID, message)

	err := h.producer.Publish(ctx, userID, message)
	if err != nil {
		return nil, err
	}

	return &pb.SendNotificationResponse{
		Status: "accepted",
	}, nil
}
