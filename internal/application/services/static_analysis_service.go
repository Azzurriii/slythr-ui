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

	"github.com/Azzurriii/slythr-go-backend/internal/application/dto/static_analysis"
	"github.com/google/uuid"
)

type StaticAnalysisService struct {
	slitherContainer string
	workspacePath    string
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
	}
}

func (s *StaticAnalysisService) AnalyzeContract(source string) (*static_analysis.AnalyzeResponse, error) {
	if !s.isContainerRunning() {
		return &static_analysis.AnalyzeResponse{
			Success: false,
			Message: "Slither analysis container is not running. Please start the container first.",
		}, fmt.Errorf("slither container not running")
	}

	analysisID := uuid.New().String()
	tempDir := filepath.Join(os.TempDir(), analysisID)

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return &static_analysis.AnalyzeResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create temp directory: %v", err),
		}, err
	}
	defer os.RemoveAll(tempDir)

	contractFile := filepath.Join(tempDir, "Contract.sol")
	if err := os.WriteFile(contractFile, []byte(source), 0644); err != nil {
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
		return &static_analysis.AnalyzeResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to write package.json: %v", err),
		}, err
	}

	copyCmd := exec.Command("docker", "cp", tempDir, fmt.Sprintf("%s:%s", s.slitherContainer, s.workspacePath))
	if err := copyCmd.Run(); err != nil {
		return &static_analysis.AnalyzeResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to copy files to container: %v", err),
		}, err
	}

	containerPath := filepath.Join(s.workspacePath, analysisID)
	installCmd := exec.Command("docker", "exec", s.slitherContainer, "bash", "-c",
		fmt.Sprintf("cd %s && npm install", containerPath))
	if err := installCmd.Run(); err != nil {
		return &static_analysis.AnalyzeResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to install dependencies: %v", err),
		}, err
	}

	slitherOutput, err := s.runSlitherInContainer(containerPath, "Contract.sol")
	if err != nil {
		return &static_analysis.AnalyzeResponse{
			Success:   false,
			Message:   fmt.Sprintf("Failed to run Slither: %v", err),
			RawOutput: slitherOutput,
		}, err
	}

	issues := s.parseSlitherOutput(slitherOutput)

	return &static_analysis.AnalyzeResponse{
		Success:     true,
		Issues:      issues,
		TotalIssues: len(issues),
		AnalyzedAt:  time.Now(),
		RawOutput:   slitherOutput,
	}, nil
}

func (s *StaticAnalysisService) isContainerRunning() bool {
	cmd := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", s.slitherContainer)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

func (s *StaticAnalysisService) runSlitherInContainer(analysisDir, contractFile string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Build Slither command to run in container
	slitherCmd := fmt.Sprintf(
		"cd %s && slither %s --solc-remaps '@openzeppelin=node_modules/@openzeppelin' --json -",
		analysisDir,
		contractFile,
	)

	// Execute in container
	cmd := exec.CommandContext(ctx, "docker", "exec", s.slitherContainer, "bash", "-c", slitherCmd)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	output := stdout.String()
	if output == "" {
		output = stderr.String()
	}

	if err != nil && output == "" {
		return "", fmt.Errorf("slither execution failed: %v", err)
	}

	return output, nil
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
				Description: detector.Description,
				Severity:    strings.ToUpper(detector.Impact),
				Confidence:  detector.Confidence,
				Reference:   detector.Reference,
			}

			if len(detector.Elements) > 0 {
				element := detector.Elements[0]
				if len(element.SourceMapping.Lines) > 0 {
					issue.Location = fmt.Sprintf("%s#L%d",
						element.SourceMapping.Filename,
						element.SourceMapping.Lines[0])
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
		words[i] = strings.Title(word)
	}
	return strings.Join(words, " ")
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
