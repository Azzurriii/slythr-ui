package static_analysis

type AnalyzeRequest struct {
	Source string `json:"source" binding:"required" validate:"required,min=1"`
}
