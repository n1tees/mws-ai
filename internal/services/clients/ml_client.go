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

type mlHTTP struct {
	baseURL string
	client  *http.Client
}

func NewMLClient(baseURL string) services.MLClient {
	return &mlHTTP{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

type mlResponse struct {
	ID           uint    `json:"id"`
	MLPredict    bool    `json:"MLPredict"`
	MLConfidence float64 `json:"MLConfidence"`
}

type mlEnvelope struct {
	Results []mlResponse `json:"results"`
}

func (m *mlHTTP) PredictBatch(
	findings []*models.Finding,
) (map[uint]services.MLResult, error) {

	payload, err := json.Marshal(findings)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Post(
		m.baseURL,
		"application/json",
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ml returned %d", resp.StatusCode)
	}

	var env mlEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, err
	}

	out := make(map[uint]services.MLResult, len(env.Results))
	for _, r := range env.Results {
		out[r.ID] = services.MLResult{
			Verdict:    r.MLPredict,
			Confidence: r.MLConfidence,
		}
	}

	return out, nil
}
