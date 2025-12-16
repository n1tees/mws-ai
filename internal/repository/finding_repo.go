package repository

import (
	"mws-ai/internal/models"
	"mws-ai/pkg/logger"

	"gorm.io/gorm"
)

type FindingRepository interface {
	BulkInsert(analysisID uint, findings []*models.Finding) error
	ListByAnalysis(analysisID uint) ([]*models.Finding, error)
	Update(f *models.Finding) error
	UpdateStatus(id uint, status string) error
	UpdateFields(id uint, fields map[string]interface{}) error
}

type findingRepository struct {
	db *gorm.DB
}

func NewFindingRepository(db *gorm.DB) FindingRepository {
	return &findingRepository{db: db}
}

func (r *findingRepository) BulkInsert(analysisID uint, findings []*models.Finding) error {
	for _, f := range findings {
		f.AnalysisID = analysisID
	}

	if len(findings) == 0 {
		logger.Log.Debug().
			Str("repo", "finding").
			Str("method", "BulkInsert").
			Uint("analysis_id", analysisID).
			Msg("no findings to insert")

		return nil
	}

	if err := r.db.Create(&findings).Error; err != nil {
		logger.Log.Error().
			Str("repo", "finding").
			Str("method", "BulkInsert").
			Uint("analysis_id", analysisID).
			Int("count", len(findings)).
			Err(err).
			Msg("failed to bulk insert findings")

		return err
	}

	logger.Log.Debug().
		Str("repo", "finding").
		Str("method", "BulkInsert").
		Uint("analysis_id", analysisID).
		Int("count", len(findings)).
		Msg("findings inserted")

	return nil
}

func (r *findingRepository) ListByAnalysis(analysisID uint) ([]*models.Finding, error) {
	var list []*models.Finding

	if err := r.db.
		Where("analysis_id = ?", analysisID).
		Order("id").
		Find(&list).
		Error; err != nil {

		logger.Log.Error().
			Str("repo", "finding").
			Str("method", "ListByAnalysis").
			Uint("analysis_id", analysisID).
			Err(err).
			Msg("failed to list findings by analysis")

		return nil, err
	}

	logger.Log.Debug().
		Str("repo", "finding").
		Str("method", "ListByAnalysis").
		Uint("analysis_id", analysisID).
		Int("count", len(list)).
		Msg("findings listed")

	return list, nil
}

func (r *findingRepository) Update(f *models.Finding) error {
	res := r.db.Model(&models.Finding{}).
		Where("id = ?", f.ID).
		Updates(f)

	if res.Error != nil {
		logger.Log.Error().
			Str("repo", "finding").
			Str("method", "Update").
			Uint("finding_id", f.ID).
			Err(res.Error).
			Msg("failed to update finding")

		return res.Error
	}

	if res.RowsAffected == 0 {
		logger.Log.Debug().
			Str("repo", "finding").
			Str("method", "Update").
			Uint("finding_id", f.ID).
			Msg("no finding found to update")
	}

	return nil
}

func (r *findingRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&models.Finding{}).
		Where("id = ?", id).
		Updates(fields).
		Error
}

func (r *findingRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Finding{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}
