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

func (p *Parser) Parse(path string) ([]models.Finding, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read sarif file: %w", err)
	}

	var sarif Sarif
	if err := json.Unmarshal(data, &sarif); err != nil {
		return nil, fmt.Errorf("parse sarif json: %w", err)
	}

	return convertToFindings(sarif), nil
}

func convertToFindings(s Sarif) []models.Finding {
	findings := make([]models.Finding, 0)

	if len(s.Runs) == 0 {
		return findings
	}

	for _, result := range s.Runs[0].Results {
		for _, loc := range result.Locations {

			line := loc.PhysicalLocation.Region.StartLine
			var lineEnd *int
			if loc.PhysicalLocation.Region.EndLine != 0 {
				v := loc.PhysicalLocation.Region.EndLine
				lineEnd = &v
			}

			value := result.Properties.Snippet
			if value == "" {
				value = result.Message.Text
			}

			f := models.Finding{
				FilePath: loc.PhysicalLocation.ArtifactLocation.URI,
				Line:     line,
				LineEnd:  lineEnd,

				Value:  value,
				RuleID: result.RuleID,

				Severity:          result.Properties.Severity,
				ScannerConfidence: result.Properties.Confidence,
			}

			findings = append(findings, f)
		}
	}

	return findings
}
