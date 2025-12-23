package repository

import (
	"errors"

	"mws-ai/internal/models"
	"mws-ai/pkg/logger"

	"gorm.io/gorm"
)

type AnalysisRepository interface {
	Create(analysis *models.Analysis) error
	GetByID(id uint) (*models.Analysis, error)
	ListByUser(userID uint) ([]models.Analysis, error)
	UpdateStatus(id uint, status string) error
	UpdateCounts(id uint, tp int, fp int) error
	
}

type analysisRepository struct {
	db *gorm.DB
}

func NewAnalysisRepository(db *gorm.DB) AnalysisRepository {
	return &analysisRepository{db: db}
}

func (r *analysisRepository) Create(analysis *models.Analysis) error {
	if err := r.db.Create(analysis).Error; err != nil {
		logger.Log.Error().
			Str("repo", "analysis").
			Str("method", "Create").
			Err(err).
			Msg("failed to create analysis")
		return err
	}
	return nil
}

func (r *analysisRepository) GetByID(id uint) (*models.Analysis, error) {
	var analysis models.Analysis

	err := r.db.First(&analysis, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		logger.Log.Error().
			Str("repo", "analysis").
			Str("method", "GetByID").
			Uint("analysis_id", id).
			Err(err).
			Msg("failed to get analysis by id")

		return nil, err
	}

	return &analysis, nil
}

func (r *analysisRepository) ListByUser(userID uint) ([]models.Analysis, error) {
	var analyses []models.Analysis

	if err := r.db.
		Where("user_id = ?", userID).
		Order("uploaded_at DESC").
		Find(&analyses).
		Error; err != nil {

		logger.Log.Error().
			Str("repo", "analysis").
			Str("method", "ListByUser").
			Uint("user_id", userID).
			Err(err).
			Msg("failed to list analyses by user")

		return nil, err
	}

	return analyses, nil
}

func (r *analysisRepository) UpdateStatus(id uint, status string) error {
	res := r.db.
		Model(&models.Analysis{}).
		Where("id = ?", id).
		Update("status", status)

	if res.Error != nil {
		logger.Log.Error().
			Str("repo", "analysis").
			Str("method", "UpdateStatus").
			Uint("analysis_id", id).
			Err(res.Error).
			Msg("failed to update analysis status")
		return res.Error
	}

	if res.RowsAffected == 0 {
		logger.Log.Debug().
			Str("repo", "analysis").
			Str("method", "UpdateStatus").
			Uint("analysis_id", id).
			Msg("no analysis found to update status")
	}

	return nil
}

func (r *analysisRepository) UpdateCounts(
	id uint,
	tp int,
	fp int,
) error {
	res := r.db.
		Model(&models.Analysis{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"tp_count": tp,
			"fp_count": fp,
		})

	if res.Error != nil {
		logger.Log.Error().
			Str("repo", "analysis").
			Str("method", "UpdateCounts").
			Uint("analysis_id", id).
			Err(res.Error).
			Msg("failed to update analysis counts")
		return res.Error
	}

	if res.RowsAffected == 0 {
		logger.Log.Debug().
			Str("repo", "analysis").
			Str("method", "UpdateCounts").
			Uint("analysis_id", id).
			Msg("no analysis found to update counts")
	}

	return nil
}
