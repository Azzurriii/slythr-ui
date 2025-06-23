package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	config "github.com/Azzurriii/slythr/config"
	"github.com/Azzurriii/slythr/internal/application/dto/analysis"
	"github.com/Azzurriii/slythr/internal/application/dto/testcase_generation"
	"github.com/Azzurriii/slythr/internal/domain/repository"
	"github.com/Azzurriii/slythr/internal/infrastructure/cache"
	"github.com/Azzurriii/slythr/internal/infrastructure/external"
	"github.com/Azzurriii/slythr/pkg/logger"
	"github.com/Azzurriii/slythr/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type TestCaseGenerationService struct {
	staticAnalysisService  *StaticAnalysisService
	dynamicAnalysisService *DynamicAnalysisService
	geminiClient           *external.GeminiClient
	logger                 *logger.Logger
	cache                  *cache.Cache
}

type TestCaseGenerationServiceOptions struct {
	GeminiModel   string
	GeminiTimeout time.Duration
}

func NewTestCaseGenerationService(
	staticAnalysisService *StaticAnalysisService,
	dynamicAnalysisService *DynamicAnalysisService,
	testCasesRepo repository.GeneratedTestCasesRepository,
	opts *TestCaseGenerationServiceOptions,
) (*TestCaseGenerationService, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	geminiOpts := &external.GeminiClientOptions{}
	if opts != nil {
		if opts.GeminiModel != "" {
			geminiOpts.Model = opts.GeminiModel
		}
		if opts.GeminiTimeout > 0 {
			geminiOpts.Timeout = opts.GeminiTimeout
		}
	}

	geminiClient := external.NewGeminiClient(cfg.Gemini, geminiOpts)

	// Setup cache
	var testCaseCache *cache.Cache
	if cfg.Redis.Addr != "" {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		testCaseCache = cache.NewCache(redisClient, nil, nil, nil, testCasesRepo)
	}

	return &TestCaseGenerationService{
		staticAnalysisService:  staticAnalysisService,
		dynamicAnalysisService: dynamicAnalysisService,
		geminiClient:           geminiClient,
		logger:                 logger.Default,
		cache:                  testCaseCache,
	}, nil
}

func (s *TestCaseGenerationService) GenerateTestCases(
	ctx context.Context,
	sourceCode, testFramework, testLanguage string,
) (*testcase_generation.TestCaseGenerateResponse, error) {
	s.logger.Infof("Starting test case generation for framework: %s, language: %s", testFramework, testLanguage)

	// Check cache first
	sourceHash := utils.GenerateSourceHash(sourceCode)
	if s.cache != nil {
		if cached, err := s.cache.GetTestCases(ctx, sourceHash); err == nil && cached != nil {
			s.logger.Infof("Returning cached test cases for source hash: %s", sourceHash)
			return cached, nil
		}
	}

	staticAnalysisResult := &analysis.StaticAnalysisResponse{}
	securityAnalysisResult := &analysis.DynamicAnalysisResponse{}

	// Get Analysis Cached Results
	if s.staticAnalysisService != nil && s.cache != nil {
		if cached, err := s.cache.GetStaticAnalysis(ctx, sourceHash); err == nil && cached != nil {
			staticAnalysisResult = cached
			s.logger.Debugf("Using cached static analysis for source hash: %s", sourceHash)
		} else {
			s.logger.Debugf("No cached static analysis found for source hash: %s", sourceHash)
		}
	}

	if s.dynamicAnalysisService != nil && s.cache != nil {
		if cached, err := s.cache.GetDynamicAnalysis(ctx, sourceHash); err == nil && cached != nil {
			securityAnalysisResult = cached
			s.logger.Debugf("Using cached dynamic analysis for source hash: %s", sourceHash)
		} else {
			s.logger.Debugf("No cached dynamic analysis found for source hash: %s", sourceHash)
		}
	}

	// Generate test cases using Gemini
	llmResponse, err := s.geminiClient.GenerateTestCases(
		ctx,
		sourceCode,
		testFramework,
		testLanguage,
		staticAnalysisResult,
		securityAnalysisResult,
	)
	if err != nil {
		errorResp := s.errorResponse(fmt.Sprintf("Failed to generate test cases: %v", err))
		errorResp.SourceHash = sourceHash
		errorResp.TestFramework = testFramework
		errorResp.TestLanguage = testLanguage

		if s.cache != nil {
			go s.cache.SetTestCases(context.Background(), sourceHash, errorResp)
		}

		return errorResp, err
	}

	testCode, llmWarnings := s.parseLLMResponse(llmResponse)

	contractName := s.extractContractName(sourceCode)
	fileName := s.generateFileName(contractName, testLanguage)

	analysisWarnings := s.generateWarningsAndRecommendations(staticAnalysisResult, securityAnalysisResult)

	allWarnings := append(analysisWarnings, llmWarnings...)

	response := &testcase_generation.TestCaseGenerateResponse{
		Success:                    true,
		TestCode:                   testCode,
		TestFramework:              testFramework,
		TestLanguage:               testLanguage,
		FileName:                   fileName,
		SourceHash:                 sourceHash,
		WarningsAndRecommendations: allWarnings,
		GeneratedAt:                time.Now(),
	}

	if s.cache != nil {
		go s.cache.SetTestCases(context.Background(), sourceHash, response)
	}

	s.logger.Infof("Successfully generated test cases for contract: %s", contractName)
	return response, nil
}

