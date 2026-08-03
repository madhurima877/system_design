package notification

import (
	"fmt"
	"system_design/food-delivery-tracker/internal/websocket"
)

type Notification struct {
	UserID  string
	Message string
}
type Service struct {
	ch  chan Notification
	hub *websocket.Hub
}

func NewService(hub *websocket.Hub) *Service {
	return &Service{ch: make(chan Notification, 2), hub: hub}
}

func (s *Service) Start() {
	for i := 1; i <= 3; i++ {
		go func(workerID int) {
			for notification := range s.ch {
				fmt.Println("Worker", workerID, "received:", notification)

				err := s.hub.SendToUser(notification.UserID, []byte(notification.Message))
				if err != nil {
					fmt.Println(err)
				}
			}
			fmt.Println("Worker", workerID, "stopped")
		}(i)
	}

}
func (s *Service) Send(notification Notification) {
	s.ch <- notification
}
func (s *Service) Stop() {
	close(s.ch)
}
