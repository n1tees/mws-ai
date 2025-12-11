package models

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Analysis struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"index" json:"user_id"`

	FilePath string `json:"file_path"`
	Status   string `json:"status"` // pending / processing / done / failed

	UploadedAt time.Time `gorm:"autoCreateTime" json:"uploaded_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Findings []Finding `gorm:"constraint:OnDelete:CASCADE;" json:"findings"`

	// SUMMARY FIELDS
	FinalVerdict    *string  `json:"final_verdict"`
	TPCount         int      `json:"tp_count"`
	FPCount         int      `json:"fp_count"`
	FinalConfidence *float64 `json:"final_confidence"`
}

type Finding struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	AnalysisID uint `gorm:"index" json:"analysis_id"`

	FilePath string `gorm:"not null" json:"file_path"`
	Line     int    `gorm:"not null" json:"line"`
	// Optional: SARIF иногда даёт диапазон строк
	LineEnd *int `json:"line_end,omitempty"`

	Value  string `gorm:"not null" json:"value"`
	RuleID string `gorm:"not null" json:"rule_id"`

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

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type ApiKey struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"index"`
	Hash       string `gorm:"uniqueIndex"`
	CreatedAt  time.Time
	LastUsedAt *time.Time
	Active     bool
}
