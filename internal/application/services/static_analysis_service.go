package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	config "github.com/Azzurriii/slythr/config"
	"github.com/Azzurriii/slythr/internal/application/dto/analysis"
	"github.com/Azzurriii/slythr/internal/domain/repository"
	"github.com/Azzurriii/slythr/internal/infrastructure/cache"
	"github.com/Azzurriii/slythr/pkg/logger"
	"github.com/Azzurriii/slythr/pkg/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	defaultContainerName = "slither"
	defaultWorkspacePath = "/workspace"
	analysisTimeout      = 5 * time.Minute
	maxDescriptionLength = 200
)

type StaticAnalysisService struct {
	containerName string
	workspacePath string
	logger        *logger.Logger
	cache         *cache.Cache
}

type StaticAnalysisServiceOptions struct {
	AnalysisTimeout time.Duration
	ContainerName   string
	WorkspacePath   string
}

func NewStaticAnalysisService(
	staticAnalysisRepo repository.StaticAnalysisRepository,
	opts *StaticAnalysisServiceOptions,
) (*StaticAnalysisService, error) {
	containerName := defaultContainerName
	workspacePath := defaultWorkspacePath

	if opts != nil {
		if opts.ContainerName != "" {
			containerName = opts.ContainerName
		}
		if opts.WorkspacePath != "" {
			workspacePath = opts.WorkspacePath
		}
	}

	if envContainer := os.Getenv("SLITHER_CONTAINER_NAME"); envContainer != "" {
		containerName = envContainer
	}
	if envWorkspace := os.Getenv("WORKSPACE_PATH"); envWorkspace != "" {
		workspacePath = envWorkspace
	}

	var analysisCache *cache.Cache
	cfg, err := config.LoadConfig()
	if err == nil {
		var redisClient *redis.Client
		if cfg.Redis.Addr != "" {
			redisClient = redis.NewClient(&redis.Options{
				Addr:     cfg.Redis.Addr,
				Password: cfg.Redis.Password,
				DB:       cfg.Redis.DB,
			})
		}

		analysisCache = cache.NewCache(redisClient, nil, nil, staticAnalysisRepo, nil)
	}

	return &StaticAnalysisService{
		containerName: containerName,
		workspacePath: workspacePath,
		logger:        logger.Default,
		cache:         analysisCache,
	}, nil
}

func (s *StaticAnalysisService) AnalyzeContract(ctx context.Context, source string) (*analysis.StaticAnalysisResponse, error) {
	if !s.isContainerRunning() {
		return s.errorResponse("Slither analysis container is not running. Please start the container first."),
			fmt.Errorf("slither container not running")
	}

	sourceHash := utils.GenerateSourceHash(source)
	if s.cache != nil {
		if cached, err := s.cache.GetStaticAnalysis(ctx, sourceHash); err == nil && cached != nil {
			s.logger.Infof("Returning cached static analysis for source hash: %s", sourceHash)
			return cached, nil
		}
	}

	analysisID := uuid.New().String()
	s.logger.Infof("Starting static analysis with ID: %s, source hash: %s", analysisID, sourceHash)

	tempDir, err := s.setupWorkspace(analysisID, source)
	if err != nil {
		return s.errorResponse(fmt.Sprintf("Failed to setup workspace: %v", err)), err
	}
	defer os.RemoveAll(tempDir)

	// Copy to container and install dependencies
	containerPath := filepath.Join(s.workspacePath, analysisID)
	if err := s.prepareContainer(tempDir, containerPath); err != nil {
		return s.errorResponse(fmt.Sprintf("Failed to prepare container: %v", err)), err
	}

	// Run analysis
	slitherOutput, err := s.runSlitherAnalysis(containerPath)
	if err != nil {
		return s.errorResponse(fmt.Sprintf("Failed to run Slither: %v", err)), err
	}

	// Parse and return results
	response := s.buildSuccessResponse(slitherOutput, sourceHash)

	// Cache the result
	if s.cache != nil {
		go s.cache.SetStaticAnalysis(context.Background(), sourceHash, response)
	}

	defer func() {
		if r := recover(); r != nil {
			s.logger.Errorf("Panic during analysis: %v", r)
			os.RemoveAll(tempDir)
		}
	}()

	defer func() {
		go func() {
			s.cleanupContainerWorkspace(containerPath)
		}()
	}()

	return response, nil
}

