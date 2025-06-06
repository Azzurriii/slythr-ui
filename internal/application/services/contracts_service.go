package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	config "github.com/Azzurriii/slythr-go-backend/config"
	"github.com/Azzurriii/slythr-go-backend/internal/application/dto/contracts"
	entity "github.com/Azzurriii/slythr-go-backend/internal/domain/entities"
	apperrors "github.com/Azzurriii/slythr-go-backend/internal/domain/errors"
	infrastructures "github.com/Azzurriii/slythr-go-backend/internal/infrastructure/external"
	"gorm.io/gorm"
)

const (
	NetworkEthereum  = "ethereum"
	NetworkPolygon   = "polygon"
	NetworkBSC       = "bsc"
	NetworkBase      = "base"
	NetworkArbitrum  = "arbitrum"
	NetworkAvalanche = "avalanche"
	NetworkOptimism  = "optimism"
	NetworkGnosis    = "gnosis"
	NetworkFantom    = "fantom"
	NetworkCelo      = "celo"
)

var (
	ErrContractNotFound = apperrors.ErrContractNotFound
	ErrInvalidAddress   = apperrors.ErrInvalidAddress
	ErrEmptySourceCode  = apperrors.ErrEmptySourceCode
)

// ContractRepository interface for database operations
type ContractRepository interface {
	GetByAddressAndNetwork(ctx context.Context, address string, network string) (*entity.Contract, error)
	Create(ctx context.Context, contract *entity.Contract) error
	CreateAsync(contract *entity.Contract) error
}

// EtherscanService interface for external API operations
type EtherscanService interface {
	GetContractSourceCode(ctx context.Context, address string, network string) (string, error)
	GetContractDetails(ctx context.Context, address string, network string) (*infrastructures.ContractInfo, error)
}

// ContractService handles contract-related business logic
type ContractService struct {
	contractRepo    ContractRepository
	etherscanClient EtherscanService
	logger          Logger
}

// Logger interface for logging operations
type Logger interface {
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
}

// DefaultLogger implements Logger interface
type DefaultLogger struct{}

func (d DefaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf("ERROR: "+format, args...)
}

func (d DefaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("INFO: "+format, args...)
}

// GormContractRepository implements ContractRepository
type GormContractRepository struct {
	db     *gorm.DB
	logger Logger
}

// NewGormContractRepository creates a new GORM-based contract repository
func NewGormContractRepository(db *gorm.DB, logger Logger) *GormContractRepository {
	return &GormContractRepository{
		db:     db,
		logger: logger,
	}
}

// GetByAddressAndNetwork retrieves a contract by its address and network
func (r *GormContractRepository) GetByAddressAndNetwork(ctx context.Context, address string, network string) (*entity.Contract, error) {
	var contract entity.Contract

	err := r.db.WithContext(ctx).Where("address = ? AND network = ?", address, network).First(&contract).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContractNotFound
		}
		return nil, fmt.Errorf("failed to get contract by address and network: %w", err)
	}

	return &contract, nil
}

// Create saves a new contract to the database
func (r *GormContractRepository) Create(ctx context.Context, contract *entity.Contract) error {
	if err := r.db.WithContext(ctx).Create(contract).Error; err != nil {
		return fmt.Errorf("failed to create contract: %w", err)
	}
	return nil
}

// CreateAsync saves a contract asynchronously
func (r *GormContractRepository) CreateAsync(contract *entity.Contract) error {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := r.Create(ctx, contract); err != nil {
			r.logger.Errorf("Failed to save contract asynchronously: %v", err)
		} else {
			r.logger.Infof("Contract saved successfully: %s", contract.Address)
		}
	}()

	return nil
}

// NewContractService creates a new contract service
func NewContractService(
	contractRepo ContractRepository,
	etherscanClient EtherscanService,
	logger Logger,
) *ContractService {
	if logger == nil {
		logger = DefaultLogger{}
	}

	return &ContractService{
		contractRepo:    contractRepo,
		etherscanClient: etherscanClient,
		logger:          logger,
	}
}

// NewContractServiceWithDefaults creates a service with default dependencies
func NewContractServiceWithDefaults(cfg *config.Config, db *gorm.DB) *ContractService {
	logger := DefaultLogger{}

	// Auto-migrate the contract table
	if err := db.AutoMigrate(&entity.Contract{}); err != nil {
		logger.Errorf("Failed to auto migrate contract table: %v", err)
	}

	contractRepo := NewGormContractRepository(db, logger)
	etherscanClient := infrastructures.NewEtherscanClient(&cfg.Etherscan)

	return NewContractService(contractRepo, etherscanClient, logger)
}

