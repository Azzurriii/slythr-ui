package services

import (
	"context"

	"github.com/Azzurriii/slythr-go-backend/internal/application/dto/contracts"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/entities"
)

// ContractServiceInterface defines the contract service interface
type ContractServiceInterface interface {
	// FetchContractSourceCode retrieves contract source code and saves it if not exists
	FetchContractSourceCode(ctx context.Context, req *contracts.GetContractSourceCodeRequest) (*contracts.GetContractSourceCodeResponse, error)

	// GetContractByAddressAndNetwork retrieves a contract by its address and network
	GetContractByAddressAndNetwork(ctx context.Context, address, network string) (*contracts.ContractResponse, error)
}

// StaticAnalysisServiceInterface defines the static analysis service interface
type StaticAnalysisServiceInterface interface {
	// CreateAnalysis creates a new static analysis
	CreateAnalysis(ctx context.Context, contractID entities.ContractID, results string) (*entities.StaticAnalysis, error)

	// GetAnalysisByContractID retrieves static analysis by contract ID
	GetAnalysisByContractID(ctx context.Context, contractID entities.ContractID) (*entities.StaticAnalysis, error)

	// GetAnalysisBySourceHash retrieves static analysis by source hash
	GetAnalysisBySourceHash(ctx context.Context, sourceHash string) (*entities.StaticAnalysis, error)
}

// DynamicAnalysisServiceInterface defines the dynamic analysis service interface
type DynamicAnalysisServiceInterface interface {
	// CreateAnalysis creates a new dynamic analysis
	CreateAnalysis(ctx context.Context, contractID entities.ContractID, llmResponse string) (*entities.DynamicAnalysis, error)

	// GetAnalysisByContractID retrieves dynamic analysis by contract ID
	GetAnalysisByContractID(ctx context.Context, contractID entities.ContractID) (*entities.DynamicAnalysis, error)

	// GetAnalysisBySourceHash retrieves dynamic analysis by source hash
	GetAnalysisBySourceHash(ctx context.Context, sourceHash string) (*entities.DynamicAnalysis, error)
}
