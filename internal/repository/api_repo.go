package repository

import (
	"errors"
	"time"

	"mws-ai/internal/models"
	"mws-ai/pkg/logger"

	"gorm.io/gorm"
)

type APIKeyRepository interface {
	Create(key *models.ApiKey) error
	FindActiveByHash(hash string) (*models.ApiKey, error)
	UpdateLastUsed(id uint, t *time.Time) error
}

type apiKeyRepository struct {
	db *gorm.DB
}

func NewAPIKeyRepository(db *gorm.DB) APIKeyRepository {
	return &apiKeyRepository{db: db}
}

func (r *apiKeyRepository) Create(key *models.ApiKey) error {
	if err := r.db.Create(key).Error; err != nil {
		logger.Log.Error().
			Str("repo", "api_key").
			Str("method", "Create").
			Err(err).
			Msg("failed to create API key")

		return err
	}

	logger.Log.Debug().
		Str("repo", "api_key").
		Str("method", "Create").
		Uint("api_key_id", key.ID).
		Uint("user_id", key.UserID).
		Msg("API key created")

	return nil
}

func (r *apiKeyRepository) FindActiveByHash(hash string) (*models.ApiKey, error) {
	var key models.ApiKey

	now := time.Now()

	err := r.db.
		Where("hash = ?", hash).
		Where("active = ?", true).
		Where(
			r.db.
				Where("expires_at IS NULL").
				Or("expires_at > ?", now),
		).
		First(&key).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		logger.Log.Error().
			Str("repo", "api_key").
			Str("method", "FindActiveByHash").
			Err(err).
			Msg("failed to find active API key by hash")

		return nil, err
	}

	return &key, nil
}

func (r *apiKeyRepository) UpdateLastUsed(id uint, t *time.Time) error {
	res := r.db.Model(&models.ApiKey{}).
		Where("id = ?", id).
		Update("last_used_at", t)

	if res.Error != nil {
		logger.Log.Error().
			Str("repo", "api_key").
			Str("method", "UpdateLastUsed").
			Uint("api_key_id", id).
			Err(res.Error).
			Msg("failed to update API key last_used_at")

		return res.Error
	}

	if res.RowsAffected == 0 {
		logger.Log.Debug().
			Str("repo", "api_key").
			Str("method", "UpdateLastUsed").
			Uint("api_key_id", id).
			Msg("no API key found to update last_used_at")
	}

	return nil
}
