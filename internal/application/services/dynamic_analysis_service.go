package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	config "github.com/Azzurriii/slythr/config"
	"github.com/Azzurriii/slythr/internal/application/dto/analysis"
	"github.com/Azzurriii/slythr/internal/domain/repository"
	"github.com/Azzurriii/slythr/internal/infrastructure/cache"
	"github.com/Azzurriii/slythr/internal/infrastructure/external"
	"github.com/Azzurriii/slythr/pkg/utils"
	"github.com/redis/go-redis/v9"
)

const (
	DefaultAnalysisTimeout = 3 * time.Minute
	MaxSourceCodeSize      = 1024 * 1024 // 1MB
)

type DynamicAnalysisService struct {
	geminiClient *external.GeminiClient
	cache        *cache.Cache
}

type ServiceOptions struct {
	AnalysisTimeout time.Duration
	GeminiOptions   *external.GeminiClientOptions
}

func NewDynamicAnalysisService(
	dynamicAnalysisRepo repository.DynamicAnalysisRepository,
	contractRepo repository.ContractRepository,
	opts *ServiceOptions,
) (*DynamicAnalysisService, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	var redisClient *redis.Client
	if cfg.Redis.Addr != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
	}

	return &DynamicAnalysisService{
		geminiClient: external.NewGeminiClient(cfg.Gemini, nil),
		cache:        cache.NewCache(redisClient, contractRepo, dynamicAnalysisRepo, nil),
	}, nil
}

func (s *DynamicAnalysisService) AnalyzeContract(ctx context.Context, source string) (*analysis.DynamicAnalysisResponse, error) {
	if strings.TrimSpace(source) == "" {
		return nil, fmt.Errorf("source code cannot be empty")
	}

	sourceHash := utils.GenerateSourceHash(source)

	if cached, err := s.cache.GetDynamicAnalysis(ctx, sourceHash); err == nil && cached != nil {
		return cached, nil
	}

	geminiResult, err := s.geminiClient.AnalyzeSmartContract(ctx, source)
	if err != nil {
		return &analysis.DynamicAnalysisResponse{
			Success: false,
			Analysis: analysis.LLMAnalysis{
				Summary: "Analysis failed: " + err.Error(),
			},
		}, err
	}

	response := &analysis.DynamicAnalysisResponse{
		Success:     geminiResult.Success,
		Analysis:    s.convertAnalysis(geminiResult.Analysis),
		TotalIssues: len(geminiResult.Analysis.Vulnerabilities),
		AnalyzedAt:  time.Now(),
	}

	go s.cache.SetDynamicAnalysis(context.Background(), sourceHash, response)

	return response, nil
}

func (s *DynamicAnalysisService) convertAnalysis(securityAssessment external.SecurityAssessment) analysis.LLMAnalysis {
	vulnerabilities := make([]analysis.Vulnerability, len(securityAssessment.Vulnerabilities))
	for i, vuln := range securityAssessment.Vulnerabilities {
		vulnerabilities[i] = analysis.Vulnerability{
			Title:          vuln.Title,
			Severity:       string(vuln.Severity),
			Description:    vuln.Description,
			Location:       vuln.Location,
			Recommendation: vuln.Recommendation,
		}
	}

	return analysis.LLMAnalysis{
		SecurityScore:   securityAssessment.SecurityScore,
		RiskLevel:       string(securityAssessment.RiskLevel),
		Summary:         securityAssessment.Summary,
		Vulnerabilities: vulnerabilities,
	}
}
