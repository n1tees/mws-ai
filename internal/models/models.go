package models

import (
	"time"
)

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

	TPCount int `json:"tp_count"`
	FPCount int `json:"fp_count"`

	UploadedAt time.Time `gorm:"autoCreateTime" json:"uploaded_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Findings []Finding `gorm:"constraint:OnDelete:CASCADE;" json:"findings"`
}

type Finding struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	AnalysisID uint `gorm:"index" json:"analysis_id"`

	FilePath string `gorm:"not null" json:"file_path"`
	Line     int    `gorm:"not null" json:"line"`
	LineEnd  *int   `json:"line_end,omitempty"`

	Value  string `gorm:"not null" json:"value"`
	RuleID string `gorm:"not null" json:"rule_id"`

	Severity          string  `json:"severity"`
	ScannerConfidence float64 `json:"scanner_confidence"`

	// Heuristic results
	HeuristicTriggered bool     `json:"heuristic_triggered"`
	HeuristicReason    *string  `json:"heuristic_reason,omitempty"`
	EntropyClass       *string  `json:"entropy_class,omitempty"`
	EntropyValue       *float64 `json:"entropy,omitempty"`

	// ML results
	MlVerdict    *string  `json:"ml_verdict"`
	MlConfidence *float64 `json:"ml_confidence"`

	// LLM results
	LlmVerdict     *string  `json:"llm_verdict"`
	LlmConfidence  *float64 `json:"llm_confidence"`
	LlmExplanation *string  `json:"llm_explanation"`

	// Final
	FinalVerdict   *string `json:"final_verdict"`
	DecisionSource string  `gorm:"type:varchar(25)" json:"decision_source"`

	// Human review
	HumanVerdict *string `json:"human_verdict"`
	HumanComment *string `json:"human_comment"`

	Status string `gorm:"type:varchar(32);default:'pending'" json:"status"` // pernding, processed, error, review

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type ApiKey struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint   `gorm:"index"`
	Hash   string `gorm:"uniqueIndex"`

	Type   string `gorm:"index"` // "service", "admin", "integration"
	Active bool   `gorm:"index"`

	CreatedAt  time.Time
	LastUsedAt *time.Time
	ExpiresAt  *time.Time
}
