package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azzurriii/slythr/internal/application/dto/analysis"
	"github.com/Azzurriii/slythr/internal/application/dto/contracts"
	"github.com/Azzurriii/slythr/internal/domain/entities"
	"github.com/Azzurriii/slythr/internal/domain/repository"
	"github.com/redis/go-redis/v9"
)

const (
	defaultCacheTTL = 30 * time.Minute

	// Cache key prefixes
	contractPrefix        = "contract"
	dynamicAnalysisPrefix = "dynamic_analysis"
	staticAnalysisPrefix  = "static_analysis"
)

/*
Cache provides caching for all entities with Redis and Database
Redis is used as a cache layer and Database is used as a persistence layer
*/
type Cache struct {
	redis               *redis.Client
	contractRepo        repository.ContractRepository
	dynamicAnalysisRepo repository.DynamicAnalysisRepository
	staticAnalysisRepo  repository.StaticAnalysisRepository
}

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

func (c *Cache) GetContract(ctx context.Context, address, network string) (*contracts.ContractResponse, error) {
	if c.redis != nil {
		key := c.buildKey(contractPrefix, address, network)
		if data, err := c.redis.Get(ctx, key).Result(); err == nil {
			var result contracts.ContractResponse
			if json.Unmarshal([]byte(data), &result) == nil {
				return &result, nil
			}
		}
	}

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

	c.setRedisAsync(c.buildKey(contractPrefix, address, network), result, defaultCacheTTL)

	return result, nil
}

func (c *Cache) SetContract(ctx context.Context, contract *entities.Contract) error {
	if c.contractRepo != nil {
		if err := c.contractRepo.SaveContract(ctx, contract); err != nil {
		}
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

	c.setRedisAsync(c.buildKey(contractPrefix, contract.Address, contract.Network), result, defaultCacheTTL)
	return nil
}

func (c *Cache) GetDynamicAnalysis(ctx context.Context, sourceHash string) (*analysis.DynamicAnalysisResponse, error) {
	if c.redis != nil {
		key := c.buildKey(dynamicAnalysisPrefix, sourceHash)
		if data, err := c.redis.Get(ctx, key).Result(); err == nil {
			var result analysis.DynamicAnalysisResponse
			if json.Unmarshal([]byte(data), &result) == nil {
				return &result, nil
			}
		}
	}

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

	c.setRedisAsync(c.buildKey(dynamicAnalysisPrefix, sourceHash), result, defaultCacheTTL)

	return &result, nil
}

func (c *Cache) SetDynamicAnalysis(ctx context.Context, sourceHash string, result *analysis.DynamicAnalysisResponse) error {
	if c.dynamicAnalysisRepo != nil {
		data, err := json.Marshal(result)
		if err == nil {
			if analysis, err := entities.NewDynamicAnalysis(sourceHash, string(data)); err == nil {
				c.dynamicAnalysisRepo.SaveAnalysis(ctx, analysis)
			}
		}
	}

	c.setRedisAsync(c.buildKey(dynamicAnalysisPrefix, sourceHash), result, defaultCacheTTL)
	return nil
}

func (c *Cache) GetStaticAnalysis(ctx context.Context, sourceHash string) (*analysis.StaticAnalysisResponse, error) {
	if c.redis != nil {
		key := c.buildKey(staticAnalysisPrefix, sourceHash)
		if data, err := c.redis.Get(ctx, key).Result(); err == nil {
			var result analysis.StaticAnalysisResponse
			if json.Unmarshal([]byte(data), &result) == nil {
				return &result, nil
			}
		}
	}

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

	c.setRedisAsync(c.buildKey(staticAnalysisPrefix, sourceHash), result, defaultCacheTTL)

	return &result, nil
}

func (c *Cache) SetStaticAnalysis(ctx context.Context, sourceHash string, result *analysis.StaticAnalysisResponse) error {
	if c.staticAnalysisRepo != nil {
		data, err := json.Marshal(result)
		if err == nil {
			if analysis, err := entities.NewStaticAnalysis(sourceHash, string(data)); err == nil {
				c.staticAnalysisRepo.SaveAnalysis(ctx, analysis)
			}
		}
	}

	c.setRedisAsync(c.buildKey(staticAnalysisPrefix, sourceHash), result, defaultCacheTTL)
	return nil
}

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

func (c *Cache) HealthCheck(ctx context.Context) error {
	if c.redis != nil {
		return c.redis.Ping(ctx).Err()
	}
	return nil
}
