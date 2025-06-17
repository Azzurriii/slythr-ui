package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/Azzurriii/slythr-go-backend/internal/application/dto/static_analysis"
	"github.com/Azzurriii/slythr-go-backend/pkg/logger"
	"github.com/google/uuid"
)

type StaticAnalysisService struct {
	slitherContainer string
	workspacePath    string
	logger           *logger.Logger
}

func NewStaticAnalysisService() *StaticAnalysisService {
	containerName := os.Getenv("SLITHER_CONTAINER_NAME")
	if containerName == "" {
		containerName = "slither"
	}

	workspacePath := os.Getenv("WORKSPACE_PATH")
	if workspacePath == "" {
		workspacePath = "/workspace"
	}

	return &StaticAnalysisService{
		slitherContainer: containerName,
		workspacePath:    workspacePath,
		logger:           logger.Default,
	}
}

func (s *StaticAnalysisService) AnalyzeContract(source string) (*static_analysis.AnalyzeResponse, error) {
	if !s.isContainerRunning() {
		s.logger.Warnf("Slither analysis container is not running")
		return &static_analysis.AnalyzeResponse{
			Success: false,
			Message: "Slither analysis container is not running. Please start the container first.",
		}, fmt.Errorf("slither container not running")
	}

	analysisID := uuid.New().String()
	s.logger.Infof("Starting static analysis with ID: %s", analysisID)

	tempDir := filepath.Join(os.TempDir(), analysisID)

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		s.logger.Errorf("Failed to create temp directory: %v", err)
		return &static_analysis.AnalyzeResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create temp directory: %v", err),
		}, err
	}
	defer os.RemoveAll(tempDir)

	contractFile := filepath.Join(tempDir, "Contract.sol")
	if err := os.WriteFile(contractFile, []byte(source), 0644); err != nil {
		s.logger.Errorf("Failed to write contract file: %v", err)
		return &static_analysis.AnalyzeResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to write contract file: %v", err),
		}, err
	}

	packageJSON := `{
		"name": "slither-analysis",
		"version": "1.0.0",
		"dependencies": {
			"@openzeppelin/contracts": "^4.9.0",
			"@openzeppelin/contracts-upgradeable": "^4.9.0",
			"@chainlink/contracts": "^0.6.1",
			"@uniswap/v2-core": "^1.0.1",
			"@uniswap/v3-core": "^1.0.0"
		}
	}`
	if err := os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0644); err != nil {
		s.logger.Errorf("Failed to write package.json: %v", err)
		return &static_analysis.AnalyzeResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to write package.json: %v", err),
		}, err
	}

	copyCmd := exec.Command("docker", "cp", tempDir, fmt.Sprintf("%s:%s", s.slitherContainer, s.workspacePath))
	if err := copyCmd.Run(); err != nil {
		s.logger.Errorf("Failed to copy files to container: %v", err)
		return &static_analysis.AnalyzeResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to copy files to container: %v", err),
		}, err
	}

	containerPath := filepath.Join(s.workspacePath, analysisID)
	s.logger.Infof("Installing dependencies in container path: %s", containerPath)

	installCmd := exec.Command("docker", "exec", s.slitherContainer, "bash", "-c",
		fmt.Sprintf("cd %s && npm install", containerPath))
	if err := installCmd.Run(); err != nil {
		s.logger.Errorf("Failed to install dependencies: %v", err)
		return &static_analysis.AnalyzeResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to install dependencies: %v", err),
		}, err
	}

	s.logger.Infof("Running Slither analysis on contract")
	slitherOutput, err := s.runSlitherInContainer(containerPath, "Contract.sol")
	if err != nil {
		s.logger.Errorf("Failed to run Slither: %v", err)
		return &static_analysis.AnalyzeResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to run Slither: %v", err),
		}, err
	}

	issues := s.parseSlitherOutput(slitherOutput)
	s.logger.Infof("Static analysis completed successfully, found %d issues", len(issues))

	severitySummary := s.calculateSeveritySummary(issues)

	response := &static_analysis.AnalyzeResponse{
		Success:         true,
		Issues:          issues,
		TotalIssues:     len(issues),
		SeveritySummary: severitySummary,
		AnalyzedAt:      time.Now(),
	}

	return response, nil
}

