package grpc

import (
	"context"
	"system_design/food-delivery-tracker/internal/model"
	"system_design/food-delivery-tracker/internal/service"
	"system_design/food-delivery-tracker/proto/user"
)

type UserServer struct {
	user.UnimplementedUserServiceServer
	repo *service.UserService
}

func NewUserServer(repo *service.UserService) *UserServer {
	return &UserServer{repo: repo}
}
func (s *UserServer) CreateUser(ctx context.Context, req *user.UserRequest) (*user.UserResponse, error) {
	userdata := model.User{
		Name:  req.Name,
		Email: req.Email,
	}
	id, err := s.repo.CreateUser(userdata)
	if err != nil {
		return nil, err
	}

	return &user.UserResponse{
		Id:      int32(id),
		Message: "User Created",
	}, nil
}
