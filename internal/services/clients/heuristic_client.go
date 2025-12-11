package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"mws-ai/internal/models"
)

type LLMClient interface {
	Analyze(f models.Finding) (string, string, error)
}

type llmHTTP struct {
	url string
}

func NewLLMClient(baseURL string) LLMClient {
	return &llmHTTP{url: baseURL}
}

func (c *llmHTTP) Analyze(f models.Finding) (string, string, error) {

	req := map[string]interface{}{
		"rule_id":   f.RuleID,
		"snippet":   f.Value,
		"message":   f.Value,
		"file_path": f.FilePath,
		"severity":  f.Severity,
	}

	body, _ := json.Marshal(req)

	resp, err := http.Post(c.url+"/analyze", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("llm service returned %d", resp.StatusCode)
	}

	var out struct {
		Verdict     string `json:"verdict"`
		Explanation string `json:"explanation"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", "", err
	}

	return out.Verdict, out.Explanation, nil
}
