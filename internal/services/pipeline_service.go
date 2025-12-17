package services

import (
	"errors"

	"mws-ai/internal/models"
	"mws-ai/pkg/logger"
)

// Interfaces clients
type HeuristicClient interface {
	AnalyzeBatch(findings []*models.Finding) (map[uint]HeuristicResult, error)
}

type MLClient interface {
	PredictBatch(findings []*models.Finding) (map[uint]MLResult, error)
}

type LLMClient interface {
	AnalyzeBatch(findings []*models.Finding) (map[uint]LLMResult, error)
}

// Results clients
type HeuristicResult struct {
	Verdict    string
	Confidence float64
	Reasons    []string
}

type MLResult struct {
	Verdict    bool // true = TP, false = FP
	Confidence float64
}

type LLMResult struct {
	Verdict     string // TP | FP (или текст, который мы нормализуем)
	Confidence  float64
	Explanation string
}

// PipelineExecutor
type PipelineExecutor interface {
	Process(findings []*models.Finding) error
}

// Realisation
type pipelineExecutor struct {
	heuristic HeuristicClient
	ml        MLClient
	llm       LLMClient
}

func NewPipelineExecutor(
	heuristic HeuristicClient,
	ml MLClient,
	llm LLMClient,
) PipelineExecutor {
	return &pipelineExecutor{
		heuristic: heuristic,
		ml:        ml,
		llm:       llm,
	}
}

// MAIN PIPELINE ENTRYPOINT
func (p *pipelineExecutor) Process(findings []*models.Finding) error {
	log := logger.Log.With().
		Str("service", "pipeline").
		Logger()

	log.Info().Msg("pipeline started")

	if len(findings) == 0 {
		return nil
	}

	// 1. HEURISTIC
	heuristicResults, err := p.heuristic.AnalyzeBatch(findings)
	if err != nil {
		return err
	}

	toML := make([]*models.Finding, 0)
	toLLM := make([]*models.Finding, 0)

	for _, f := range findings {
		if res, ok := heuristicResults[f.ID]; ok {
			f.HeuristicVerdict = &res.Verdict
			f.HeuristicConfidence = &res.Confidence
			f.HeuristicReasons = res.Reasons
		}

		// stop-rule heuristic
		if f.HeuristicVerdict != nil &&
			*f.HeuristicVerdict == "TP" &&
			f.HeuristicConfidence != nil &&
			*f.HeuristicConfidence >= 0.8 {

			final := "TP"
			conf := *f.HeuristicConfidence
			f.FinalVerdict = &final
			f.FinalConfidence = &conf
			continue
		}

		toML = append(toML, f)
	}

	// 2. ML
	if len(toML) > 0 {
		mlResults, err := p.ml.PredictBatch(toML)
		if err != nil {
			return err
		}

		for _, f := range toML {
			if res, ok := mlResults[f.ID]; ok {
				verdict := "FP"
				if res.Verdict {
					verdict = "TP"
				}
				f.MlVerdict = &verdict
				f.MlConfidence = &res.Confidence
			}

			// stop-rule ML
			if f.MlVerdict != nil &&
				*f.MlVerdict == "TP" &&
				f.MlConfidence != nil &&
				*f.MlConfidence >= 0.9 {

				final := "TP"
				conf := *f.MlConfidence
				f.FinalVerdict = &final
				f.FinalConfidence = &conf
				continue
			}

			toLLM = append(toLLM, f)
		}
	}

	// 3. LLM
	if len(toLLM) > 0 {
		llmResults, err := p.llm.AnalyzeBatch(toLLM)
		if err != nil {
			return err
		}

		for _, f := range toLLM {
			res, ok := llmResults[f.ID]
			if !ok {
				return errors.New("missing LLM result")
			}

			f.LlmVerdict = &res.Verdict
			f.LlmConfidence = &res.Confidence
			f.LlmExplanation = &res.Explanation

			final := "FP"
			if res.Verdict == "TP" {
				final = "TP"
			}

			conf := res.Confidence
			f.FinalVerdict = &final
			f.FinalConfidence = &conf
		}
	}

	log.Info().Msg("pipeline finished")
	return nil
}
