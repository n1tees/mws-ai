package services

import (
	"fmt"

	"mws-ai/internal/models"
	"mws-ai/internal/repository"
	"mws-ai/internal/services/clients"
	"mws-ai/pkg/logger"
)

type PipelineExecutor interface {
	Process(findings []*models.Finding) error
}

type pipelineExecutor struct {
	heuristic   clients.HeuristicClient
	ml          clients.MLClient
	llm         clients.LLMClient
	findingRepo repository.FindingRepository
}

func NewPipeline(
	heuristic clients.HeuristicClient,
	ml clients.MLClient,
	llm clients.LLMClient,
	findingRepo repository.FindingRepository,
) PipelineExecutor {
	return &pipelineExecutor{
		heuristic:   heuristic,
		ml:          ml,
		llm:         llm,
		findingRepo: findingRepo,
	}
}

//
// ========== MAIN PIPELINE ==========
//

func (p *pipelineExecutor) Process(findings []*models.Finding) error {

	// ===============================
	// 1. HEURISTIC
	// ===============================
	heuristicResults, err := p.heuristic.Analyze(findings)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Heuristic batch failed")
		return fmt.Errorf("heuristic failed: %w", err)
	}

	var toML []*models.Finding

	for _, f := range findings {
		res, ok := heuristicResults[f.ID]
		if !ok {
			logger.Log.Warn().Uint("finding_id", f.ID).Msg("Missing heuristic result")
			continue
		}

		f.RuleVerdict = &res.Verdict
		f.RuleConfidence = &res.Confidence
		p.saveFinding(f)

		// STOP RULE: strong TP
		if res.Verdict == "TP" && res.Confidence >= 0.8 {
			// финальный вывод
			f.FinalVerdict = &res.Verdict
			f.FinalConfidence = &res.Confidence
			p.saveFinding(f)

			logger.Log.Info().Uint("id", f.ID).Msg("Heuristic says strong TP → STOP")
			continue
		}

		toML = append(toML, f)
	}

	if len(toML) == 0 {
		return nil
	}

	// ===============================
	// 2. ML
	// ===============================
	mlResults, err := p.ml.Predict(toML)
	if err != nil {
		logger.Log.Error().Err(err).Msg("ML batch failed")
		return fmt.Errorf("ml failed: %w", err)
	}

	var toLLM []*models.Finding

	for _, f := range toML {
		res, ok := mlResults[f.ID]
		if !ok {
			logger.Log.Warn().Uint("finding_id", f.ID).Msg("Missing ML result")
			continue
		}

		verdict := "FP"
		if res.Verdict {
			verdict = "TP"
		}

		f.MlVerdict = &verdict
		f.MlConfidence = &res.Confidence
		p.saveFinding(f)

		// STOP RULE: ML confident TP
		if res.Verdict && res.Confidence >= 0.9 {
			f.FinalVerdict = &verdict
			f.FinalConfidence = &res.Confidence
			p.saveFinding(f)

			logger.Log.Info().Uint("id", f.ID).Msg("ML says strong TP → STOP")
			continue
		}

		toLLM = append(toLLM, f)
	}

	if len(toLLM) == 0 {
		return nil
	}

	// ===============================
	// 3. LLM
	// ===============================
	llmResults, err := p.llm.AnalyzeBatch(toLLM)
	if err != nil {
		logger.Log.Error().Err(err).Msg("LLM batch failed")
		return fmt.Errorf("llm failed: %w", err)
	}

	for _, f := range toLLM {
		res, ok := llmResults[f.ID]
		if !ok {
			logger.Log.Warn().Uint("finding_id", f.ID).Msg("Missing LLM result")
			continue
		}

		f.LlmVerdict = &res.Verdict
		f.LlmConfidence = &res.Confidence
		f.LlmExplanation = &res.Explanation

		f.FinalVerdict = &res.Verdict
		f.FinalConfidence = &res.Confidence

		p.saveFinding(f)
	}

	return nil
}

//
// ========== HELPERS ==========
//

func (p *pipelineExecutor) saveFinding(f *models.Finding) {
	if err := p.findingRepo.Update(f); err != nil {
		logger.Log.Error().Err(err).
			Uint("finding_id", f.ID).
			Msg("Failed to update finding in DB")
	}
}