func (s *TestCaseGenerationService) generateWarningsAndRecommendations(
	staticResult *analysis.StaticAnalysisResponse,
	securityResult *analysis.DynamicAnalysisResponse,
) []string {
	var warnings []string

	if staticResult != nil && staticResult.Success {
		if staticResult.TotalIssues > 0 {
			warnings = append(warnings, fmt.Sprintf("Static analysis found %d security issues. Ensure your tests cover these vulnerabilities.", staticResult.TotalIssues))
		}

		if staticResult.SeveritySummary.High > 0 {
			warnings = append(warnings, fmt.Sprintf("Found %d high-severity issues. Priority testing recommended for these vulnerabilities.", staticResult.SeveritySummary.High))
		}
		if staticResult.SeveritySummary.Medium > 0 {
			warnings = append(warnings, fmt.Sprintf("Found %d medium-severity issues. Include edge case testing.", staticResult.SeveritySummary.Medium))
		}
	} else {
		warnings = append(warnings, "No static analysis data available. Consider running Slither analysis for better security coverage.")
	}

	if securityResult != nil && securityResult.Success {
		if securityResult.TotalIssues > 0 {
			warnings = append(warnings, fmt.Sprintf("AI security analysis identified %d potential vulnerabilities. Review test coverage for these areas.", securityResult.TotalIssues))
		}

		switch securityResult.Analysis.RiskLevel {
		case "HIGH":
			warnings = append(warnings, "High risk level detected. Implement comprehensive security testing including access control, reentrancy, and overflow tests.")
		case "MEDIUM":
			warnings = append(warnings, "Medium risk level detected. Include boundary testing and input validation tests.")
		case "LOW":
			warnings = append(warnings, "Low risk level detected. Focus on functionality and edge case testing.")
		}
	} else {
		warnings = append(warnings, "No AI security analysis data available. Tests will focus on general smart contract best practices.")
	}

	warnings = append(warnings, "Always run tests against multiple scenarios including edge cases and boundary conditions.")
	warnings = append(warnings, "Consider using fuzzing and property-based testing for comprehensive coverage.")

	warnings = append(warnings, s.getFrameworkSpecificRecommendations())

	return warnings
}

func (s *TestCaseGenerationService) getFrameworkSpecificRecommendations() string {
	return "Ensure your test environment matches production conditions including gas limits and network conditions."
}

func (s *TestCaseGenerationService) errorResponse(message string) *testcase_generation.TestCaseGenerateResponse {
	return &testcase_generation.TestCaseGenerateResponse{
		Success:       false,
		Message:       message,
		TestCode:      "",
		TestFramework: "",
		TestLanguage:  "",
		FileName:      "",
		SourceHash:    "",
		WarningsAndRecommendations: []string{
			"Test case generation failed. Please check your source code and try again.",
			"Ensure your contract compiles successfully before generating tests.",
		},
		GeneratedAt: time.Now(),
	}
}

func (s *TestCaseGenerationService) extractContractName(sourceCode string) string {
	contractRegex := regexp.MustCompile(`contract\s+(\w+)`)
	matches := contractRegex.FindStringSubmatch(sourceCode)
	if len(matches) > 1 {
		return matches[1]
	}
	return "Contract"
}

func (s *TestCaseGenerationService) generateFileName(contractName, testLanguage string) string {
	ext := s.getFileExtension(testLanguage)
	return fmt.Sprintf("%s.test.%s", contractName, ext)
}

