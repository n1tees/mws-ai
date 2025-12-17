package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"mws-ai/internal/models"
	"mws-ai/internal/services"
)

type mlHTTP struct {
	baseURL string
	client  *http.Client
}

func NewMLClient(baseURL string) services.MLClient {
	return &mlHTTP{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type mlRequest struct {
	ID       uint   `json:"id"`
	FilePath string `json:"file_path"`
	Value    string `json:"value"`
	RuleID   string `json:"rule_id"`
	Severity string `json:"severity"`
}

type mlResponse struct {
	ID         uint    `json:"id"`
	Verdict    bool    `json:"verdict"` // true = TP
	Confidence float64 `json:"confidence"`
}

func (m *mlHTTP) PredictBatch(
	findings []*models.Finding,
) (map[uint]services.MLResult, error) {

	reqBody := make([]mlRequest, 0, len(findings))
	for _, f := range findings {
		reqBody = append(reqBody, mlRequest{
			ID:       f.ID,
			FilePath: f.FilePath,
			Value:    f.Value,
			RuleID:   f.RuleID,
			Severity: f.Severity,
		})
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Post(
		fmt.Sprintf("%s/predict", m.baseURL),
		"application/json",
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res []mlResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	out := make(map[uint]services.MLResult, len(res))
	for _, r := range res {
		out[r.ID] = services.MLResult{
			Verdict:    r.Verdict,
			Confidence: r.Confidence,
		}
	}

	return out, nil
}
