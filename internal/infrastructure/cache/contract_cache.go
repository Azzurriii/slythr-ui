package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Azzurriii/slythr-go-backend/internal/application/dto/contracts"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/entities"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/repository"
	"github.com/redis/go-redis/v9"
)

type ContractCache struct {
	redis *redis.Client
	repo  repository.ContractRepository
}

func NewContractCache(redis *redis.Client, repo repository.ContractRepository) *ContractCache {
	return &ContractCache{
		redis: redis,
		repo:  repo,
	}
}

func (c *ContractCache) Get(ctx context.Context, address, network string) (*contracts.ContractResponse, error) {
	// L1: Try Redis first
	if c.redis != nil {
		key := "contract:" + address + ":" + network
		data, err := c.redis.Get(ctx, key).Result()
		if err == nil {
			var result contracts.ContractResponse
			if json.Unmarshal([]byte(data), &result) == nil {
				return &result, nil
			}
		}
	}

	// L2: Try Database
	contract, err := c.repo.FindByAddressAndNetwork(ctx, address, network)
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
	if c.redis != nil {
		go func() {
			key := "contract:" + address + ":" + network
			data, _ := json.Marshal(result)
			c.redis.Set(context.Background(), key, data, 15*time.Minute)
		}()
	}

	return result, nil
}

func (c *ContractCache) Set(ctx context.Context, contract *entities.Contract) error {
	// Save to Database (L2)
	if err := c.repo.SaveContract(ctx, contract); err != nil {
		// Ignore duplicate errors
		return nil
	}

	// Save to Redis (L1)
	if c.redis != nil {
		go func() {
			key := "contract:" + contract.Address + ":" + contract.Network
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
			data, _ := json.Marshal(result)
			c.redis.Set(context.Background(), key, data, 15*time.Minute)
		}()
	}

	return nil
}
