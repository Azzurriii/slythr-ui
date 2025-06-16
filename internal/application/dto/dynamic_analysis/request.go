package dynamic_analysis

type AnalyzeRequest struct {
	SourceCode string `json:"source_code" binding:"required" validate:"required,min=1"`
}
