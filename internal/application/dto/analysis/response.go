package analysis

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

type DynamicAnalysisResponse struct {
	Success        bool        `json:"success"`
	Message        string      `json:"message,omitempty"`
	Analysis       LLMAnalysis `json:"analysis,omitempty"`
	TotalIssues    int         `json:"total_issues"`
	AnalyzedAt     time.Time   `json:"analyzed_at"`
	RawLLMResponse string      `json:"raw_llm_response,omitempty"`
}

type SlitherIssue struct {
	Type        string `json:"type" example:"detector"`
	Title       string `json:"title" example:"Shadowing Local"`
	Description string `json:"description" example:"Variable shadows another variable"`
	Severity    string `json:"severity" example:"LOW" enums:"HIGH,MEDIUM,LOW,INFORMATIONAL"`
	Confidence  string `json:"confidence" example:"High" enums:"High,Medium,Low"`
	Location    string `json:"location" example:"Contract.sol:L42"`
	Reference   string `json:"reference" example:"Contract.sol#L42"`
}

type SeveritySummary struct {
	High          int `json:"high" example:"0"`
	Medium        int `json:"medium" example:"1"`
	Low           int `json:"low" example:"2"`
	Informational int `json:"informational" example:"5"`
}

type StaticAnalysisResponse struct {
	Success         bool            `json:"success" example:"true"`
	Message         string          `json:"message,omitempty" example:""`
	Issues          []SlitherIssue  `json:"issues"`
	TotalIssues     int             `json:"total_issues" example:"8"`
	SeveritySummary SeveritySummary `json:"severity_summary"`
	AnalyzedAt      time.Time       `json:"analyzed_at" example:"2025-06-17T16:19:00.579573+07:00"`
}