func (s *TestCaseGenerationService) getFileExtension(testLanguage string) string {
	switch strings.ToLower(testLanguage) {
	case "javascript", "js":
		return "js"
	case "typescript", "ts":
		return "ts"
	case "solidity", "sol":
		return "sol"
	case "python", "py":
		return "py"
	default:
		return "js"
	}
}

func (s *TestCaseGenerationService) GetTestCases(ctx context.Context, sourceHash string) (*testcase_generation.TestCaseGenerateResponse, error) {
	if strings.TrimSpace(sourceHash) == "" {
		return nil, fmt.Errorf("source hash cannot be empty")
	}

	s.logger.Infof("Getting test cases for source hash: %s", sourceHash)

	if s.cache != nil {
		if cached, err := s.cache.GetTestCases(ctx, sourceHash); err == nil && cached != nil {
			s.logger.Infof("Returning cached test cases for source hash: %s", sourceHash)
			return cached, nil
		}
	}

	return nil, fmt.Errorf("test cases not found for source hash: %s", sourceHash)
}

func (s *TestCaseGenerationService) parseLLMResponse(llmResponse string) (string, []string) {
	var testCode string
	var warnings []string

	// Clean up the response first
	cleanResponse := s.cleanLLMResponse(llmResponse)

	testCodeSection := s.extractSection(cleanResponse, "TEST CODE")
	warningsSection := s.extractSection(cleanResponse, "WARNINGS AND RECOMMENDATIONS")

	if testCodeSection != "" {
		testCode = s.extractCodeFromSection(testCodeSection)
		// Validate extracted test code
		if !s.validateTestCode(testCode) {
			s.logger.Warnf("Generated test code appears to be invalid or incomplete")
			warnings = append(warnings, "Generated test code may be incomplete - please review carefully")
		}
	} else {
		s.logger.Warnf("No TEST CODE section found in LLM response")
		warnings = append(warnings, "No test code section found in response")
	}

	if warningsSection != "" {
		warnings = append(warnings, s.extractWarningsFromSection(warningsSection)...)
	} else {
		s.logger.Debugf("No WARNINGS section found in LLM response")
	}

	s.logger.Debugf("Parsed test code (%d chars) and %d warnings from LLM response", len(testCode), len(warnings))

	return testCode, warnings
}

func (s *TestCaseGenerationService) cleanLLMResponse(response string) string {
	// Remove excessive whitespace and normalize line endings
	response = strings.ReplaceAll(response, "\r\n", "\n")
	response = strings.ReplaceAll(response, "\r", "\n")

	// Remove any leading/trailing whitespace
	return strings.TrimSpace(response)
}

func (s *TestCaseGenerationService) extractSection(response, header string) string {
	lines := strings.Split(response, "\n")

	// Look for section headers with flexible matching
	headerPatterns := []string{
		"## " + header,
		"##" + header,
		"# " + header,
		"**" + header + "**",
		header + ":",
		strings.ToUpper(header),
		strings.ToLower(header),
	}

	for i, line := range lines {
		cleanLine := strings.TrimSpace(line)

		for _, pattern := range headerPatterns {
			if strings.Contains(strings.ToUpper(cleanLine), strings.ToUpper(pattern)) {
				var sectionLines []string
				for j := i + 1; j < len(lines); j++ {
					nextLine := strings.TrimSpace(lines[j])
					// Stop at next major section
					if s.isNewSection(nextLine) {
						break
					}
					sectionLines = append(sectionLines, lines[j])
				}
				return strings.Join(sectionLines, "\n")
			}
		}
	}
	return ""
}

func (s *TestCaseGenerationService) isNewSection(line string) bool {
	line = strings.TrimSpace(strings.ToUpper(line))
	sectionHeaders := []string{
		"## TEST CODE",
		"## WARNINGS",
		"## RECOMMENDATIONS",
		"## ANALYSIS",
		"# TEST CODE",
		"# WARNINGS",
		"**TEST CODE**",
		"**WARNINGS",
	}

	for _, header := range sectionHeaders {
		if strings.Contains(line, header) {
			return true
		}
	}
	return false
}

