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
	url    string
	client *http.Client
}

func NewLLMClient(url string) services.LLMClient {
	return &llmHTTP{
		url: url,
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

type llmFinding struct {
	ID       uint   `json:"id"`
	FilePath string `json:"file_path"`
	Line     int    `json:"line"`
	Value    string `json:"value"`
	RuleID   string `json:"rule_id"`
}

type llmRequest struct {
	Findings []llmFinding `json:"findings"`
}

type llmResponse struct {
	ID          uint    `json:"id"`
	Verdict     string  `json:"llm_verdict"`
	Confidence  float64 `json:"llm_confidence"`
	Explanation string  `json:"llm_explanation"`
}

type llmResponseEnvelope struct {
	Results []llmResponse `json:"results"`
}

// CLIENT
func (c *llmHTTP) AnalyzeBatch(
	findings []*models.Finding,
) (map[uint]services.LLMResult, error) {

	// build request
	req := llmRequest{
		Findings: make([]llmFinding, 0, len(findings)),
	}

	for _, f := range findings {
		req.Findings = append(req.Findings, llmFinding{
			ID:       f.ID,
			FilePath: f.FilePath,
			Line:     f.Line,
			Value:    f.Value,
			RuleID:   f.RuleID,
		})
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// send request
	resp, err := c.client.Post(
		c.url,
		"application/json",
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("llm request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("llm returned status %d", resp.StatusCode)
	}

	// decode response
	var env llmResponseEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, fmt.Errorf("llm decode error: %w", err)
	}

	// normalize
	out := make(map[uint]services.LLMResult, len(env.Results))
	for _, r := range env.Results {
		out[r.ID] = services.LLMResult{
			Verdict:     r.Verdict,
			Confidence:  r.Confidence,
			Explanation: r.Explanation,
		}
	}

	return out, nil
}
