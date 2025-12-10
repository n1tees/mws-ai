package repository

import (
	"mws-ai/internal/models"
	"time"

	"gorm.io/gorm"
)

type APIKeyRepository interface {
	Create(key *models.ApiKey) error
	FindByHash(hash string) (*models.ApiKey, error)
	UpdateLastUsed(id uint, t *time.Time) error
}

type apiKeyRepository struct {
	db *gorm.DB
}

func NewAPIKeyRepository(db *gorm.DB) APIKeyRepository {
	return &apiKeyRepository{db: db}
}

func (r *apiKeyRepository) Create(key *models.ApiKey) error {
	return r.db.Create(key).Error
}

func (r *apiKeyRepository) FindByHash(hash string) (*models.ApiKey, error) {
	var key models.ApiKey
	err := r.db.Where("hash = ?", hash).First(&key).Error
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (r *apiKeyRepository) UpdateLastUsed(id uint, t *time.Time) error {
	return r.db.Model(&models.ApiKey{}).
		Where("id = ?", id).
		Update("last_used_at", t).
		Error
}
