package handlers

import (
	"net/http"
	"strings"

	"github.com/Azzurriii/slythr-go-backend/internal/application/dto/analysis"
	"github.com/Azzurriii/slythr-go-backend/internal/application/services"
	"github.com/gin-gonic/gin"
)

type StaticAnalysisHandler struct {
	service *services.StaticAnalysisService
}

func NewStaticAnalysisHandler(service *services.StaticAnalysisService) *StaticAnalysisHandler {
	return &StaticAnalysisHandler{
		service: service,
	}
}

// AnalyzeContract godoc
// @Summary Analyze Solidity contract for security vulnerabilities
// @Description Performs static security analysis on Solidity source code using Slither
// @Tags static-analysis
// @Accept json
// @Produce json
// @Param request body analysis.AnalyzeRequest true "Contract source code"
// @Router /static-analysis [post]
func (h *StaticAnalysisHandler) AnalyzeContract(c *gin.Context) {
	var req analysis.AnalyzeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	if strings.TrimSpace(req.SourceCode) == "" {
		h.respondError(c, http.StatusBadRequest, "Source code cannot be empty", nil)
		return
	}

	result, err := h.service.AnalyzeContract(c.Request.Context(), req.SourceCode)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "Failed to analyze contract", err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *StaticAnalysisHandler) respondError(c *gin.Context, status int, message string, err error) {
	response := gin.H{
		"success": false,
		"message": message,
	}

	if err != nil {
		response["error"] = err.Error()
	}

	c.JSON(status, response)
}
