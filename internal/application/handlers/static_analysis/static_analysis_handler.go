package handlers

import (
	"net/http"

	"github.com/Azzurriii/slythr-go-backend/internal/application/dto/static_analysis"
	"github.com/Azzurriii/slythr-go-backend/internal/application/services"
	"github.com/gin-gonic/gin"
)

type StaticAnalysisHandler struct {
	staticAnalysisService *services.StaticAnalysisService
}

func NewStaticAnalysisHandler() *StaticAnalysisHandler {
	return &StaticAnalysisHandler{
		staticAnalysisService: services.NewStaticAnalysisService(),
	}
}

// AnalyzeContract godoc
// @Summary Analyze Solidity contract for security vulnerabilities
// @Description Performs static security analysis on Solidity source code using Slither
// @Tags static-analysis
// @Accept json
// @Produce json
// @Param request body static_analysis.AnalyzeRequest true "Contract source code"
// @Router /static-analysis [post]
func (h *StaticAnalysisHandler) AnalyzeContract(c *gin.Context) {
	var req static_analysis.AnalyzeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Validate source code is not empty
	if len(req.SourceCode) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Source code cannot be empty",
		})
		return
	}

	// Perform static analysis
	result, err := h.staticAnalysisService.AnalyzeContract(req.SourceCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to analyze contract",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
