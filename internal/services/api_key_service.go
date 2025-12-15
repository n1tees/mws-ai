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
	ErrInactiveAPIKey = errors.New("api key inactive")
	ErrEmptyAPIKey    = errors.New("empty api key")
)

type APIKeyService struct {
	repo repository.APIKeyRepository
}

func NewAPIKeyService(repo repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{
		repo: repo,
	}
}

// Generate creates a new API key for a user.
// Returns the raw key string (only once!) and stores the hash in DB.
func (s *APIKeyService) Generate(userID uint) (string, error) {
	log := logger.Log.With().
		Str("service", "api_key").
		Str("method", "Generate").
		Uint("user_id", userID).
		Logger()

	log.Debug().Msg("API key generation started")

	// 1. Генерируем "сырой" ключ — показывается пользователю один раз
	rawKey := "mws_sk_" + uuid.New().String()

	// 2. Хеш SHA-256 — безопасно, ключ не хранится в открытом виде
	hash := sha256.Sum256([]byte(rawKey))
	hashStr := hex.EncodeToString(hash[:])

	// 3. Подготовка модели
	apiKey := &models.ApiKey{
		UserID:    userID,
		Hash:      hashStr,
		Active:    true,
		CreatedAt: time.Now(),
	}

	// 4. Сохраняем хеш в БД
	if err := s.repo.Create(apiKey); err != nil {
		log.Error().
			Err(err).
			Msg("failed to store API key hash")

		return "", err
	}

	log.Info().
		Uint("api_key_id", apiKey.ID).
		Msg("API key generated successfully")

	// 5. Возвращаем ключ — его нужно сохранить клиенту!
	return rawKey, nil
}

// Validate checks whether the provided API key is valid.
// Returns the user ID if valid.
func (s *APIKeyService) Validate(rawKey string) (uint, error) {
	log := logger.Log.With().
		Str("service", "api_key").
		Str("method", "Validate").
		Logger()

	if rawKey == "" {
		log.Debug().
			Msg("API key validation failed: empty key")

		return 0, ErrEmptyAPIKey
	}

	// 1. Превращаем ключ в хеш
	hash := sha256.Sum256([]byte(rawKey))
	hashStr := hex.EncodeToString(hash[:])

	log.Debug().
		Msg("API key hash calculated")

	// 2. Пытаемся найти ключ в БД
	record, err := s.repo.FindByHash(hashStr)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to lookup API key")

		return 0, ErrInvalidAPIKey
	}

	if record == nil {
		log.Info().
			Msg("API key not found")

		return 0, ErrInvalidAPIKey
	}

	if !record.Active {
		log.Info().
			Uint("api_key_id", record.ID).
			Uint("user_id", record.UserID).
			Msg("API key inactive")

		return 0, ErrInactiveAPIKey
	}

	// 3. Обновляем время последнего использования (best-effort)
	now := time.Now()
	if err := s.repo.UpdateLastUsed(record.ID, &now); err != nil {
		log.Warn().
			Uint("api_key_id", record.ID).
			Err(err).
			Msg("failed to update API key last_used_at")
	}

	log.Debug().
		Uint("api_key_id", record.ID).
		Uint("user_id", record.UserID).
		Msg("API key validated successfully")

	// 4. Возвращаем userID
	return record.UserID, nil
}
