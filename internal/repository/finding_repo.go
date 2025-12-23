package repository

import (
	"mws-ai/internal/models"
	"mws-ai/pkg/logger"

	"gorm.io/gorm"
)

type FindingRepository interface {
	BulkInsert(findings []models.Finding) error
	UpdateFields(id uint, fields map[string]interface{}) error
	ListByAnalysis(analysisID uint) ([]models.Finding, error)
}

type findingRepository struct {
	db *gorm.DB
}

func NewFindingRepository(db *gorm.DB) FindingRepository {
	return &findingRepository{db: db}
}

func (r *findingRepository) BulkInsert(findings []models.Finding) error {
	if len(findings) == 0 {
		return nil
	}

	if err := r.db.Create(&findings).Error; err != nil {
		logger.Log.Error().
			Str("repo", "finding").
			Str("method", "BulkInsert").
			Err(err).
			Msg("failed to bulk insert findings")
		return err
	}

	return nil
}

func (r *findingRepository) UpdateFields(
	id uint,
	fields map[string]interface{},
) error {

	res := r.db.Model(&models.Finding{}).
		Where("id = ?", id).
		Updates(fields)

	if res.Error != nil {
		logger.Log.Error().
			Str("repo", "finding").
			Str("method", "UpdateFields").
			Uint("finding_id", id).
			Err(res.Error).
			Msg("failed to update finding fields")

		return res.Error
	}

	return nil
}

func (r *findingRepository) ListByAnalysis(analysisID uint) ([]models.Finding, error) {
	var findings []models.Finding

	if err := r.db.
		Where("analysis_id = ?", analysisID).
		Find(&findings).
		Error; err != nil {

		logger.Log.Error().
			Str("repo", "finding").
			Str("method", "ListByAnalysis").
			Uint("analysis_id", analysisID).
			Err(err).
			Msg("failed to list findings by analysis")

		return nil, err
	}

	return findings, nil
}