// GetStaticAnalysis retrieves static analysis result by source hash
func (s *StaticAnalysisService) GetStaticAnalysis(ctx context.Context, sourceHash string) (*analysis.StaticAnalysisResponse, error) {
	if strings.TrimSpace(sourceHash) == "" {
		return nil, fmt.Errorf("source hash cannot be empty")
	}

	if s.cache != nil {
		if cached, err := s.cache.GetStaticAnalysis(ctx, sourceHash); err == nil && cached != nil {
			s.logger.Infof("Returning cached static analysis for source hash: %s", sourceHash)
			return cached, nil
		}
	}

	return nil, fmt.Errorf("static analysis not found for source hash: %s", sourceHash)
}

func (s *StaticAnalysisService) errorResponse(message string) *analysis.StaticAnalysisResponse {
	return &analysis.StaticAnalysisResponse{
		Success: false,
		Message: message,
	}
}

func (s *StaticAnalysisService) setupWorkspace(analysisID, source string) (string, error) {
	tempDir := filepath.Join(os.TempDir(), analysisID)

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}

	// Write contract file
	contractFile := filepath.Join(tempDir, "Contract.sol")
	if err := os.WriteFile(contractFile, []byte(source), 0644); err != nil {
		return "", fmt.Errorf("failed to write contract file: %v", err)
	}

	// Write package.json
	packageJSON := s.getPackageJSON()
	if err := os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0644); err != nil {
		return "", fmt.Errorf("failed to write package.json: %v", err)
	}

	return tempDir, nil
}

func (s *StaticAnalysisService) getPackageJSON() string {
	return `{
		"name": "slither-analysis",
		"version": "1.0.0",
		"dependencies": {
			"@openzeppelin/contracts": "^4.9.0",
			"@openzeppelin/contracts-upgradeable": "^4.9.0",
			"@chainlink/contracts": "^0.6.1",
			"@uniswap/v2-core": "^1.0.1",
			"@uniswap/v3-core": "^1.0.0",
			"@aave/protocol-v2": "^1.0.0",
			"@aave/core-v3": "^1.16.2"
		}
	}`
}

func (s *StaticAnalysisService) prepareContainer(tempDir, containerPath string) error {
	go func() {
		copyCmd := exec.Command("docker", "cp", tempDir, fmt.Sprintf("%s:%s", s.containerName, s.workspacePath))
		if err := copyCmd.Run(); err != nil {
			s.logger.Errorf("failed to copy files to container: %v", err)
		}

		installCmd := exec.Command("docker", "exec", s.containerName, "bash", "-c",
			fmt.Sprintf("cd %s && npm install", containerPath))
		if err := installCmd.Run(); err != nil {
			s.logger.Errorf("failed to install dependencies: %v", err)
		}
	}()

	return nil
}

func (s *StaticAnalysisService) runSlitherAnalysis(containerPath string) (string, error) {
	solcVersion, err := s.detectSolidityVersion(containerPath, "Contract.sol")
	if err != nil {
		s.logger.Warnf("Failed to detect Solidity version, using default: %v", err)
		solcVersion = "0.8.20"
	}

	return s.runSlitherInContainer(containerPath, "Contract.sol", solcVersion)
}

