package handlers

import (
	"net/http"
	"strings"

	"github.com/Azzurriii/slythr/internal/application/dto/analysis"
	"github.com/Azzurriii/slythr/internal/application/services"
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

// GetStaticAnalysis godoc
// @Summary Get static analysis result by source hash
// @Description Retrieves static analysis result from cache or database using source hash
// @Tags static-analysis
// @Accept json
// @Produce json
// @Param sourceHash path string true "Source hash"
// @Router /static-analysis/{sourceHash} [get]
func (h *StaticAnalysisHandler) GetStaticAnalysis(c *gin.Context) {
	sourceHash := c.Param("sourceHash")

	if strings.TrimSpace(sourceHash) == "" {
		h.respondError(c, http.StatusBadRequest, "Source hash cannot be empty", nil)
		return
	}

	result, err := h.service.GetStaticAnalysis(c.Request.Context(), sourceHash)
	if err != nil {
		h.respondError(c, http.StatusNotFound, "Static analysis not found", err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetContainerStatus godoc
// @Summary Get Slither container status
// @Description Check if Slither container is running and accessible
// @Tags static-analysis
// @Accept json
// @Produce json
// @Router /static-analysis/status [get]
func (h *StaticAnalysisHandler) GetContainerStatus(c *gin.Context) {
	status, err := h.service.GetContainerStatus()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Failed to get container status",
			"error":   err.Error(),
			"data":    status,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Container status retrieved successfully",
		"data":    status,
	})
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
