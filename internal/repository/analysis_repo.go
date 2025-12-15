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
	UpdateSummary(id uint, verdict string, tp int, fp int, conf *float64) error
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

	logger.Log.Debug().
		Str("repo", "analysis").
		Str("method", "Create").
		Uint("analysis_id", analysis.ID).
		Uint("user_id", analysis.UserID).
		Msg("analysis created")

	return nil
}

func (r *analysisRepository) GetByID(id uint) (*models.Analysis, error) {
	var analysis models.Analysis

	err := r.db.
		Preload("Findings").
		First(&analysis, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ❗ это НЕ ошибка — просто нет записи
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
	var list []models.Analysis

	if err := r.db.
		Where("user_id = ?", userID).
		Order("uploaded_at DESC").
		Find(&list).
		Error; err != nil {

		logger.Log.Error().
			Str("repo", "analysis").
			Str("method", "ListByUser").
			Uint("user_id", userID).
			Err(err).
			Msg("failed to list analyses by user")

		return nil, err
	}

	logger.Log.Debug().
		Str("repo", "analysis").
		Str("method", "ListByUser").
		Uint("user_id", userID).
		Int("count", len(list)).
		Msg("analyses listed")

	return list, nil
}

func (r *analysisRepository) UpdateStatus(id uint, status string) error {
	res := r.db.Model(&models.Analysis{}).
		Where("id = ?", id).
		Update("status", status)

	if res.Error != nil {
		logger.Log.Error().
			Str("repo", "analysis").
			Str("method", "UpdateStatus").
			Uint("analysis_id", id).
			Str("status", status).
			Err(res.Error).
			Msg("failed to update analysis status")

		return res.Error
	}

	if res.RowsAffected == 0 {
		// ❗ не ошибка, но полезно знать
		logger.Log.Debug().
			Str("repo", "analysis").
			Str("method", "UpdateStatus").
			Uint("analysis_id", id).
			Msg("no analysis found to update status")
	}

	return nil
}

func (r *analysisRepository) UpdateSummary(
	id uint,
	verdict string,
	tp int,
	fp int,
	conf *float64,
) error {
	res := r.db.Model(&models.Analysis{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"final_verdict":    verdict,
			"tp_count":         tp,
			"fp_count":         fp,
			"final_confidence": conf,
		})

	if res.Error != nil {
		logger.Log.Error().
			Str("repo", "analysis").
			Str("method", "UpdateSummary").
			Uint("analysis_id", id).
			Err(res.Error).
			Msg("failed to update analysis summary")

		return res.Error
	}

	if res.RowsAffected == 0 {
		logger.Log.Debug().
			Str("repo", "analysis").
			Str("method", "UpdateSummary").
			Uint("analysis_id", id).
			Msg("no analysis found to update summary")
	}

	return nil
}
