package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Azzurriii/slythr-go-backend/internal/application/dto/dynamic_analysis"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/entities"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/repository"
	"github.com/redis/go-redis/v9"
)

type AnalysisCache struct {
	redis *redis.Client
	repo  repository.DynamicAnalysisRepository
}

func NewAnalysisCache(redis *redis.Client, repo repository.DynamicAnalysisRepository) *AnalysisCache {
	return &AnalysisCache{
		redis: redis,
		repo:  repo,
	}
}

func (c *AnalysisCache) Get(ctx context.Context, sourceHash string) (*dynamic_analysis.AnalyzeResponse, error) {
	// L1: Try Redis first
	if c.redis != nil {
		key := "analysis:" + sourceHash
		data, err := c.redis.Get(ctx, key).Result()
		if err == nil {
			var result dynamic_analysis.AnalyzeResponse
			if json.Unmarshal([]byte(data), &result) == nil {
				return &result, nil
			}
		}
	}

	// L2: Try Database
	exists, err := c.repo.ExistsBySourceHash(ctx, sourceHash)
	if err != nil || !exists {
		return nil, nil
	}

	dbAnalysis, err := c.repo.FindBySourceHash(ctx, sourceHash)
	if err != nil {
		return nil, nil
	}

	var result dynamic_analysis.AnalyzeResponse
	if json.Unmarshal([]byte(dbAnalysis.LLMResponse), &result) != nil {
		return nil, nil
	}

	// Warm Redis cache
	if c.redis != nil {
		go func() {
			key := "analysis:" + sourceHash
			data, _ := json.Marshal(result)
			c.redis.Set(context.Background(), key, data, 15*time.Minute)
		}()
	}

	return &result, nil
}

func (c *AnalysisCache) Set(ctx context.Context, sourceHash string, result *dynamic_analysis.AnalyzeResponse) error {
	// Save to Database (L2)
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	analysis, err := entities.NewDynamicAnalysis(sourceHash, string(data))
	if err != nil {
		return err
	}

	if err := c.repo.SaveAnalysis(ctx, analysis); err != nil {
		// Ignore duplicate errors
		return nil
	}

	// Save to Redis (L1)
	if c.redis != nil {
		go func() {
			key := "analysis:" + sourceHash
			c.redis.Set(context.Background(), key, data, 15*time.Minute)
		}()
	}

	return nil
}
