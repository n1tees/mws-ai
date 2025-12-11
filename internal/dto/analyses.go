package dto

type UploadAnalysisResponse struct {
	AnalysisID uint `json:"analysis_id" example:"42"`
}

type FindingItem struct {
	ID              uint    `json:"id" example:"10"`
	FilePath        string  `json:"file_path"`
	Line            int     `json:"line"`
	RuleID          string  `json:"rule_id"`
	Severity        string  `json:"severity"`
	FinalVerdict    string  `json:"final_verdict"`
	FinalConfidence float64 `json:"final_confidence"`
}

type AnalysisResponse struct {
	ID         uint          `json:"id" example:"42"`
	UserID     uint          `json:"user_id"`
	UploadedAt string        `json:"uploaded_at"`
	Findings   []FindingItem `json:"findings"`
}

type AnalysisListItem struct {
	ID         uint   `json:"id" example:"42"`
	Status     string `json:"status" example:"done"`
	UploadedAt string `json:"uploaded_at"`
}
