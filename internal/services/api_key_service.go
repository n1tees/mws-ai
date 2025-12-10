package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"

	"mws-ai/internal/models"
	"mws-ai/internal/repository"
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
	// 1. Генерируем "сырой" ключ — он будет показан пользователю 1 раз
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
	err := s.repo.Create(apiKey)
	if err != nil {
		return "", err
	}

	// 5. Возвращаем ключ — его нужно сохранить клиенту!
	return rawKey, nil
}

// Validate checks whether the provided API key is valid.
// Returns the user ID if valid.
func (s *APIKeyService) Validate(rawKey string) (uint, error) {
	if rawKey == "" {
		return 0, errors.New("empty api key")
	}

	// 1. Превращаем ключ в хеш
	hash := sha256.Sum256([]byte(rawKey))
	hashStr := hex.EncodeToString(hash[:])

	// 2. Пытаемся найти ключ в БД
	record, err := s.repo.FindByHash(hashStr)
	if err != nil || record == nil {
		return 0, errors.New("invalid api key")
	}

	if !record.Active {
		return 0, errors.New("api key inactive")
	}

	// 3. Обновляем время последнего использования
	now := time.Now()
	_ = s.repo.UpdateLastUsed(record.ID, &now)

	// 4. Возвращаем userID, от имени которого делается запрос
	return record.UserID, nil
}
