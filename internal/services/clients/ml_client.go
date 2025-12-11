package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"mws-ai/internal/models"
)

type MLClient interface {
	Predict(f models.Finding) (string, float64, error)
}

type mlHTTP struct {
	url string
}

func NewMLClient(baseURL string) MLClient {
	return &mlHTTP{url: baseURL}
}

func (c *mlHTTP) Predict(f models.Finding) (string, float64, error) {

	req := map[string]interface{}{
		"rule_id": f.RuleID,
		"snippet": f.Value,
	}

	body, _ := json.Marshal(req)

	resp, err := http.Post(c.url+"/predict", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", 0, fmt.Errorf("ml service returned %d", resp.StatusCode)
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
