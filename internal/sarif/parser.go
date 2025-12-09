package sarif

import (
	"encoding/json"
	"fmt"
	"os"

	"mws-ai/internal/models"
)

func ParseFile(path string, analysisID uint) ([]models.Finding, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read sarif file: %w", err)
	}

	var sarif Sarif
	if err := json.Unmarshal(data, &sarif); err != nil {
		return nil, fmt.Errorf("parse sarif json: %w", err)
	}

	return ConvertToFindings(sarif, analysisID), nil
}

func ConvertToFindings(s Sarif, analysisID uint) []models.Finding {
	findings := []models.Finding{}

	if len(s.Runs) == 0 {
		return findings
	}

	for _, result := range s.Runs[0].Results {

		// SARIF может иметь несколько locations (но обычно 1)
		for _, loc := range result.Locations {

			f := models.Finding{
				AnalysisID: analysisID,

				FilePath: loc.PhysicalLocation.ArtifactLocation.URI,
				Line:     loc.PhysicalLocation.Region.StartLine,

				Value:  result.Properties.Snippet,
				RuleID: result.RuleID,

				Severity:          result.Properties.Severity,
				ScannerConfidence: result.Properties.Confidence,

				// Алгоритм:
				// RuleVerdict = aiVerdict (предварительно)
				RuleVerdict: &result.Properties.AIVerdict,
			}

			findings = append(findings, f)
		}
	}

	return findings
}
