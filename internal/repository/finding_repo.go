package repository

import (
	"mws-ai/internal/models"

	"gorm.io/gorm"
)

type FindingRepository interface {
	BulkInsert(analysisID uint, findings []models.Finding) error
	ListByAnalysis(analysisID uint) ([]models.Finding, error)
	Update(finding *models.Finding) error
}

type findingRepository struct {
	db *gorm.DB
}

func NewFindingRepository(db *gorm.DB) FindingRepository {
	return &findingRepository{db: db}
}

func (r *findingRepository) BulkInsert(analysisID uint, findings []models.Finding) error {
	if len(findings) == 0 {
		return nil
	}

	for i := range findings {
		findings[i].AnalysisID = analysisID
	}

	return r.db.Create(&findings).Error
}

func (r *findingRepository) ListByAnalysis(analysisID uint) ([]models.Finding, error) {
	var list []models.Finding

	err := r.db.
		Where("analysis_id = ?", analysisID).
		Order("id").
		Find(&list).Error

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r *findingRepository) Update(finding *models.Finding) error {
	return r.db.Save(finding).Error
}
