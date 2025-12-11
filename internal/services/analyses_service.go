package services

import (
	"time"

	"mws-ai/internal/models"
	"mws-ai/internal/repository"
	"mws-ai/pkg/logger"
)

// Interfaces
type SarifParser interface {
	Parse(filePath string) ([]models.Finding, error)
}

// AnalysisService
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

// Public API: Upload
func (s *AnalysisService) Upload(userID uint, filePath string) (*models.Analysis, error) {

	analysis := &models.Analysis{
		UserID:     userID,
		UploadedAt: time.Now(),
		Status:     "pending",
		FilePath:   filePath,
	}

	// Создаём Analysis
	if err := s.repo.Create(analysis); err != nil {
		return nil, err
	}

	logger.Log.Info().
		Uint("analysis_id", analysis.ID).
		Msg("Analysis created, starting pipeline...")

	// Асинхронный запуск
	go s.processAnalysis(analysis.ID, filePath)

	return analysis, nil
}

// Public API: Queries
func (s *AnalysisService) GetByID(id uint) (*models.Analysis, error) {
	return s.repo.GetByID(id)
}

func (s *AnalysisService) ListByUser(userID uint) ([]models.Analysis, error) {
	return s.repo.ListByUser(userID)
}

// INTERNAL: Background worker
func (s *AnalysisService) processAnalysis(analysisID uint, filePath string) {

	_ = s.repo.UpdateStatus(analysisID, "processing")

	// Parse SARIF
	parsedFindings, err := s.sarifParser.Parse(filePath)
	if err != nil {
		logger.Log.Error().Err(err).Uint("analysis_id", analysisID).Msg("SARIF parsing failed")
		_ = s.repo.UpdateStatus(analysisID, "failed")
		return
	}

	// Convert slice of values → slice of pointers
	findings := make([]*models.Finding, len(parsedFindings))
	for i := range parsedFindings {
		findings[i] = &parsedFindings[i]
		findings[i].AnalysisID = analysisID
	}

	// Bulk insert
	if err := s.findingRepo.BulkInsert(analysisID, findings); err != nil {
		logger.Log.Error().Err(err).Uint("analysis_id", analysisID).Msg("Saving findings failed")
		_ = s.repo.UpdateStatus(analysisID, "failed")
		return
	}

	// Summary accumulators
	tpCount := 0
	fpCount := 0
	var confSum float64
	var confCount int

	// Process pipeline
	for _, f := range findings {

		if err := s.pipeline.Process(f); err != nil {
			logger.Log.Error().Err(err).Uint("finding_id", f.ID).Msg("Pipeline failed")
		}

		if err := s.findingRepo.Update(f); err != nil {
			logger.Log.Error().Err(err).Uint("finding_id", f.ID).Msg("Failed to update finding")
		}

		if f.FinalVerdict != nil {
			if *f.FinalVerdict == "TP" {
				tpCount++
			} else if *f.FinalVerdict == "FP" {
				fpCount++
			}
		}

		if f.FinalConfidence != nil {
			confSum += *f.FinalConfidence
			confCount++
		}
	}

	// COMPUTE FINAL SUMMARY
	finalVerdict := "FP"
	if tpCount > fpCount {
		finalVerdict = "TP"
	}

	var avgConf *float64
	if confCount > 0 {
		c := confSum / float64(confCount)
		avgConf = &c
	}

	// Save summary
	if err := s.repo.UpdateSummary(analysisID, finalVerdict, tpCount, fpCount, avgConf); err != nil {
		logger.Log.Error().Err(err).Uint("analysis_id", analysisID).Msg("Failed to update summary")
	}

	// DONE
	_ = s.repo.UpdateStatus(analysisID, "done")

	logger.Log.Info().Uint("analysis_id", analysisID).Msg("Analysis successfully processed")
}
