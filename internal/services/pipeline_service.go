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

	//---------------------------------------------------
	// 1. HEURISTIC (жёсткий фильтр)
	//---------------------------------------------------
	rVerdict, rConf, err := p.heuristic.Evaluate(*f)
	if err != nil {
		logger.Log.Warn().Err(err).Msg("Heuristic failed — continuing to ML")
	} else {
		f.RuleVerdict = &rVerdict
		f.RuleConfidence = &rConf

		// Если эвристика уверена (FP или TP) → завершаем анализ
		if rVerdict == "FP" || rVerdict == "TP" {
			f.FinalVerdict = &rVerdict
			f.FinalConfidence = &rConf
			return nil
		}
	}

	//---------------------------------------------------
	// 2. ML (основной модуль)
	//---------------------------------------------------
	mlVerdict, mlConf, err := p.ml.Predict(*f)
	if err != nil {
		logger.Log.Warn().Err(err).Msg("ML failed — fallback to heuristic")

		// Если ML упал и эвристика не дала verdict → FP по умолчанию
		defaultVerdict := "FP"
		defaultConf := 0.5

		f.FinalVerdict = &defaultVerdict
		f.FinalConfidence = &defaultConf
		return nil
	}

	f.MlVerdict = &mlVerdict
	f.MlConfidence = &mlConf

	//---------------------------------------------------
	// Если ML дал уверенность >= 0.7 → это итог
	//---------------------------------------------------
	if mlConf >= 0.7 {
		f.FinalVerdict = &mlVerdict
		f.FinalConfidence = &mlConf
		return nil
	}

	//---------------------------------------------------
	// 3. ML < 0.7 → LLM
	//---------------------------------------------------
	llmVerdict, explanation, err := p.llm.Analyze(*f)
	if err != nil {
		logger.Log.Warn().Err(err).Msg("LLM failed — using ML verdict")
		f.FinalVerdict = &mlVerdict
		f.FinalConfidence = &mlConf
		return nil
	}

	f.LlmVerdict = &llmVerdict
	f.LlmExplanation = &explanation

	//---------------------------------------------------
	// Final verdict = LLM verdict
	//---------------------------------------------------
	f.FinalVerdict = &llmVerdict

	// Confidence = усреднение эвристики и ML
	conf := computeFinalConfidence(f)
	f.FinalConfidence = &conf

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

	// LLM confidence у нас нет → считаем как 1.0 или пропускаем
	// Лучше пропускать, иначе будет ложное усиление
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
	verdict := "FP"
	conf := 0.5

	f.FinalVerdict = &verdict
	f.FinalConfidence = &conf

	return nil
}
