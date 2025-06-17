package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azzurriii/slythr-go-backend/internal/application/dto/analysis"
	"github.com/Azzurriii/slythr-go-backend/internal/application/dto/contracts"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/entities"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/repository"
	"github.com/redis/go-redis/v9"
)

const (
	defaultCacheTTL = 30 * time.Minute

	// Cache key prefixes
	contractPrefix        = "contract"
	dynamicAnalysisPrefix = "dynamic_analysis"
	staticAnalysisPrefix  = "static_analysis"
)

// Cache provides caching for all entities with Redis L1 and Database L2
type Cache struct {
	redis               *redis.Client
	contractRepo        repository.ContractRepository
	dynamicAnalysisRepo repository.DynamicAnalysisRepository
	staticAnalysisRepo  repository.StaticAnalysisRepository
}

// NewCache creates a new unified cache instance
func NewCache(
	redis *redis.Client,
	contractRepo repository.ContractRepository,
	dynamicAnalysisRepo repository.DynamicAnalysisRepository,
	staticAnalysisRepo repository.StaticAnalysisRepository,
) *Cache {
	return &Cache{
		redis:               redis,
		contractRepo:        contractRepo,
		dynamicAnalysisRepo: dynamicAnalysisRepo,
		staticAnalysisRepo:  staticAnalysisRepo,
	}
}

// Contract Cache Methods

func (c *Cache) GetContract(ctx context.Context, address, network string) (*contracts.ContractResponse, error) {
	// L1: Try Redis first
	if c.redis != nil {
		key := c.buildKey(contractPrefix, address, network)
		if data, err := c.redis.Get(ctx, key).Result(); err == nil {
			var result contracts.ContractResponse
			if json.Unmarshal([]byte(data), &result) == nil {
				return &result, nil
			}
		}
	}

	// L2: Try Database
	if c.contractRepo == nil {
		return nil, nil
	}

	contract, err := c.contractRepo.FindByAddressAndNetwork(ctx, address, network)
	if err != nil {
		return nil, nil
	}

	result := &contracts.ContractResponse{
		Address:         contract.Address,
		Network:         contract.Network,
		SourceCode:      contract.SourceCode,
		ContractName:    contract.ContractName,
		CompilerVersion: contract.CompilerVersion,
		SourceHash:      contract.SourceHash,
		CreatedAt:       contract.CreatedAt,
		UpdatedAt:       contract.UpdatedAt,
	}

	// Warm Redis cache
	c.setRedisAsync(c.buildKey(contractPrefix, address, network), result, defaultCacheTTL)

	return result, nil
}

func (c *Cache) SetContract(ctx context.Context, contract *entities.Contract) error {
	// Save to Database (L2)
	if c.contractRepo != nil {
		if err := c.contractRepo.SaveContract(ctx, contract); err != nil {
			// Log error but continue to cache in Redis
		}
	}

	// Save to Redis (L1)
	result := &contracts.ContractResponse{
		Address:         contract.Address,
		Network:         contract.Network,
		SourceCode:      contract.SourceCode,
		ContractName:    contract.ContractName,
		CompilerVersion: contract.CompilerVersion,
		SourceHash:      contract.SourceHash,
		CreatedAt:       contract.CreatedAt,
		UpdatedAt:       contract.UpdatedAt,
	}

	c.setRedisAsync(c.buildKey(contractPrefix, contract.Address, contract.Network), result, defaultCacheTTL)
	return nil
}

// Dynamic Analysis Cache Methods

func (c *Cache) GetDynamicAnalysis(ctx context.Context, sourceHash string) (*analysis.DynamicAnalysisResponse, error) {
	// L1: Try Redis first
	if c.redis != nil {
		key := c.buildKey(dynamicAnalysisPrefix, sourceHash)
		if data, err := c.redis.Get(ctx, key).Result(); err == nil {
			var result analysis.DynamicAnalysisResponse
			if json.Unmarshal([]byte(data), &result) == nil {
				return &result, nil
			}
		}
	}

	// L2: Try Database
	if c.dynamicAnalysisRepo == nil {
		return nil, nil
	}

	exists, err := c.dynamicAnalysisRepo.ExistsBySourceHash(ctx, sourceHash)
	if err != nil || !exists {
		return nil, nil
	}

	dbAnalysis, err := c.dynamicAnalysisRepo.FindBySourceHash(ctx, sourceHash)
	if err != nil {
		return nil, nil
	}

	var result analysis.DynamicAnalysisResponse
	if json.Unmarshal([]byte(dbAnalysis.LLMResponse), &result) != nil {
		return nil, nil
	}

	// Warm Redis cache
	c.setRedisAsync(c.buildKey(dynamicAnalysisPrefix, sourceHash), result, defaultCacheTTL)

	return &result, nil
}

