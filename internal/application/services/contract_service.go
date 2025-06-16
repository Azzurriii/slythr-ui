package services

import (
	"context"
	"errors"
	"fmt"

	config "github.com/Azzurriii/slythr-go-backend/config"
	"github.com/Azzurriii/slythr-go-backend/internal/application/dto/contracts"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/entities"
	domainerrors "github.com/Azzurriii/slythr-go-backend/internal/domain/errors"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/repository"
	"github.com/Azzurriii/slythr-go-backend/internal/infrastructure/cache"
	"github.com/Azzurriii/slythr-go-backend/internal/infrastructure/external"
	"github.com/Azzurriii/slythr-go-backend/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type ContractService struct {
	client external.EtherscanService
	cache  *cache.ContractCache
}

func NewContractService(repo repository.ContractRepository, client external.EtherscanService) ContractServiceInterface {
	cfg, _ := config.LoadConfig()

	var redisClient *redis.Client
	if cfg.Redis.Addr != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
	}

	return &ContractService{
		client: client,
		cache:  cache.NewContractCache(redisClient, repo),
	}
}

func (s *ContractService) FetchContractSourceCode(ctx context.Context, req *contracts.GetContractSourceCodeRequest) (*contracts.GetContractSourceCodeResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Check cache (Redis L1 + Database L2)
	if cached, err := s.cache.Get(ctx, req.Address, req.Network); err == nil && cached != nil {
		return &contracts.GetContractSourceCodeResponse{
			Address:    cached.Address,
			SourceCode: cached.SourceCode,
			SourceHash: cached.SourceHash,
			Network:    cached.Network,
			CachedAt:   &cached.CreatedAt,
		}, nil
	}

	// Fetch from external API
	sourceCode, err := s.client.GetContractSourceCode(ctx, req.Address, req.Network)
	if err != nil {
		return nil, s.handleExternalError(err)
	}

	details, err := s.client.GetContractDetails(ctx, req.Address, req.Network)
	if err != nil {
		return nil, s.handleExternalError(err)
	}

	// Create and save contract
	sourceHash := utils.GenerateSourceHash(sourceCode)
	contract, err := entities.NewContract(req.Address, req.Network, sourceCode, details.ContractName, details.CompilerVersion, sourceHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract: %w", err)
	}

	// Save to cache
	go s.cache.Set(context.Background(), contract)

	return &contracts.GetContractSourceCodeResponse{
		Address:    contract.Address,
		SourceCode: contract.SourceCode,
		SourceHash: contract.SourceHash,
		Network:    contract.Network,
	}, nil
}

func (s *ContractService) GetContractByAddressAndNetwork(ctx context.Context, address, network string) (*contracts.ContractResponse, error) {
	contract, err := s.cache.Get(ctx, address, network)
	if err != nil || contract == nil {
		return nil, fmt.Errorf("contract not found")
	}
	return contract, nil
}

func (s *ContractService) handleExternalError(err error) error {
	var etherscanErr external.EtherscanError
	if errors.As(err, &etherscanErr) {
		return domainerrors.ErrContractNotFound
	}
	return fmt.Errorf("external service error: %w", err)
}
