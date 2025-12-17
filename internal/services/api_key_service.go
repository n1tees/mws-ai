package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"

	"mws-ai/internal/models"
	"mws-ai/internal/repository"
	"mws-ai/pkg/logger"
)

var (
	ErrInvalidAPIKey  = errors.New("invalid api key")
	ErrExpiredAPIKey  = errors.New("api key expired")
	ErrInactiveAPIKey = errors.New("api key inactive")
	ErrEmptyAPIKey    = errors.New("empty api key")
)

type APIKeyService struct {
	repo repository.APIKeyRepository
}

func NewAPIKeyService(repo repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{repo: repo}
}

// Generate creates a new API key (returned ONCE)
func (s *APIKeyService) Generate(
	userID uint,
	keyType string,
	ttl time.Duration,
) (string, error) {

	log := logger.Log.With().
		Str("service", "api_key").
		Str("method", "Generate").
		Uint("user_id", userID).
		Str("type", keyType).
		Logger()

	rawKey := "mws_sk_" + uuid.New().String()

	hash := sha256.Sum256([]byte(rawKey))
	hashStr := hex.EncodeToString(hash[:])

	expiresAt := time.Now().Add(ttl)

	apiKey := &models.ApiKey{
		UserID:    userID,
		Hash:      hashStr,
		Type:      keyType,
		Active:    true,
		CreatedAt: time.Now(),
		ExpiresAt: &expiresAt,
	}

	if err := s.repo.Create(apiKey); err != nil {
		log.Error().Err(err).Msg("failed to store API key")
		return "", err
	}

	log.Info().
		Uint("api_key_id", apiKey.ID).
		Msg("API key generated")

	return rawKey, nil
}

// Validate checks API key and returns user_id
func (s *APIKeyService) Validate(rawKey string) (uint, error) {
	log := logger.Log.With().
		Str("service", "api_key").
		Str("method", "Validate").
		Logger()

	if rawKey == "" {
		return 0, ErrEmptyAPIKey
	}

	hash := sha256.Sum256([]byte(rawKey))
	hashStr := hex.EncodeToString(hash[:])

	key, err := s.repo.FindActiveByHash(hashStr)
	if err != nil {
		log.Error().Err(err).Msg("api key lookup failed")
		return 0, ErrInvalidAPIKey
	}

	if key == nil {
		return 0, ErrInvalidAPIKey
	}

	now := time.Now()
	_ = s.repo.UpdateLastUsed(key.ID, &now)

	log.Debug().
		Uint("api_key_id", key.ID).
		Uint("user_id", key.UserID).
		Str("type", key.Type).
		Msg("API key validated")

	return key.UserID, nil
}