func (c *Cache) SetDynamicAnalysis(ctx context.Context, sourceHash string, result *analysis.DynamicAnalysisResponse) error {
	// Save to Database (L2)
	if c.dynamicAnalysisRepo != nil {
		data, err := json.Marshal(result)
		if err == nil {
			if analysis, err := entities.NewDynamicAnalysis(sourceHash, string(data)); err == nil {
				c.dynamicAnalysisRepo.SaveAnalysis(ctx, analysis)
			}
		}
	}

	// Save to Redis (L1)
	c.setRedisAsync(c.buildKey(dynamicAnalysisPrefix, sourceHash), result, defaultCacheTTL)
	return nil
}

// Static Analysis Cache Methods

func (c *Cache) GetStaticAnalysis(ctx context.Context, sourceHash string) (*analysis.StaticAnalysisResponse, error) {
	// L1: Try Redis first
	if c.redis != nil {
		key := c.buildKey(staticAnalysisPrefix, sourceHash)
		if data, err := c.redis.Get(ctx, key).Result(); err == nil {
			var result analysis.StaticAnalysisResponse
			if json.Unmarshal([]byte(data), &result) == nil {
				return &result, nil
			}
		}
	}

	// L2: Try Database
	if c.staticAnalysisRepo == nil {
		return nil, nil
	}

	exists, err := c.staticAnalysisRepo.ExistsBySourceHash(ctx, sourceHash)
	if err != nil || !exists {
		return nil, nil
	}

	dbAnalysis, err := c.staticAnalysisRepo.FindBySourceHash(ctx, sourceHash)
	if err != nil {
		return nil, nil
	}

	var result analysis.StaticAnalysisResponse
	if json.Unmarshal([]byte(dbAnalysis.SlitherOutput), &result) != nil {
		return nil, nil
	}

	// Warm Redis cache
	c.setRedisAsync(c.buildKey(staticAnalysisPrefix, sourceHash), result, defaultCacheTTL)

	return &result, nil
}

func (c *Cache) SetStaticAnalysis(ctx context.Context, sourceHash string, result *analysis.StaticAnalysisResponse) error {
	// Save to Database (L2)
	if c.staticAnalysisRepo != nil {
		data, err := json.Marshal(result)
		if err == nil {
			if analysis, err := entities.NewStaticAnalysis(sourceHash, string(data)); err == nil {
				c.staticAnalysisRepo.SaveAnalysis(ctx, analysis)
			}
		}
	}

	// Save to Redis (L1)
	c.setRedisAsync(c.buildKey(staticAnalysisPrefix, sourceHash), result, defaultCacheTTL)
	return nil
}

// Helper Methods

func (c *Cache) buildKey(prefix string, parts ...string) string {
	key := prefix
	for _, part := range parts {
		key += ":" + part
	}
	return key
}

func (c *Cache) setRedisAsync(key string, data interface{}, ttl time.Duration) {
	if c.redis == nil {
		return
	}

	go func() {
		if jsonData, err := json.Marshal(data); err == nil {
			c.redis.Set(context.Background(), key, jsonData, ttl)
		}
	}()
}

// Invalidation Methods

func (c *Cache) InvalidateContract(address, network string) error {
	if c.redis != nil {
		key := c.buildKey(contractPrefix, address, network)
		return c.redis.Del(context.Background(), key).Err()
	}
	return nil
}

func (c *Cache) InvalidateDynamicAnalysis(sourceHash string) error {
	if c.redis != nil {
		key := c.buildKey(dynamicAnalysisPrefix, sourceHash)
		return c.redis.Del(context.Background(), key).Err()
	}
	return nil
}

func (c *Cache) InvalidateStaticAnalysis(sourceHash string) error {
	if c.redis != nil {
		key := c.buildKey(staticAnalysisPrefix, sourceHash)
		return c.redis.Del(context.Background(), key).Err()
	}
	return nil
}

// Bulk Operations

func (c *Cache) InvalidateAll() error {
	if c.redis != nil {
		patterns := []string{
			fmt.Sprintf("%s:*", contractPrefix),
			fmt.Sprintf("%s:*", dynamicAnalysisPrefix),
			fmt.Sprintf("%s:*", staticAnalysisPrefix),
		}

		for _, pattern := range patterns {
			keys, err := c.redis.Keys(context.Background(), pattern).Result()
			if err != nil {
				continue
			}
			if len(keys) > 0 {
				c.redis.Del(context.Background(), keys...)
			}
		}
	}
	return nil
}

// Health Check

func (c *Cache) HealthCheck(ctx context.Context) error {
	if c.redis != nil {
		return c.redis.Ping(ctx).Err()
	}
	return nil
}
