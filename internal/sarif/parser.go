package sarif

import (
	"encoding/json"
	"fmt"
	"os"

	"mws-ai/internal/models"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

// Реализация интерфейса SarifParser
func (p *Parser) Parse(path string) ([]models.Finding, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read sarif file: %w", err)
	}

	var sarif Sarif
	if err := json.Unmarshal(data, &sarif); err != nil {
		return nil, fmt.Errorf("parse sarif json: %w", err)
	}

	return ConvertToFindings(sarif), nil
}

func ConvertToFindings(s Sarif) []models.Finding {
	findings := []models.Finding{}

	if len(s.Runs) == 0 {
		return findings
	}

	for _, result := range s.Runs[0].Results {

		for _, loc := range result.Locations {

			f := models.Finding{
				// AnalysisID НЕ заполняем — сервис сделает сам!

				FilePath: loc.PhysicalLocation.ArtifactLocation.URI,
				Line:     loc.PhysicalLocation.Region.StartLine,

				Value:  result.Properties.Snippet,
				RuleID: result.RuleID,

				Severity:          result.Properties.Severity,
				ScannerConfidence: result.Properties.Confidence,

				RuleVerdict: &result.Properties.AIVerdict,
			}

			findings = append(findings, f)
		}
	}

	return findings
}
