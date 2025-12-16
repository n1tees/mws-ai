package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mws-ai/internal/models"
)

type HeuristicResult struct {
	ID         uint     `json:"id"`
	Verdict    string   `json:"verdict"`
	Confidence float64  `json:"confidence"`
	Reasons    []string `json:"reasons"`
}

type HeuristicClient interface {
	Analyze(findings []*models.Finding) (map[uint]HeuristicResult, error)
}

type heuristicHTTP struct {
	baseURL string
	client  *http.Client
}

func NewHeuristicClient(baseURL string) HeuristicClient {
	return &heuristicHTTP{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 20 * time.Second},
	}
}

func (h *heuristicHTTP) Analyze(findings []*models.Finding) (map[uint]HeuristicResult, error) {

	// === Формируем JSON-массив ===
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

	// === POST в контейнер эвристики ===
	resp, err := h.client.Post(h.baseURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("heuristic request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("heuristic service returned %d", resp.StatusCode)
	}

	// === Ответ ===
	var out struct {
		Results []HeuristicResult `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode heuristic response: %w", err)
	}

	// === Конвертируем в map[id]result ===
	resultMap := make(map[uint]HeuristicResult)

	for _, r := range out.Results {
		resultMap[r.ID] = r
	}

	return resultMap, nil
}
