package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mws-ai/internal/models"
	"mws-ai/internal/services"
)

type llmHTTP struct {
	baseURL string
	client  *http.Client
}

func NewLLMClient(baseURL string) services.LLMClient {
	return &llmHTTP{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// ВАЖНО: метод реализует services.LLMClient
func (c *llmHTTP) AnalyzeBatch(
	findings []*models.Finding,
) (map[uint]services.LLMResult, error) {

	// === Формируем payload ===
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

	resp, err := c.client.Post(
		c.baseURL,
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("llm request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LLM returned %d", resp.StatusCode)
	}

	// === Декод ответа ===
	var out struct {
		Results []struct {
			ID          uint    `json:"id"`
			Verdict     string  `json:"llm_verdict"`
			Confidence  float64 `json:"llm_confidence"`
			Explanation string  `json:"llm_explanation"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode LLM error: %w", err)
	}

	results := make(map[uint]services.LLMResult, len(out.Results))
	for _, r := range out.Results {
		results[r.ID] = services.LLMResult{
			Verdict:     r.Verdict,
			Confidence:  r.Confidence,
			Explanation: r.Explanation,
		}
	}

	return results, nil
}