func (s *StaticAnalysisService) buildSuccessResponse(slitherOutput string, sourceHash string) *analysis.StaticAnalysisResponse {
	issues := s.parseSlitherOutput(slitherOutput)
	severitySummary := s.calculateSeveritySummary(issues)

	s.logger.Infof("Static analysis completed successfully, found %d issues", len(issues))

	return &analysis.StaticAnalysisResponse{
		Success:         true,
		Issues:          issues,
		TotalIssues:     len(issues),
		SeveritySummary: severitySummary,
		AnalyzedAt:      time.Now(),
		SourceHash:      sourceHash,
	}
}

func (s *StaticAnalysisService) isContainerRunning() bool {
	cmd := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", s.containerName)
	output, err := cmd.Output()
	if err != nil {
		s.logger.Errorf("Failed to check container status: %v", err)
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

func (s *StaticAnalysisService) runSlitherInContainer(analysisDir, contractFile, solcVersion string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), analysisTimeout)
	defer cancel()

	// Build Slither command with --solc-solcs-select flag
	slitherCmd := fmt.Sprintf(
		"cd %s && slither %s --solc-remaps '@openzeppelin=node_modules/@openzeppelin' --solc-solcs-select %s --json -",
		analysisDir,
		contractFile,
		solcVersion,
	)

	s.logger.Infof("Executing Slither command: %s", slitherCmd)

	// Execute in container
	cmd := exec.CommandContext(ctx, "docker", "exec", s.containerName, "bash", "-c", slitherCmd)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	output := stdout.String()
	errorOutput := stderr.String()

	s.logger.Infof("Slither stdout: %s", output)
	s.logger.Infof("Slither stderr: %s", errorOutput)

	if output == "" && errorOutput != "" {
		output = errorOutput
	}

	if err != nil {
		if output != "" && (strings.Contains(output, `"success": true`) || strings.Contains(output, `"results"`)) {
			s.logger.Infof("Slither completed with findings (exit code %v), parsing results", err)
			return output, nil
		}

		s.logger.Errorf("Slither command failed with error: %v, stderr: %s", err, errorOutput)
		return output, fmt.Errorf("slither execution failed: %v, stderr: %s", err, errorOutput)
	}

	return output, nil
}

var (
	pragmaVersionRegex = regexp.MustCompile(`pragma\s+solidity\s+[^\d]*(\d+\.\d+)`)

	versionMap = map[string]string{
		"0.8": "0.8.20",
		"0.7": "0.7.6",
		"0.6": "0.6.12",
		"0.5": "0.5.16",
	}
)

func (s *StaticAnalysisService) detectSolidityVersion(analysisDir, contractFile string) (string, error) {
	cmd := exec.Command("docker", "exec", s.containerName, "head", "-c", "1024", filepath.Join(analysisDir, contractFile))
	output, err := cmd.Output()
	if err != nil {
		return "0.8.20", nil
	}

	matches := pragmaVersionRegex.FindSubmatch(output)
	if len(matches) < 2 {
		s.logger.Debug("No pragma solidity found, using default version")
		return "0.8.20", nil
	}

	version := string(matches[1])
	majorMinor := version
	if len(version) > 3 {
		majorMinor = version[:3]
	}

	if ltsVersion, ok := versionMap[majorMinor]; ok {
		s.logger.Infof("Detected Solidity version: %s", ltsVersion)
		return ltsVersion, nil
	}

	// Check if exact version is installed (from Dockerfile)
	installedVersions := []string{"0.8.20", "0.8.19", "0.8.0", "0.7.6", "0.6.12", "0.5.16"}
	for _, installed := range installedVersions {
		if strings.HasPrefix(installed, version) {
			s.logger.Infof("Detected Solidity version: %s", installed)
			return installed, nil
		}
	}

	s.logger.Warnf("Unknown version %s, using default", version)
	return "0.8.20", nil
}

