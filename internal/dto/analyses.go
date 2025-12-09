package dto

type UploadAnalysisResponse struct {
	AnalysisID uint `json:"analysis_id"`
}

type AnalysisListItem struct {
	ID       uint   `json:"id"`
	Uploaded string `json:"uploaded_at"`
	Findings int    `json:"findings_count"`
}

type FindingResponse struct {
	ID              uint     `json:"id"`
	FilePath        string   `json:"file_path"`
	Line            int      `json:"line"`
	Value           string   `json:"value"`
	RuleID          string   `json:"rule_id"`
	Severity        string   `json:"severity"`
	FinalVerdict    *string  `json:"final_verdict"`
	FinalConfidence *float64 `json:"final_confidence"`
}
