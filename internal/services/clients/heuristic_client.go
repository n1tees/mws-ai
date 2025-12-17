package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"mws-ai/internal/models"
	"mws-ai/internal/services"
)

type heuristicHTTP struct {
	baseURL string
	client  *http.Client
}

func NewHeuristicClient(baseURL string) services.HeuristicClient {
	return &heuristicHTTP{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type heuristicRequest struct {
	ID       uint   `json:"id"`
	FilePath string `json:"file_path"`
	Value    string `json:"value"`
	RuleID   string `json:"rule_id"`
	Severity string `json:"severity"`
}

type heuristicResponse struct {
	ID         uint     `json:"id"`
	Verdict    string   `json:"verdict"` // TP / FP
	Confidence float64  `json:"confidence"`
	Reasons    []string `json:"reasons"`
}

func (h *heuristicHTTP) AnalyzeBatch(
	findings []*models.Finding,
) (map[uint]services.HeuristicResult, error) {

	reqBody := make([]heuristicRequest, 0, len(findings))
	for _, f := range findings {
		reqBody = append(reqBody, heuristicRequest{
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

	resp, err := h.client.Post(
		fmt.Sprintf("%s/analyze", h.baseURL),
		"application/json",
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res []heuristicResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	out := make(map[uint]services.HeuristicResult, len(res))
	for _, r := range res {
		out[r.ID] = services.HeuristicResult{
			Verdict:    r.Verdict,
			Confidence: r.Confidence,
			Reasons:    r.Reasons,
		}
	}

	return out, nil
}