func (s *StaticAnalysisService) parseSlitherOutput(output string) []analysis.SlitherIssue {
	var slitherResult struct {
		Success bool `json:"success"`
		Results struct {
			Detectors []struct {
				Check       string `json:"check"`
				Impact      string `json:"impact"`
				Confidence  string `json:"confidence"`
				Description string `json:"description"`
				Elements    []struct {
					SourceMapping struct {
						Lines []int `json:"lines"`
					} `json:"source_mapping"`
				} `json:"elements"`
				Reference string `json:"first_markdown_element"`
			} `json:"detectors"`
		} `json:"results"`
	}

	if err := json.Unmarshal([]byte(output), &slitherResult); err != nil || !slitherResult.Success {
		s.logger.Warnf("Failed to parse Slither JSON output: %v", err)
		return []analysis.SlitherIssue{}
	}

	issues := make([]analysis.SlitherIssue, 0, len(slitherResult.Results.Detectors))

	for _, detector := range slitherResult.Results.Detectors {
		issue := analysis.SlitherIssue{
			Type:        "detector",
			Title:       s.formatTitle(detector.Check),
			Description: s.cleanDescription(detector.Description),
			Severity:    strings.ToUpper(detector.Impact),
			Confidence:  detector.Confidence,
			Reference:   detector.Reference,
			Location:    s.formatLocation(detector.Elements),
		}
		issues = append(issues, issue)
	}

	return issues
}

func (s *StaticAnalysisService) formatLocation(elements []struct {
	SourceMapping struct {
		Lines []int `json:"lines"`
	} `json:"source_mapping"`
}) string {
	if len(elements) == 0 || len(elements[0].SourceMapping.Lines) == 0 {
		return ""
	}

	lines := elements[0].SourceMapping.Lines
	if len(lines) == 1 {
		return fmt.Sprintf("Contract.sol:L%d", lines[0])
	}

	if lines[0] == lines[len(lines)-1] {
		return fmt.Sprintf("Contract.sol:L%d", lines[0])
	}

	return fmt.Sprintf("Contract.sol:L%d-L%d", lines[0], lines[len(lines)-1])
}

func (s *StaticAnalysisService) cleanDescription(description string) string {
	cleaned := strings.ReplaceAll(description, "\\n", " ")
	cleaned = strings.ReplaceAll(cleaned, "\\t", " ")
	cleaned = strings.Join(strings.Fields(cleaned), " ")

	if len(cleaned) > maxDescriptionLength {
		cleaned = cleaned[:maxDescriptionLength-3] + "..."
	}

	return cleaned
}

func (s *StaticAnalysisService) formatTitle(check string) string {
	words := strings.Split(check, "-")
	for i, word := range words {
		words[i] = cases.Title(language.English).String(word)
	}
	return strings.Join(words, " ")
}

func (s *StaticAnalysisService) calculateSeveritySummary(issues []analysis.SlitherIssue) analysis.SeveritySummary {
	summary := analysis.SeveritySummary{}

	for _, issue := range issues {
		switch strings.ToUpper(issue.Severity) {
		case "HIGH":
			summary.High++
		case "MEDIUM":
			summary.Medium++
		case "LOW":
			summary.Low++
		default:
			summary.Informational++
		}
	}

	return summary
}

func (s *StaticAnalysisService) GetContainerStatus() (map[string]interface{}, error) {
	status := make(map[string]interface{})

	cmd := exec.Command("docker", "inspect", s.containerName)
	output, err := cmd.Output()
	if err != nil {
		status["exists"] = false
		status["running"] = false
		status["error"] = err.Error()
		return status, err
	}

	var containerInfo []map[string]interface{}
	if err := json.Unmarshal(output, &containerInfo); err != nil {
		return status, err
	}

	if len(containerInfo) > 0 {
		state := containerInfo[0]["State"].(map[string]interface{})
		status["exists"] = true
		status["running"] = state["Running"]
		status["status"] = state["Status"]
		status["started_at"] = state["StartedAt"]
	}

	return status, nil
}

func (s *StaticAnalysisService) cleanupContainerWorkspace(path string) {
	// Clean all folder save contract file, keep node_modules, package.json, package-lock.json
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "exec", s.containerName, "rm", "-rf", path)
	cmd.Run()
}
