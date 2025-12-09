package repository

import (
	"mws-ai/internal/models"

	"gorm.io/gorm"
)

type AnalysisRepository interface {
	Create(analysis *models.Analysis) error
	GetByID(id uint) (*models.Analysis, error)
	ListByUser(userID uint) ([]models.Analysis, error)
}

type analysisRepository struct {
	db *gorm.DB
}

func NewAnalysisRepository(db *gorm.DB) AnalysisRepository {
	return &analysisRepository{db: db}
}

func (r *analysisRepository) Create(analysis *models.Analysis) error {
	return r.db.Create(analysis).Error
}

func (r *analysisRepository) GetByID(id uint) (*models.Analysis, error) {
	var analysis models.Analysis
	err := r.db.Preload("Findings").
		First(&analysis, id).Error

	if err != nil {
		return nil, err
	}
	return &analysis, nil
}

func (r *analysisRepository) ListByUser(userID uint) ([]models.Analysis, error) {
	var list []models.Analysis

	err := r.db.
		Where("user_id = ?", userID).
		Order("uploaded_at DESC").
		Find(&list).Error

	if err != nil {
		return nil, err
	}

	return list, nil
}
