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

type AnalysisService struct {
	analysisRepo repository.AnalysisRepository
	findingRepo  repository.FindingRepository
	parser       SarifParser
	pipeline     PipelineExecutor
}

func NewAnalysisService(
	analysisRepo repository.AnalysisRepository,
	findingRepo repository.FindingRepository,
	parser SarifParser,
	pipeline PipelineExecutor,
) *AnalysisService {
	return &AnalysisService{
		analysisRepo: analysisRepo,
		findingRepo:  findingRepo,
		parser:       parser,
		pipeline:     pipeline,
	}
}

// =====================
// UPLOAD ENTRYPOINT
// =====================
func (s *AnalysisService) Upload(
	userID uint,
	filePath string,
) (*models.Analysis, error) {

	log := logger.Log.With().
		Str("service", "analysis").
		Str("method", "Upload").
		Uint("user_id", userID).
		Logger()

	analysis := &models.Analysis{
		UserID: userID,
		Status: "processing",
	}

	if err := s.analysisRepo.Create(analysis); err != nil {
		return nil, err
	}

	log.Info().
		Uint("analysis_id", analysis.ID).
		Msg("analysis created")

	go s.processAnalysis(analysis.ID, filePath)

	return analysis, nil
}

// =====================
// PROCESS ANALYSIS
// =====================
func (s *AnalysisService) processAnalysis(
	analysisID uint,
	filePath string,
) {

	log := logger.Log.With().
		Str("service", "analysis").
		Str("method", "processAnalysis").
		Uint("analysis_id", analysisID).
		Logger()

	start := time.Now()

	// ---------- PARSE SARIF ----------
	parsed, err := s.parser.Parse(filePath)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse sarif")
		_ = s.analysisRepo.UpdateStatus(analysisID, "failed")
		return
	}

	findings := make([]models.Finding, len(parsed))
	for i := range parsed {
		parsed[i].AnalysisID = analysisID
		findings[i] = parsed[i]
	}

	if err := s.findingRepo.BulkInsert(findings); err != nil {
		log.Error().Err(err).Msg("failed to insert findings")
		_ = s.analysisRepo.UpdateStatus(analysisID, "failed")
		return
	}

	log.Debug().
		Int("findings_count", len(findings)).
		Msg("findings inserted")

	// ---------- PIPELINE (IN-MEMORY) ----------
	ptrs := make([]*models.Finding, len(findings))
	for i := range findings {
		ptrs[i] = &findings[i]
	}

	if err := s.pipeline.Process(ptrs); err != nil {
		log.Error().Err(err).Msg("pipeline failed")
		_ = s.analysisRepo.UpdateStatus(analysisID, "failed")
		return
	}

	// ---------- SAVE FINDINGS ----------
	tpCount := 0
	fpCount := 0

	for _, f := range ptrs {
		f.Status = "processed"

		if f.FinalVerdict != nil {
			switch *f.FinalVerdict {
			case "TP":
				tpCount++
			case "FP":
				fpCount++
			}
		}

		_ = s.findingRepo.UpdateFields(f.ID, map[string]interface{}{
			"heuristic_triggered": f.HeuristicTriggered,
			"heuristic_reason":    f.HeuristicReason,
			"entropy_class":       f.EntropyClass,
			"entropy_value":       f.EntropyValue,
			"ml_verdict":          f.MlVerdict,
			"ml_confidence":       f.MlConfidence,
			"llm_verdict":         f.LlmVerdict,
			"llm_confidence":      f.LlmConfidence,
			"llm_explanation":     f.LlmExplanation,
			"final_verdict":       f.FinalVerdict,
			"status":              f.Status,
			"decision_source":     f.DecisionSource,
		})
	}

	// ---------- UPDATE ANALYSIS ----------
	_ = s.analysisRepo.UpdateCounts(analysisID, tpCount, fpCount)
	_ = s.analysisRepo.UpdateStatus(analysisID, "done")

	log.Info().
		Int("tp", tpCount).
		Int("fp", fpCount).
		Dur("duration", time.Since(start)).
		Msg("analysis completed")
}

func (s *AnalysisService) ListByUser(userID uint) ([]models.Analysis, error) {
	return s.analysisRepo.ListByUser(userID)
}

func (s *AnalysisService) GetByID(id uint) (*models.Analysis, error) {
	return s.analysisRepo.GetByID(id)
}

func (s *AnalysisService) GetDetails(
	id uint,
) (*models.Analysis, []models.Finding, error) {

	analysis, err := s.analysisRepo.GetByID(id)
	if err != nil || analysis == nil {
		return nil, nil, err
	}

	findings, err := s.findingRepo.ListByAnalysis(id)
	if err != nil {
		return nil, nil, err
	}

	return analysis, findings, nil
}
