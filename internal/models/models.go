package models

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`

	CreatedAt time.Time `json:"created_at"`
}

type Analysis struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	UploadedAt time.Time `json:"uploaded_at"`

	Findings []Finding `json:"findings"`
}

type Finding struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	AnalysisID uint `gorm:"index" json:"analysis_id"`

	FilePath string `gorm:"not null" json:"file_path"`
	Line     int    `gorm:"not null" json:"line"`
	Value    string `gorm:"not null" json:"value"`
	RuleID   string `gorm:"not null" json:"rule_id"`

	Severity          string  `json:"severity"`
	ScannerConfidence float64 `json:"scanner_confidence"`

	RuleVerdict    *string  `json:"rule_verdict"`
	RuleConfidence *float64 `json:"rule_confidence"`

	MlVerdict    *string  `json:"ml_verdict"`
	MlConfidence *float64 `json:"ml_confidence"`

	LlmVerdict     *string `json:"llm_verdict"`
	LlmExplanation *string `json:"llm_explanation"`

	FinalVerdict    *string  `json:"final_verdict"`
	FinalConfidence *float64 `json:"final_confidence"`

	HumanVerdict *string `json:"human_verdict"`
	HumanComment *string `json:"human_comment"`

	CreatedAt time.Time `json:"created_at"`
}
