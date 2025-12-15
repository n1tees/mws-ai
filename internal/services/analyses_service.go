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

	log := logger.Log.With().
		Str("service", "analysis").
		Str("method", "Upload").
		Uint("user_id", userID).
		Str("file_path", filePath).
		Logger()

	log.Debug().Msg("upload analysis request received")

	analysis := &models.Analysis{
		UserID:     userID,
		UploadedAt: time.Now(),
		Status:     "pending",
		FilePath:   filePath,
	}

	// Создаём Analysis
	if err := s.repo.Create(analysis); err != nil {
		log.Error().
			Err(err).
			Msg("failed to create analysis")

		return nil, err
	}

	log.Info().
		Uint("analysis_id", analysis.ID).
		Msg("analysis created, starting pipeline")

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

	start := time.Now()

	log := logger.Log.With().
		Str("service", "analysis").
		Str("method", "processAnalysis").
		Uint("analysis_id", analysisID).
		Str("file_path", filePath).
		Logger()

	log.Info().Msg("analysis processing started")

	_ = s.repo.UpdateStatus(analysisID, "processing")

	// Parse SARIF
	parsedFindings, err := s.sarifParser.Parse(filePath)
	if err != nil {
		log.Error().
			Err(err).
			Msg("SARIF parsing failed")

		_ = s.repo.UpdateStatus(analysisID, "failed")
		return
	}

	log.Debug().
		Int("findings_count", len(parsedFindings)).
		Msg("SARIF parsed successfully")

	// Convert slice of values → slice of pointers
	findings := make([]*models.Finding, len(parsedFindings))
	for i := range parsedFindings {
		findings[i] = &parsedFindings[i]
		findings[i].AnalysisID = analysisID
	}

	// Bulk insert
	if err := s.findingRepo.BulkInsert(analysisID, findings); err != nil {
		log.Error().
			Err(err).
			Msg("saving findings failed")

		_ = s.repo.UpdateStatus(analysisID, "failed")
		return
	}

	log.Debug().
		Int("findings_count", len(findings)).
		Msg("findings saved")

	// Summary accumulators
	tpCount := 0
	fpCount := 0
	var confSum float64
	var confCount int

	// Process pipeline
	for _, f := range findings {

		if err := s.pipeline.Process(f); err != nil {
			logger.Log.Error().
				Err(err).
				Uint("analysis_id", analysisID).
				Uint("finding_id", f.ID).
				Msg("pipeline processing failed")
		}

		if err := s.findingRepo.Update(f); err != nil {
			logger.Log.Error().
				Err(err).
				Uint("analysis_id", analysisID).
				Uint("finding_id", f.ID).
				Msg("failed to update finding")
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

	log.Debug().
		Int("tp_count", tpCount).
		Int("fp_count", fpCount).
		Int("confidence_count", confCount).
		Msg("pipeline processing finished")

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
		log.Error().
			Err(err).
			Msg("failed to update analysis summary")
	}

	// DONE
	_ = s.repo.UpdateStatus(analysisID, "done")

	log.Info().
		Str("final_verdict", finalVerdict).
		Int("tp", tpCount).
		Int("fp", fpCount).
		Dur("duration", time.Since(start)).
		Msg("analysis successfully processed")
}
