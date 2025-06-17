package dynamic_analysis

import "time"

type Vulnerability struct {
	Title          string      `json:"title"`
	Severity       string      `json:"severity"`
	Description    string      `json:"description"`
	Location       interface{} `json:"location,omitempty"`
	Recommendation interface{} `json:"recommendation"`
}

type LLMAnalysis struct {
	SecurityScore   int             `json:"security_score"`
	RiskLevel       string          `json:"risk_level"`
	Summary         string          `json:"summary"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	GoodPractices   interface{}     `json:"good_practices"`
	Recommendations interface{}     `json:"recommendations"`
}

type AnalyzeResponse struct {
	Success        bool        `json:"success"`
	Message        string      `json:"message,omitempty"`
	Analysis       LLMAnalysis `json:"analysis,omitempty"`
	TotalIssues    int         `json:"total_issues"`
	AnalyzedAt     time.Time   `json:"analyzed_at"`
	RawLLMResponse string      `json:"raw_llm_response,omitempty"`
}
