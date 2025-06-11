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

type AnalyzeResponse struct {
	Success     bool           `json:"success"`
	Message     string         `json:"message,omitempty"`
	Issues      []SlitherIssue `json:"issues"`
	TotalIssues int            `json:"total_issues"`
	AnalyzedAt  time.Time      `json:"analyzed_at"`
	RawOutput   string         `json:"raw_output,omitempty"`
}
