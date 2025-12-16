package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mws-ai/internal/models"
)

type MLResult struct {
	ID         uint
	Verdict    bool
	Confidence float64
}

type MLClient interface {
	Predict(findings []*models.Finding) (map[uint]MLResult, error)
}

type mlHTTP struct {
	baseURL string
	client  *http.Client
}

func NewMLClient(baseURL string) MLClient {
	return &mlHTTP{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 25 * time.Second},
	}
}

func (m *mlHTTP) Predict(findings []*models.Finding) (map[uint]MLResult, error) {

	// === 1. Формируем JSON-массив ===
	payload := make([]map[string]interface{}, len(findings))

	for i, f := range findings {
		payload[i] = map[string]interface{}{
			"id":                 f.ID,
			"rule_id":            f.RuleID,
			"file_path":          f.FilePath,
			"line":               f.Line,
			"value":              f.Value,
			"severity":           f.Severity,
			"scanner_confidence": f.ScannerConfidence,
		}
	}

	body, _ := json.Marshal(payload)

	// === 2. POST на ML-контейнер ===
	resp, err := m.client.Post(m.baseURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ml request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ml service returned %d", resp.StatusCode)
	}

	// === 3. Читаем ответ ===
	var out struct {
		Results []struct {
			ID         uint    `json:"id"`
			Predict    bool    `json:"MLPredict"`
			Confidence float64 `json:"MLConfidence"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode ml response: %w", err)
	}

	// === 4. Преобразуем ответ в map[id]MLResult ===
	result := make(map[uint]MLResult)

	for _, r := range out.Results {
		result[r.ID] = MLResult{
			ID:         r.ID,
			Verdict:    r.Predict,
			Confidence: r.Confidence,
		}
	}

	return result, nil
}