func (s *StaticAnalysisService) isContainerRunning() bool {
	cmd := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", s.slitherContainer)
	output, err := cmd.Output()
	if err != nil {
		s.logger.Errorf("Failed to check container status: %v", err)
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

func (s *StaticAnalysisService) runSlitherInContainer(analysisDir, contractFile string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Detect Solidity version từ contract source
	solcVersion, err := s.detectSolidityVersion(analysisDir, contractFile)
	if err != nil {
		s.logger.Warnf("Failed to detect Solidity version, using default: %v", err)
		solcVersion = "0.8.20" // fallback
	}

	// Build Slither command với --solc-solcs-select flag
	slitherCmd := fmt.Sprintf(
		"cd %s && slither %s --solc-remaps '@openzeppelin=node_modules/@openzeppelin' --solc-solcs-select %s --json -",
		analysisDir,
		contractFile,
		solcVersion,
	)

	s.logger.Infof("Executing Slither command: %s", slitherCmd)

	// Execute in container
	cmd := exec.CommandContext(ctx, "docker", "exec", s.slitherContainer, "bash", "-c", slitherCmd)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	output := stdout.String()
	errorOutput := stderr.String()

	s.logger.Infof("Slither stdout: %s", output)
	s.logger.Infof("Slither stderr: %s", errorOutput)

	// Return stderr nếu stdout trống
	if output == "" && errorOutput != "" {
		output = errorOutput
	}

	// Slither trả về exit code 255 khi có findings - đây KHÔNG phải lỗi
	// Chỉ return error nếu không có output hoặc có stderr mà không có stdout
	if err != nil {
		// Kiểm tra nếu có JSON output hợp lệ thì coi như success
		if output != "" && (strings.Contains(output, `"success": true`) || strings.Contains(output, `"results"`)) {
			s.logger.Infof("Slither completed with findings (exit code %v), parsing results", err)
			return output, nil
		}

		s.logger.Errorf("Slither command failed with error: %v, stderr: %s", err, errorOutput)
		return output, fmt.Errorf("slither execution failed: %v, stderr: %s", err, errorOutput)
	}

	return output, nil
}

func (s *StaticAnalysisService) detectSolidityVersion(analysisDir, contractFile string) (string, error) {
	readCmd := exec.Command("docker", "exec", s.slitherContainer, "cat", filepath.Join(analysisDir, contractFile))
	output, err := readCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to read contract file: %v", err)
	}

	content := string(output)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "pragma solidity") {
			version := s.extractVersionFromPragma(line)
			if version != "" {
				s.logger.Infof("Detected Solidity version: %s", version)
				return version, nil
			}
		}
	}

	return "", fmt.Errorf("no pragma solidity found")
}

func (s *StaticAnalysisService) extractVersionFromPragma(pragma string) string {
	pragma = strings.TrimPrefix(pragma, "pragma solidity")
	pragma = strings.TrimSpace(pragma)
	pragma = strings.TrimSuffix(pragma, ";")
	pragma = strings.TrimSpace(pragma)

	if strings.HasPrefix(pragma, "^") {
		if strings.HasPrefix(pragma, "^0.8") {
			return "0.8.20"
		} else if strings.HasPrefix(pragma, "^0.7") {
			return "0.7.6"
		} else if strings.HasPrefix(pragma, "^0.6") {
			return "0.6.12"
		} else if strings.HasPrefix(pragma, "^0.5") {
			return "0.5.16"
		}
	} else if strings.HasPrefix(pragma, ">=") {
		if strings.Contains(pragma, "0.8") {
			return "0.8.20"
		} else if strings.Contains(pragma, "0.7") {
			return "0.7.6"
		}
	} else {
		version := strings.Fields(pragma)[0]
		if len(version) > 0 && (version[0] >= '0' && version[0] <= '9') {
			return version
		}
	}

	return "0.8.20"
}

