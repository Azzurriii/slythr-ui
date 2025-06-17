package static_analysis

import "time"

type SlitherIssue struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Confidence  string `json:"confidence"`
	Location    string `json:"location"`
	Reference   string `json:"reference"`
}

type SeveritySummary struct {
	High          int `json:"high"`
	Medium        int `json:"medium"`
	Low           int `json:"low"`
	Informational int `json:"informational"`
}

type AnalyzeResponse struct {
	Success         bool            `json:"success"`
	Message         string          `json:"message,omitempty"`
	Issues          []SlitherIssue  `json:"issues"`
	TotalIssues     int             `json:"total_issues"`
	SeveritySummary SeveritySummary `json:"severity_summary"`
	AnalyzedAt      time.Time       `json:"analyzed_at"`
}
