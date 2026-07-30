package service

import (
	"errors"
	"time"

	"github.com/madhurima877/chat-system/internal/auth"
	"github.com/madhurima877/chat-system/internal/models"
	"github.com/madhurima877/chat-system/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *AuthService) Register(req *models.RegisterRequest) error {
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return err
	}
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,

		CreatedAt: time.Now(),
	}

	return s.repo.CreateUser(user)
}

func (s *AuthService) Login(req *models.LoginRequest) (string, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := auth.GenerateJWT(user.ID.Hex())
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	return token, nil
}