// GetAndSaveContractSourceCode retrieves contract source code and saves it if not exists
func (s *ContractService) GetAndSaveContractSourceCode(ctx context.Context, req *contracts.GetContractSourceCodeRequest) (*contracts.GetContractSourceCodeResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Check if contract already exists
	existingContract, err := s.contractRepo.GetByAddressAndNetwork(ctx, req.Address, req.Network)
	if err != nil && !errors.Is(err, ErrContractNotFound) {
		return nil, fmt.Errorf("failed to check existing contract: %w", err)
	}

	// If contract exists, return cached data
	if existingContract != nil {
		s.logger.Infof("Returning cached contract source code for address: %s on network: %s", req.Address, req.Network)
		return &contracts.GetContractSourceCodeResponse{
			Address:    existingContract.Address,
			SourceCode: existingContract.SourceCode,
			SourceHash: existingContract.SourceHash,
			Network:    existingContract.Network,
			CachedAt:   &existingContract.CreatedAt,
		}, nil
	}

	// Fetch contract details from Etherscan
	contractInfo, err := s.etherscanClient.GetContractDetails(ctx, req.Address, req.Network)
	if err != nil {
		var etherscanErr infrastructures.EtherscanError
		if errors.As(err, &etherscanErr) {
			// Etherscan API can return an error for various reasons (e.g. source not verified).
			// We'll treat these as a "not found" case.
			return nil, ErrContractNotFound
		}
		return nil, fmt.Errorf("failed to get contract details: %w", err)
	}

	if contractInfo.SourceCode == "" {
		return nil, ErrEmptySourceCode
	}

	// Generate hash for source code
	sourceHash := s.generateSourceHash(contractInfo.SourceCode)

	// Create contract entity
	contract := &entity.Contract{
		Address:         req.Address,
		Network:         req.Network,
		SourceCode:      contractInfo.SourceCode,
		ContractName:    contractInfo.ContractName,
		CompilerVersion: contractInfo.CompilerVersion,
		SourceHash:      sourceHash,
	}

	// Save contract asynchronously
	if err := s.contractRepo.CreateAsync(contract); err != nil {
		s.logger.Errorf("Failed to initiate async save for contract %s: %v", req.Address, err)
		// Don't return error here as the main operation (getting source code) succeeded
	}

	return &contracts.GetContractSourceCodeResponse{
		Address:    req.Address,
		SourceCode: contractInfo.SourceCode,
		SourceHash: sourceHash,
		Network:    req.Network,
		CachedAt:   nil, // Not cached yet
	}, nil
}

// GetContractByAddressAndNetwork retrieves a contract by its address and network
func (s *ContractService) GetContractByAddressAndNetwork(ctx context.Context, address string, network string) (*contracts.ContractResponse, error) {
	if address == "" {
		return nil, ErrInvalidAddress
	}

	contract, err := s.contractRepo.GetByAddressAndNetwork(ctx, address, network)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}

	return &contracts.ContractResponse{
		Address:         contract.Address,
		Network:         contract.Network,
		SourceCode:      contract.SourceCode,
		ContractName:    contract.ContractName,
		CompilerVersion: contract.CompilerVersion,
		SourceHash:      contract.SourceHash,
		CreatedAt:       contract.CreatedAt,
		UpdatedAt:       contract.UpdatedAt,
	}, nil
}

// generateSourceHash creates a SHA256 hash of the source code
func (s *ContractService) generateSourceHash(sourceCode string) string {
	hash := sha256.Sum256([]byte(sourceCode))
	return hex.EncodeToString(hash[:])
}

// RefreshContractSourceCode forces a refresh of contract source code from Etherscan
func (s *ContractService) RefreshContractSourceCode(ctx context.Context, req *contracts.RefreshContractSourceCodeRequest) (*contracts.GetContractSourceCodeResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Get fresh source code from Etherscan
	sourceCode, err := s.etherscanClient.GetContractSourceCode(ctx, req.Address, req.Network)
	if err != nil {
		var etherscanErr infrastructures.EtherscanError
		if errors.As(err, &etherscanErr) {
			return nil, ErrContractNotFound
		}
		return nil, fmt.Errorf("failed to refresh contract source code: %w", err)
	}

	if sourceCode == "" {
		return nil, ErrEmptySourceCode
	}

	sourceHash := s.generateSourceHash(sourceCode)

	// Check if contract exists and update or create
	existingContract, err := s.contractRepo.GetByAddressAndNetwork(ctx, req.Address, req.Network)
	if err != nil && !errors.Is(err, ErrContractNotFound) {
		return nil, fmt.Errorf("failed to check existing contract: %w", err)
	}

	contract := &entity.Contract{
		Address:    req.Address,
		Network:    req.Network,
		SourceCode: sourceCode,
		SourceHash: sourceHash,
	}

	if existingContract != nil {
		// Update existing contract
		contract.ID = existingContract.ID
		contract.CreatedAt = existingContract.CreatedAt
	}

	if err := s.contractRepo.CreateAsync(contract); err != nil {
		s.logger.Errorf("Failed to save refreshed contract %s: %v", req.Address, err)
	}

	return &contracts.GetContractSourceCodeResponse{
		Address:    req.Address,
		SourceCode: sourceCode,
		SourceHash: sourceHash,
		Network:    req.Network,
		CachedAt:   nil,
	}, nil
}