func (s *TestCaseGenerationService) extractCodeFromSection(section string) string {
	lines := strings.Split(section, "\n")
	var codeLines []string
	inCodeBlock := false
	codeBlockCount := 0

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Handle code block markers
		if strings.HasPrefix(trimmedLine, "```") {
			if !inCodeBlock {
				inCodeBlock = true
				codeBlockCount++
				continue
			} else {
				inCodeBlock = false
				// If we found code, break after first complete code block
				if len(codeLines) > 0 {
					break
				}
				continue
			}
		}

		if inCodeBlock {
			// Skip language identifier lines
			if i > 0 && strings.HasPrefix(strings.TrimSpace(lines[i-1]), "```") {
				if s.isLanguageIdentifier(trimmedLine) {
					continue
				}
			}
			codeLines = append(codeLines, line)
		}
	}

	code := strings.TrimSpace(strings.Join(codeLines, "\n"))

	// If no code blocks found, try to extract code without markers
	if code == "" {
		code = s.extractCodeWithoutMarkers(section)
	}

	return code
}

func (s *TestCaseGenerationService) isLanguageIdentifier(line string) bool {
	line = strings.ToLower(strings.TrimSpace(line))
	languages := []string{
		"javascript", "js", "typescript", "ts", "solidity", "sol",
		"python", "py", "go", "rust", "java", "c++", "cpp",
	}

	for _, lang := range languages {
		if line == lang {
			return true
		}
	}
	return false
}

func (s *TestCaseGenerationService) extractCodeWithoutMarkers(section string) string {
	lines := strings.Split(section, "\n")
	var codeLines []string

	// Look for lines that appear to be code (indented or contain code patterns)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Skip lines that look like comments or explanations
		if strings.HasPrefix(trimmed, "-") || strings.HasPrefix(trimmed, "*") ||
			strings.HasPrefix(trimmed, "Note:") || strings.HasPrefix(trimmed, "//") {
			continue
		}

		// Include lines that look like code
		if s.looksLikeCode(line) {
			codeLines = append(codeLines, line)
		}
	}

	return strings.Join(codeLines, "\n")
}

func (s *TestCaseGenerationService) looksLikeCode(line string) bool {
	trimmed := strings.TrimSpace(line)

	// Check for common code patterns
	codePatterns := []string{
		"const ", "let ", "var ", "function", "describe(", "it(", "expect(",
		"require(", "import ", "contract ", "pragma ", "beforeEach", "await ",
		"assert", "=", "{", "}", "(", ")", ";",
	}

	for _, pattern := range codePatterns {
		if strings.Contains(trimmed, pattern) {
			return true
		}
	}

	// Check if line is indented (likely code)
	return len(line) > 0 && (line[0] == ' ' || line[0] == '\t') && len(strings.TrimSpace(line)) > 0
}

func (s *TestCaseGenerationService) validateTestCode(code string) bool {
	if strings.TrimSpace(code) == "" {
		return false
	}

	// Basic validation - check for common test patterns
	testPatterns := []string{
		"describe", "it(", "test(", "expect", "assert", "contract", "beforeEach",
	}

	codeUpper := strings.ToUpper(code)
	for _, pattern := range testPatterns {
		if strings.Contains(codeUpper, strings.ToUpper(pattern)) {
			return true
		}
	}

	// If no test patterns found, it might still be valid but warn
	return len(code) > 50 // Arbitrary minimum length
}

func (s *TestCaseGenerationService) extractWarningsFromSection(section string) []string {
	var warnings []string
	lines := strings.Split(section, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Handle different bullet point styles and markdown
		prefixes := []string{"- ", "* ", "• ", "◦ ", "▪ ", "▫ "}
		var warning string

		for _, prefix := range prefixes {
			if strings.HasPrefix(line, prefix) {
				warning = strings.TrimSpace(strings.TrimPrefix(line, prefix))
				break
			}
		}

		// Handle numbered lists
		if warning == "" {
			re := regexp.MustCompile(`^\d+\.\s+(.+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				warning = strings.TrimSpace(matches[1])
			}
		}

		// Handle lines without bullets if they look like warnings
		if warning == "" && s.looksLikeWarning(line) {
			warning = line
		}

		if warning != "" {
			// Clean up markdown formatting
			warning = s.cleanWarningText(warning)
			warnings = append(warnings, warning)
		}
	}

	return warnings
}

func (s *TestCaseGenerationService) looksLikeWarning(line string) bool {
	line = strings.ToLower(line)
	warningKeywords := []string{
		"warning", "recommendation", "note", "important", "caution",
		"consider", "ensure", "avoid", "check", "review", "test",
	}

	for _, keyword := range warningKeywords {
		if strings.Contains(line, keyword) {
			return true
		}
	}
	return false
}

func (s *TestCaseGenerationService) cleanWarningText(text string) string {
	// Remove markdown formatting
	text = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`\*(.*?)\*`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile("`(.*?)`").ReplaceAllString(text, "$1")

	// Remove extra spaces
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}
