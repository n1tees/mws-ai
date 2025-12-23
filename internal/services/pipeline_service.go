package services

import (
	"errors"
	"mws-ai/internal/models"
	"mws-ai/pkg/logger"
)

// Interfaces clients
type HeuristicClient interface {
	AnalyzeBatch(findings []*models.Finding) (map[uint]*HeuristicFacts, error)
}
type MLClient interface {
	PredictBatch(findings []*models.Finding) (map[uint]MLResult, error)
}

type LLMClient interface {
	AnalyzeBatch(findings []*models.Finding) (map[uint]LLMResult, error)
}

// Results clients

type HeuristicFacts struct {
	HeuristicTriggered bool
	HeuristicReason    *string
	EntropyClass       *string
	EntropyValue       *float64
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

	// 1. HEURISTIC (NEW)
	heuristicResults, err := p.heuristic.AnalyzeBatch(findings)
	if err != nil {
		return err
	}

	toML := make([]*models.Finding, 0)

	for _, f := range findings {
		h, ok := heuristicResults[f.ID]
		if !ok {
			continue
		}

		// сохраняем факты эвристики
		f.HeuristicTriggered = h.HeuristicTriggered
		f.HeuristicReason = h.HeuristicReason
		f.EntropyClass = h.EntropyClass
		f.EntropyValue = h.EntropyValue

		// эвристика сработала
		if f.HeuristicTriggered {
			final := "FP"
			f.FinalVerdict = &final
			f.DecisionSource = "heuristic (static)"
			continue
		}
		//  недостаточная энтропия
		if f.EntropyClass != nil && *f.EntropyClass != "acceptable" {
			final := "FP"
			f.FinalVerdict = &final
			f.DecisionSource = "heuristic (entropy)"
			continue
		}

		// идём дальше в ML
		toML = append(toML, f)
	}

	// 2. ML
	toLLM := make([]*models.Finding, 0)

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

			// STOP-RULE ML
			if f.MlVerdict != nil &&
				*f.MlVerdict == "TP" &&
				f.MlConfidence != nil &&
				*f.MlConfidence >= 0.8 {

				final := "TP"
				f.FinalVerdict = &final
				f.DecisionSource = "ml"
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
			if res.Verdict == "real_secret" {
				final = "TP"
			} else if res.Verdict == "unknown" {
				final = "LLM_UNKNOWN_ANSWER"
			}

			f.FinalVerdict = &final
			f.DecisionSource = "llm"
		}
	}

	log.Info().Msg("pipeline finished")
	return nil
}
