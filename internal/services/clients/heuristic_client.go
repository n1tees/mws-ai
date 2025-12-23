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

type heuristicHTTP struct {
	baseURL string
	client  *http.Client
}

func NewHeuristicClient(baseURL string) services.HeuristicClient {
	return &heuristicHTTP{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// DTO
type heuristicRequest struct {
	ID       uint   `json:"id"`
	FilePath string `json:"file_path"`
	Value    string `json:"value"`
}

type heuristicResponse struct {
	ID                 uint    `json:"id"`
	HeuristicTriggered bool    `json:"heuristic_triggered"`
	HeuristicReason    *string `json:"heuristic_reason"`
	Metrics            struct {
		Entropy      *float64 `json:"entropy"`
		Length       int      `json:"length"`
		EntropyClass *string  `json:"entropy_class"`
	} `json:"metrics"`
}

type heuristicResponseEnvelope struct {
	Results []heuristicResponse `json:"results"`
}

// CLIENT

func (h *heuristicHTTP) AnalyzeBatch(
	findings []*models.Finding,
) (map[uint]*services.HeuristicFacts, error) {

	// build request
	req := make([]heuristicRequest, 0, len(findings))
	for _, f := range findings {
		req = append(req, heuristicRequest{
			ID:       f.ID,
			FilePath: f.FilePath,
			Value:    f.Value,
		})
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := h.client.Post(
		h.baseURL,
		"application/json",
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("heuristic returned status %d", resp.StatusCode)
	}

	var env heuristicResponseEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, err
	}

	out := make(map[uint]*services.HeuristicFacts, len(env.Results))
	for _, r := range env.Results {

		out[r.ID] = &services.HeuristicFacts{
			HeuristicTriggered: r.HeuristicTriggered,
			HeuristicReason:    r.HeuristicReason,
			EntropyClass:       r.Metrics.EntropyClass,
			EntropyValue:       r.Metrics.Entropy,
		}
	}

	return out, nil
}
