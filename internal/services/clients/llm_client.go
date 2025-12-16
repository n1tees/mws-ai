package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mws-ai/internal/models"
)

type LLMResult struct {
	ID          uint    `json:"id"`
	Verdict     string  `json:"llm_verdict"`
	Confidence  float64 `json:"llm_confidence"`
	Explanation string  `json:"llm_explanation"`
}

type LLMClient interface {
	AnalyzeBatch(findings []*models.Finding) (map[uint]LLMResult, error)
}

type llmHTTP struct {
	baseURL string
	client  *http.Client
}

func NewLLMClient(baseURL string) LLMClient {
	return &llmHTTP{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *llmHTTP) AnalyzeBatch(findings []*models.Finding) (map[uint]LLMResult, error) {

	// === Формируем JSON как требует LLM ===
	req := struct {
		Findings []map[string]interface{} `json:"findings"`
	}{
		Findings: make([]map[string]interface{}, len(findings)),
	}

	for i, f := range findings {
		req.Findings[i] = map[string]interface{}{
			"id":        f.ID,
			"file_path": f.FilePath,
			"line":      f.Line,
			"value":     f.Value,
			"rule_id":   f.RuleID,
		}
	}

	body, _ := json.Marshal(req)

	// === POST ===
	resp, err := c.client.Post(c.baseURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("llm request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("LLM returned %d", resp.StatusCode)
	}

	// === Декод ===
	var out struct {
		Results []LLMResult `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode LLM error: %w", err)
	}

	// === Преобразуем в map[id]LLMResult ===
	results := make(map[uint]LLMResult)
	for _, r := range out.Results {
		results[r.ID] = r
	}

	return results, nil
}
