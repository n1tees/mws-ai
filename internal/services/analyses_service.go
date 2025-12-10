package services

import (
	"time"

	"mws-ai/internal/models"
	"mws-ai/internal/repository"
	"mws-ai/pkg/logger"
)

type SarifParser interface {
	Parse(filePath string) ([]models.Finding, error)
}

type PipelineExecutor interface {
	Run(analysisID uint) error
}

type AnalysisService struct {
	repo        repository.AnalysisRepository
	findingRepo repository.FindingRepository
	sarifParser SarifParser
	pipeline    PipelineExecutor
}

func NewAnalysisService(
	repo repository.AnalysisRepository,
	findingRepo repository.FindingRepository,
	parser SarifParser,
	pipeline PipelineExecutor,
) *AnalysisService {
	return &AnalysisService{
		repo:        repo,
		findingRepo: findingRepo,
		sarifParser: parser,
		pipeline:    pipeline,
	}
}

// Upload: создаёт анализ и запускает pipeline
func (s *AnalysisService) Upload(userID uint, filePath string) (*models.Analysis, error) {

	analysis := &models.Analysis{
		UserID:     userID,
		UploadedAt: time.Now(),
		Status:     "pending",
		FilePath:   filePath,
	}

	// 1. Сохраняем analysis в БД
	if err := s.repo.Create(analysis); err != nil {
		return nil, err
	}

	logger.Log.Info().
		Uint("analysis_id", analysis.ID).
		Msg("Analysis created, starting pipeline...")

	// 2. Запускаем асинхронный pipeline
	go s.processAnalysis(analysis.ID, filePath)

	return analysis, nil
}

func (s *AnalysisService) GetByID(id uint) (*models.Analysis, error) {
	return s.repo.GetByID(id)
}

func (s *AnalysisService) ListByUser(userID uint) ([]models.Analysis, error) {
	return s.repo.ListByUser(userID)
}

// processAnalysis запускается в фоне
func (s *AnalysisService) processAnalysis(analysisID uint, filePath string) {

	// Обновляем статус на "processing"
	_ = s.repo.UpdateStatus(analysisID, "processing")

	// 1. Парсим SARIF
	findings, err := s.sarifParser.Parse(filePath)
	if err != nil {
		logger.Log.Error().Err(err).
			Uint("analysis_id", analysisID).
			Msg("SARIF parsing failed")

		_ = s.repo.UpdateStatus(analysisID, "failed")
		return
	}

	// 2. Сохраняем findings
	if err := s.findingRepo.BulkInsert(analysisID, findings); err != nil {
		logger.Log.Error().Err(err).
			Uint("analysis_id", analysisID).
			Msg("Saving findings failed")

		_ = s.repo.UpdateStatus(analysisID, "failed")
		return
	}

	// 3. Запустить pipeline (heuristic + ml + llm)
	if err := s.pipeline.Run(analysisID); err != nil {
		logger.Log.Error().Err(err).
			Uint("analysis_id", analysisID).
			Msg("Pipeline failed")

		_ = s.repo.UpdateStatus(analysisID, "failed")
		return
	}

	// 4. Обновляем статус на "done"
	_ = s.repo.UpdateStatus(analysisID, "done")

	logger.Log.Info().
		Uint("analysis_id", analysisID).
		Msg("Analysis successfully processed")
}
