package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"mws-ai/internal/dto"
	"mws-ai/internal/models"
	"mws-ai/internal/repository"
	"mws-ai/pkg/jwt"
	"mws-ai/pkg/logger"
)

var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrUserAlreadyExists = errors.New("user already exists")

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
	log := logger.Log.With().
		Str("service", "auth").
		Str("method", "Register").
		Str("email", req.Email).
		Logger()

	log.Debug().Msg("registration started")

	existing, err := s.users.FindByEmail(req.Email)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to check existing user")

		return nil, err
	}

	if existing != nil {
		log.Info().
			Msg("registration rejected: user already exists")

		return nil, ErrUserAlreadyExists
	}

	log.Debug().
		Msg("user not found, proceeding with registration")

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to hash password")

		return nil, err
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	if err := s.users.Create(user); err != nil {
		log.Error().
			Err(err).
			Msg("failed to create user")

		return nil, err
	}

	log.Info().
		Uint("user_id", user.ID).
		Msg("user registered successfully")

	return user, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	log := logger.Log.With().
		Str("service", "auth").
		Str("method", "Login").
		Str("email", req.Email).
		Logger()

	log.Debug().Msg("login attempt")

	user, err := s.users.FindByEmail(req.Email)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to find user during login")

		return nil, ErrInvalidCredentials
	}

	if user == nil {
		log.Info().
			Msg("login failed: user not found")

		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {

		log.Info().
			Uint("user_id", user.ID).
			Msg("login failed: invalid password")

		return nil, ErrInvalidCredentials
	}

	access, err := s.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		log.Error().
			Uint("user_id", user.ID).
			Err(err).
			Msg("failed to generate access token")

		return nil, err
	}

	refresh, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Error().
			Uint("user_id", user.ID).
			Err(err).
			Msg("failed to generate refresh token")

		return nil, err
	}

	log.Info().
		Uint("user_id", user.ID).
		Msg("login successful")

	return &dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
