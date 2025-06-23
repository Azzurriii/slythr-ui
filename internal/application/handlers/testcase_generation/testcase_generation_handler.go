package handlers

import (
	"net/http"
	"strings"

	"github.com/Azzurriii/slythr/internal/application/dto/testcase_generation"
	"github.com/Azzurriii/slythr/internal/application/services"
	"github.com/gin-gonic/gin"
)

type TestCaseGenerationHandler struct {
	service *services.TestCaseGenerationService
}

func NewTestCaseGenerationHandler(service *services.TestCaseGenerationService) *TestCaseGenerationHandler {
	return &TestCaseGenerationHandler{
		service: service,
	}
}

// GenerateTestCases godoc
// @Summary Generate test cases for Solidity contract
// @Description Generates test cases for Solidity contract using Gemini AI with comprehensive analysis
// @Tags testcase-generation
// @Accept json
// @Produce json
// @Param request body testcase_generation.TestCaseGenerateRequest true "Contract source code and test preferences"
// @Router /test-cases/generate [post]
func (h *TestCaseGenerationHandler) GenerateTestCases(c *gin.Context) {
	var req testcase_generation.TestCaseGenerateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	// Validate required fields
	if strings.TrimSpace(req.SourceCode) == "" {
		h.respondError(c, http.StatusBadRequest, "Source code cannot be empty", nil)
		return
	}

	if strings.TrimSpace(req.TestFramework) == "" {
		h.respondError(c, http.StatusBadRequest, "Test framework cannot be empty", nil)
		return
	}

	if strings.TrimSpace(req.TestLanguage) == "" {
		h.respondError(c, http.StatusBadRequest, "Test language cannot be empty", nil)
		return
	}

	// Validate supported test frameworks and languages
	if !h.isValidTestFramework(req.TestFramework) {
		h.respondError(c, http.StatusBadRequest, "Unsupported test framework. Supported frameworks: hardhat, truffle, foundry, brownie", nil)
		return
	}

	if !h.isValidTestLanguage(req.TestLanguage) {
		h.respondError(c, http.StatusBadRequest, "Unsupported test language. Supported languages: javascript, typescript, solidity, python", nil)
		return
	}

	result, err := h.service.GenerateTestCases(
		c.Request.Context(),
		req.SourceCode,
		req.TestFramework,
		req.TestLanguage,
	)

	// Even if there's an error, if we have a structured response, return it
	if result != nil {
		if result.Success {
			c.JSON(http.StatusOK, result)
		} else {
			c.JSON(http.StatusInternalServerError, result)
		}
		return
	}

	// Fallback error response
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "Failed to generate test cases", err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetTestCases godoc
// @Summary Get test cases for Solidity contract
// @Description Gets test cases for Solidity contract by source hash
// @Tags testcase-generation
// @Accept json
// @Produce json
// @Param sourceHash path string true "Source hash"
// @Router /test-cases/{sourceHash} [get]
func (h *TestCaseGenerationHandler) GetTestCases(c *gin.Context) {
	sourceHash := c.Param("sourceHash")

	if strings.TrimSpace(sourceHash) == "" {
		h.respondError(c, http.StatusBadRequest, "Source hash cannot be empty", nil)
		return
	}

	testCases, err := h.service.GetTestCases(c.Request.Context(), sourceHash)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.respondError(c, http.StatusNotFound, "Test cases not found for the given source hash", err)
		} else {
			h.respondError(c, http.StatusInternalServerError, "Failed to get test cases", err)
		}
		return
	}

	c.JSON(http.StatusOK, testCases)
}

func (h *TestCaseGenerationHandler) isValidTestFramework(framework string) bool {
	validFrameworks := map[string]bool{
		"hardhat": true,
		"truffle": true,
		"foundry": true,
		"brownie": true,
	}
	return validFrameworks[strings.ToLower(framework)]
}

func (h *TestCaseGenerationHandler) isValidTestLanguage(language string) bool {
	validLanguages := map[string]bool{
		"javascript": true,
		"js":         true,
		"typescript": true,
		"ts":         true,
		"solidity":   true,
		"sol":        true,
		"python":     true,
		"py":         true,
	}
	return validLanguages[strings.ToLower(language)]
}

func (h *TestCaseGenerationHandler) respondError(c *gin.Context, status int, message string, err error) {
	response := gin.H{
		"success": false,
		"message": message,
	}

	if err != nil {
		response["error"] = err.Error()
	}

	c.JSON(status, response)
}
