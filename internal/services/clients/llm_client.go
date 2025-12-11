package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"mws-ai/internal/models"
)

type HeuristicClient interface {
	Evaluate(f models.Finding) (string, float64, error)
}

type heuristicHTTP struct {
	url string
}

func NewHeuristicClient(baseURL string) HeuristicClient {
	return &heuristicHTTP{url: baseURL}
}

func (c *heuristicHTTP) Evaluate(f models.Finding) (string, float64, error) {

	req := map[string]interface{}{
		"rule_id":   f.RuleID,
		"snippet":   f.Value,
		"file_path": f.FilePath,
		"severity":  f.Severity,
	}

	body, _ := json.Marshal(req)

	resp, err := http.Post(c.url+"/evaluate", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", 0, fmt.Errorf("heuristic service returned %d", resp.StatusCode)
	}

	var out struct {
		Verdict    string  `json:"verdict"`
		Confidence float64 `json:"confidence"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", 0, err
	}

	return out.Verdict, out.Confidence, nil
}
