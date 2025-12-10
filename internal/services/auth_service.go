package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"mws-ai/internal/dto"
	"mws-ai/internal/models"
	"mws-ai/internal/repository"
	"mws-ai/pkg/jwt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type AuthService struct {
	users repository.UserRepository
	jwt   *jwt.JWTManager
}

func NewAuthService(users repository.UserRepository, jwt *jwt.JWTManager) *AuthService {
	return &AuthService{
		users: users,
		jwt:   jwt,
	}
}

func (s *AuthService) Register(req dto.RegisterRequest) (*models.User, error) {
	// Проверка на существующего пользователя
	_, err := s.users.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	err = s.users.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.users.FindByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	access, err := s.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refresh, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
