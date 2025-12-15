package services

import (
	"mws-ai/internal/models"
	"mws-ai/pkg/logger"
)

// Интерфейсы для вызова внешних сервисов
type HeuristicClient interface {
	Evaluate(f models.Finding) (verdict string, confidence float64, err error)
}

type MLClient interface {
	Predict(f models.Finding) (verdict string, confidence float64, err error)
}

type LLMClient interface {
	Analyze(f models.Finding) (verdict string, explanation string, err error)
}

type PipelineExecutor interface {
	Process(f *models.Finding) error
}

type pipelineExecutor struct {
	heuristic HeuristicClient
	ml        MLClient
	llm       LLMClient
}

func NewPipeline(heuristic HeuristicClient, ml MLClient, llm LLMClient) PipelineExecutor {
	return &pipelineExecutor{
		heuristic: heuristic,
		ml:        ml,
		llm:       llm,
	}
}

// Основная логика обработки одного finding
func (p *pipelineExecutor) Process(f *models.Finding) error {

	log := logger.Log.With().
		Str("service", "pipeline").
		Str("method", "Process").
		Uint("finding_id", f.ID).
		Logger()

	log.Debug().Msg("pipeline processing started")

	//---------------------------------------------------
	// 1. HEURISTIC (жёсткий фильтр)
	//---------------------------------------------------
	rVerdict, rConf, err := p.heuristic.Evaluate(*f)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("heuristic evaluation failed, continuing to ML")
	} else {
		f.RuleVerdict = &rVerdict
		f.RuleConfidence = &rConf

		log.Debug().
			Str("verdict", rVerdict).
			Float64("confidence", rConf).
			Msg("heuristic verdict received")

		// Если эвристика уверена (FP или TP) → завершаем анализ
		if rVerdict == "FP" || rVerdict == "TP" {
			f.FinalVerdict = &rVerdict
			f.FinalConfidence = &rConf

			log.Debug().
				Str("final_verdict", rVerdict).
				Float64("final_confidence", rConf).
				Msg("final verdict decided by heuristic")

			return nil
		}
	}

	//---------------------------------------------------
	// 2. ML (основной модуль)
	//---------------------------------------------------
	mlVerdict, mlConf, err := p.ml.Predict(*f)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("ML prediction failed, fallback to default FP")

		// Если ML упал и эвристика не дала verdict → FP по умолчанию
		defaultVerdict := "FP"
		defaultConf := 0.5

		f.FinalVerdict = &defaultVerdict
		f.FinalConfidence = &defaultConf

		log.Debug().
			Str("final_verdict", defaultVerdict).
			Float64("final_confidence", defaultConf).
			Msg("final verdict set by fallback")

		return nil
	}

	f.MlVerdict = &mlVerdict
	f.MlConfidence = &mlConf

	log.Debug().
		Str("verdict", mlVerdict).
		Float64("confidence", mlConf).
		Msg("ML verdict received")

	//---------------------------------------------------
	// Если ML дал уверенность >= 0.7 → это итог
	//---------------------------------------------------
	if mlConf >= 0.7 {
		f.FinalVerdict = &mlVerdict
		f.FinalConfidence = &mlConf

		log.Debug().
			Str("final_verdict", mlVerdict).
			Float64("final_confidence", mlConf).
			Msg("final verdict decided by ML")

		return nil
	}

	//---------------------------------------------------
	// 3. ML < 0.7 → LLM
	//---------------------------------------------------
	llmVerdict, explanation, err := p.llm.Analyze(*f)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("LLM analysis failed, using ML verdict")

		f.FinalVerdict = &mlVerdict
		f.FinalConfidence = &mlConf

		log.Debug().
			Str("final_verdict", mlVerdict).
			Float64("final_confidence", mlConf).
			Msg("final verdict taken from ML after LLM failure")

		return nil
	}

	f.LlmVerdict = &llmVerdict
	f.LlmExplanation = &explanation

	log.Debug().
		Str("verdict", llmVerdict).
		Msg("LLM verdict received")

	//---------------------------------------------------
	// Final verdict = LLM verdict
	//---------------------------------------------------
	f.FinalVerdict = &llmVerdict

	// Confidence = усреднение эвристики и ML
	conf := computeFinalConfidence(f)
	f.FinalConfidence = &conf

	log.Debug().
		Str("final_verdict", llmVerdict).
		Float64("final_confidence", conf).
		Msg("final verdict decided by LLM")

	return nil
}

func computeFinalConfidence(f *models.Finding) float64 {
	sum := 0.0
	count := 0

	if f.RuleConfidence != nil {
		sum += *f.RuleConfidence
		count++
	}
	if f.MlConfidence != nil {
		sum += *f.MlConfidence
		count++
	}

	if count == 0 {
		return 0.0
	}
	return sum / float64(count)
}

// ____________________________ ЗАГЛУШКА
type DummyPipeline struct{}

func NewDummyPipeline() PipelineExecutor {
	return &DummyPipeline{}
}

func (p *DummyPipeline) Process(f *models.Finding) error {
	log := logger.Log.With().
		Str("service", "pipeline").
		Str("method", "DummyProcess").
		Uint("finding_id", f.ID).
		Logger()

	verdict := "FP"
	conf := 0.5

	f.FinalVerdict = &verdict
	f.FinalConfidence = &conf

	log.Debug().
		Str("final_verdict", verdict).
		Float64("final_confidence", conf).
		Msg("dummy pipeline applied")

	return nil
}
