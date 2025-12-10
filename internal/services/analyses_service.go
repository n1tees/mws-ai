package services

import (
	"time"

	"mws-ai/internal/models"
	"mws-ai/internal/repository"
)

type AnalysisService struct {
	repo repository.AnalysisRepository
}

func NewAnalysisService(repo repository.AnalysisRepository) *AnalysisService {
	return &AnalysisService{repo: repo}
}

// CreateAnalysis создаёт запись об анализе
func (s *AnalysisService) CreateAnalysis(userID uint, filePath string) (*models.Analysis, error) {
	analysis := &models.Analysis{
		UserID:     userID,
		UploadedAt: time.Now(),
	}

	err := s.repo.Create(analysis)
	if err != nil {
		return nil, err
	}

	return analysis, nil
}
