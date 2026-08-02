package service

import (
	"system_design/food-delivery-tracker/internal/model"
	"system_design/food-delivery-tracker/internal/repository"
	"system_design/food-delivery-tracker/internal/websocket"
)

type UserService struct {
	repo *repository.UserRepository
	hub  *websocket.Hub
}

func NewUserService(repo *repository.UserRepository, hub *websocket.Hub) *UserService {
	return &UserService{repo: repo, hub: hub}
}
func (s *UserService) CreateUser(user model.User) (int, error) {
	id, err := s.repo.CreateUser(user)
	if err != nil {
		return 0, err
	}

	return id, nil

}

func (s *UserService) NotifyUser(userID string, message string) error {
	return s.hub.SendToUser(userID, []byte(message))
}
