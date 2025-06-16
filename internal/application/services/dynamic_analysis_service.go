package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	config "github.com/Azzurriii/slythr-go-backend/config"
	"github.com/Azzurriii/slythr-go-backend/internal/application/dto/dynamic_analysis"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/repository"
	"github.com/Azzurriii/slythr-go-backend/internal/infrastructure/cache"
	"github.com/Azzurriii/slythr-go-backend/internal/infrastructure/external"
	"github.com/Azzurriii/slythr-go-backend/pkg/utils"
	"github.com/redis/go-redis/v9"
)

const (
	DefaultAnalysisTimeout = 3 * time.Minute
	MaxSourceCodeSize      = 1024 * 1024 // 1MB
)

type DynamicAnalysisService struct {
	geminiClient *external.GeminiClient
	cache        *cache.AnalysisCache
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
		cache:        cache.NewAnalysisCache(redisClient, dynamicAnalysisRepo),
	}, nil
}

func (s *DynamicAnalysisService) AnalyzeContract(ctx context.Context, source string) (*dynamic_analysis.AnalyzeResponse, error) {
	if strings.TrimSpace(source) == "" {
		return nil, fmt.Errorf("source code cannot be empty")
	}

	sourceHash := utils.GenerateSourceHash(source)

	if cached, err := s.cache.Get(ctx, sourceHash); err == nil && cached != nil {
		return cached, nil
	}

	analysis, err := s.geminiClient.AnalyzeSmartContract(ctx, source)
	if err != nil {
		return &dynamic_analysis.AnalyzeResponse{
			Success: false,
			Analysis: dynamic_analysis.LLMAnalysis{
				Summary: "Analysis failed: " + err.Error(),
			},
		}, err
	}

	response := &dynamic_analysis.AnalyzeResponse{
		Success:     analysis.Success,
		Analysis:    s.convertAnalysis(analysis.Analysis),
		TotalIssues: len(analysis.Analysis.Vulnerabilities),
		AnalyzedAt:  time.Now(),
	}

	go s.cache.Set(context.Background(), sourceHash, response)

	return response, nil
}

func (s *DynamicAnalysisService) convertAnalysis(analysis external.SecurityAssessment) dynamic_analysis.LLMAnalysis {
	vulnerabilities := make([]dynamic_analysis.Vulnerability, len(analysis.Vulnerabilities))
	for i, vuln := range analysis.Vulnerabilities {
		vulnerabilities[i] = dynamic_analysis.Vulnerability{
			Title:          vuln.Title,
			Severity:       string(vuln.Severity),
			Description:    vuln.Description,
			Location:       vuln.Location,
			Recommendation: vuln.Recommendation,
		}
	}

	return dynamic_analysis.LLMAnalysis{
		SecurityScore:   analysis.SecurityScore,
		RiskLevel:       string(analysis.RiskLevel),
		Summary:         analysis.Summary,
		Vulnerabilities: vulnerabilities,
	}
}
