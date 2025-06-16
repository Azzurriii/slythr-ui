package handlers

import (
	"net/http"

	"github.com/Azzurriii/slythr-go-backend/internal/application/dto/dynamic_analysis"
	"github.com/Azzurriii/slythr-go-backend/internal/application/services"
	"github.com/gin-gonic/gin"
)

type DynamicAnalysisHandler struct {
	dynamicAnalysisService *services.DynamicAnalysisService
}

func NewDynamicAnalysisHandler(dynamicAnalysisService *services.DynamicAnalysisService) *DynamicAnalysisHandler {
	return &DynamicAnalysisHandler{
		dynamicAnalysisService: dynamicAnalysisService,
	}
}

// AnalyzeContract godoc
// @Summary Analyze Solidity contract for security vulnerabilities using LLM
// @Description Performs dynamic security analysis on Solidity source code using Gemini LLM
// @Tags dynamic-analysis
// @Accept json
// @Produce json
// @Param request body dynamic_analysis.AnalyzeRequest true "Contract source code"
// @Router /dynamic-analysis [post]
func (h *DynamicAnalysisHandler) AnalyzeContract(c *gin.Context) {
	var req dynamic_analysis.AnalyzeRequest

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

	// Perform dynamic analysis using LLM with context
	result, err := h.dynamicAnalysisService.AnalyzeContract(c.Request.Context(), req.SourceCode)
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
