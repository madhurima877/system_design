package service

import (
	"system_design/food-delivery-tracker/internal/notification"
	"system_design/food-delivery-tracker/internal/repository"
	"system_design/food-delivery-tracker/internal/websocket"
)

type OrderService struct {
	repo                *repository.OrderRepository
	hub                 *websocket.Hub
	notificationService *notification.Service
}

func NewOrderService(repo *repository.OrderRepository, hub *websocket.Hub, notificationService *notification.Service) *OrderService {
	return &OrderService{repo: repo, hub: hub, notificationService: notificationService}
}
func (s *OrderService) UpdateOrderStatus(orderId int, customerID string, status string) error {
	err := s.repo.UpdateOrderStatus(orderId, status)
	if err != nil {
		return err
	}

	s.notificationService.Send(notification.Notification{UserID: customerID, Message: status})

	return nil
}