func (s *StaticAnalysisService) parseSlitherOutput(output string) []static_analysis.SlitherIssue {
	var issues []static_analysis.SlitherIssue

	var slitherResult struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		Results struct {
			Detectors []struct {
				Check       string `json:"check"`
				Impact      string `json:"impact"`
				Confidence  string `json:"confidence"`
				Description string `json:"description"`
				Elements    []struct {
					Name          string `json:"name"`
					SourceMapping struct {
						Filename string `json:"filename"`
						Lines    []int  `json:"lines"`
					} `json:"source_mapping"`
				} `json:"elements"`
				Reference string `json:"first_markdown_element"`
			} `json:"detectors"`
		} `json:"results"`
	}

	if err := json.Unmarshal([]byte(output), &slitherResult); err == nil && slitherResult.Success {
		for _, detector := range slitherResult.Results.Detectors {
			issue := static_analysis.SlitherIssue{
				Type:        "detector",
				Title:       s.formatTitle(detector.Check),
				Description: s.cleanDescription(detector.Description),
				Severity:    strings.ToUpper(detector.Impact),
				Confidence:  detector.Confidence,
				Reference:   detector.Reference,
			}

			// Format location đẹp hơn
			if len(detector.Elements) > 0 {
				element := detector.Elements[0]
				if len(element.SourceMapping.Lines) > 0 {
					line := element.SourceMapping.Lines[0]
					if len(element.SourceMapping.Lines) > 1 {
						lastLine := element.SourceMapping.Lines[len(element.SourceMapping.Lines)-1]
						if lastLine != line {
							issue.Location = fmt.Sprintf("Contract.sol:L%d-L%d", line, lastLine)
						} else {
							issue.Location = fmt.Sprintf("Contract.sol:L%d", line)
						}
					} else {
						issue.Location = fmt.Sprintf("Contract.sol:L%d", line)
					}
				}
			}

			issues = append(issues, issue)
		}
	} else {
		issues = s.parseTextOutput(output)
	}

	return issues
}

func (s *StaticAnalysisService) parseTextOutput(output string) []static_analysis.SlitherIssue {
	var issues []static_analysis.SlitherIssue
	lines := strings.Split(output, "\n")

	var currentIssue *static_analysis.SlitherIssue

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, "INFO:Detectors:") {
			if currentIssue != nil {
				issues = append(issues, *currentIssue)
			}

			description := strings.TrimPrefix(line, "INFO:Detectors:")
			currentIssue = &static_analysis.SlitherIssue{
				Type:        "detector",
				Description: strings.TrimSpace(description),
				Severity:    s.extractSeverity(line),
			}
			currentIssue.Title = s.extractTitle(currentIssue.Description)
		} else if strings.HasPrefix(line, "Reference:") && currentIssue != nil {
			currentIssue.Reference = strings.TrimSpace(strings.TrimPrefix(line, "Reference:"))
		}
	}

	if currentIssue != nil {
		issues = append(issues, *currentIssue)
	}

	return issues
}

func (s *StaticAnalysisService) formatTitle(check string) string {
	words := strings.Split(check, "-")
	for i, word := range words {
		words[i] = cases.Title(language.English).String(word)
	}
	return strings.Join(words, " ")
}

func (s *StaticAnalysisService) cleanDescription(description string) string {
	// Bỏ các ký tự xuống dòng và tab không cần thiết
	cleaned := strings.ReplaceAll(description, "\\n", " ")
	cleaned = strings.ReplaceAll(cleaned, "\\t", " ")

	// Bỏ khoảng trắng thừa
	cleaned = strings.Join(strings.Fields(cleaned), " ")

	// Giới hạn độ dài nếu quá dài
	if len(cleaned) > 200 {
		cleaned = cleaned[:197] + "..."
	}

	return strings.TrimSpace(cleaned)
}

func (s *StaticAnalysisService) extractTitle(description string) string {
	parts := strings.SplitN(description, ".", 2)
	if len(parts) > 0 {
		title := strings.TrimSpace(parts[0])
		if len(title) > 100 {
			title = title[:100] + "..."
		}
		return title
	}
	return "Security Issue"
}

func (s *StaticAnalysisService) extractSeverity(line string) string {
	line = strings.ToLower(line)
	if strings.Contains(line, "high") || strings.Contains(line, "reentrancy") {
		return "HIGH"
	}
	if strings.Contains(line, "medium") || strings.Contains(line, "timestamp") {
		return "MEDIUM"
	}
	if strings.Contains(line, "low") || strings.Contains(line, "optimization") {
		return "LOW"
	}
	return "INFO"
}

func (s *StaticAnalysisService) GetContainerStatus() (map[string]interface{}, error) {
	status := make(map[string]interface{})

	cmd := exec.Command("docker", "inspect", s.slitherContainer)
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

func (s *StaticAnalysisService) calculateSeveritySummary(issues []static_analysis.SlitherIssue) static_analysis.SeveritySummary {
	summary := static_analysis.SeveritySummary{}

	for _, issue := range issues {
		switch strings.ToUpper(issue.Severity) {
		case "HIGH":
			summary.High++
		case "MEDIUM":
			summary.Medium++
		case "LOW":
			summary.Low++
		case "INFORMATIONAL":
			summary.Informational++
		default:
			summary.Informational++ // fallback
		}
	}

	return summary
}
