package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/Azzurriii/slythr-go-backend/internal/application/dto/contracts"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/entities"
	domainerrors "github.com/Azzurriii/slythr-go-backend/internal/domain/errors"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/repository"
	"github.com/Azzurriii/slythr-go-backend/internal/infrastructure/external"
)

type ContractService struct {
	contractRepo    repository.ContractRepository
	etherscanClient external.EtherscanService
	hashGenerator   HashGenerator
	logger          Logger
}

type Logger interface {
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

type defaultHashGenerator struct{}

func (h *defaultHashGenerator) GenerateSourceHash(sourceCode string) string {
	hash := sha256.Sum256([]byte(sourceCode))
	return hex.EncodeToString(hash[:])
}

func NewContractService(
	contractRepo repository.ContractRepository,
	etherscanClient external.EtherscanService,
	logger Logger,
) ContractServiceInterface {
	return &ContractService{
		contractRepo:    contractRepo,
		etherscanClient: etherscanClient,
		hashGenerator:   &defaultHashGenerator{},
		logger:          logger,
	}
}

func (s *ContractService) FetchContractSourceCode(ctx context.Context, req *contracts.GetContractSourceCodeRequest) (*contracts.GetContractSourceCodeResponse, error) {
	if err := req.Validate(); err != nil {
		s.logger.Warnf("Invalid request for contract source code: %v", err)
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	existingContract, err := s.contractRepo.FindByAddressAndNetwork(ctx, req.Address, req.Network)
	if err != nil && !errors.Is(err, domainerrors.ErrContractNotFound) {
		s.logger.Errorf("Failed to check existing contract: %v", err)
		return nil, fmt.Errorf("failed to check existing contract: %w", err)
	}

	if existingContract != nil {
		s.logger.Infof("Returning cached contract source code for address: %s on network: %s", req.Address, req.Network)
		return s.buildSourceCodeResponse(existingContract, true), nil
	}

	contractInfo, err := s.fetchContractFromExternal(ctx, req.Address, req.Network)
	if err != nil {
		return nil, err
	}

	contract, err := s.createContractEntity(req.Address, req.Network, contractInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract entity: %w", err)
	}

	go s.saveContractAsync(contract)

	return s.buildSourceCodeResponse(contract, false), nil
}

func (s *ContractService) GetContractByAddressAndNetwork(ctx context.Context, address, network string) (*contracts.ContractResponse, error) {
	contract, err := s.contractRepo.FindByAddressAndNetwork(ctx, address, network)
	if err != nil {
		s.logger.Errorf("Failed to get contract by address and network: %v", err)
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}

	return s.buildContractResponse(contract), nil
}

func (s *ContractService) fetchContractFromExternal(ctx context.Context, address, network string) (*external.ContractInfo, error) {
	processedSourceCode, err := s.etherscanClient.GetContractSourceCode(ctx, address, network)
	if err != nil {
		var etherscanErr external.EtherscanError
		if errors.As(err, &etherscanErr) {
			s.logger.Warnf("Contract not found on external service: %s", address)
			return nil, domainerrors.ErrContractNotFound
		}
		s.logger.Errorf("Failed to get contract source code from external service: %v", err)
		return nil, fmt.Errorf("failed to get contract source code: %w", err)
	}

	if processedSourceCode == "" {
		s.logger.Warnf("Empty source code returned for contract: %s", address)
		return nil, domainerrors.ErrEmptySourceCode
	}

	contractInfo, err := s.etherscanClient.GetContractDetails(ctx, address, network)
	if err != nil {
		var etherscanErr external.EtherscanError
		if errors.As(err, &etherscanErr) {
			s.logger.Warnf("Contract details not found on external service: %s", address)
			return nil, domainerrors.ErrContractNotFound
		}
		s.logger.Errorf("Failed to get contract details from external service: %v", err)
		return nil, fmt.Errorf("failed to get contract details: %w", err)
	}

	return contractInfo, nil
}

func (s *ContractService) createContractEntity(address, network string, contractInfo *external.ContractInfo) (*entities.Contract, error) {
	sourceHash := s.hashGenerator.GenerateSourceHash(contractInfo.SourceCode)

	return entities.NewContract(
		address,
		network,
		contractInfo.SourceCode,
		contractInfo.ContractName,
		contractInfo.CompilerVersion,
		sourceHash,
	)
}

func (s *ContractService) buildSourceCodeResponse(contract *entities.Contract, cached bool) *contracts.GetContractSourceCodeResponse {
	response := &contracts.GetContractSourceCodeResponse{
		Address:    contract.Address,
		SourceCode: contract.SourceCode,
		SourceHash: contract.SourceHash,
		Network:    contract.Network,
	}

	if cached {
		response.CachedAt = &contract.CreatedAt
	}

	return response
}

func (s *ContractService) buildContractResponse(contract *entities.Contract) *contracts.ContractResponse {
	return &contracts.ContractResponse{
		Address:         contract.Address,
		Network:         contract.Network,
		SourceCode:      contract.SourceCode,
		ContractName:    contract.ContractName,
		CompilerVersion: contract.CompilerVersion,
		SourceHash:      contract.SourceHash,
		CreatedAt:       contract.CreatedAt,
		UpdatedAt:       contract.UpdatedAt,
	}
}

func (s *ContractService) saveContractAsync(contract *entities.Contract) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.contractRepo.SaveContract(ctx, contract); err != nil {
		s.logger.Errorf("Failed to save contract asynchronously: %v", err)
	} else {
		s.logger.Infof("Contract saved successfully: %s", contract.Address)
	}
}
