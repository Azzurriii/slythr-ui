package repository

import (
	"context"

	"github.com/Azzurriii/slythr/internal/domain/entities"
)

// DynamicAnalysisRepository defines the contract for dynamic analysis data access
type DynamicAnalysisRepository interface {
	// SaveAnalysis creates a new dynamic analysis
	SaveAnalysis(ctx context.Context, analysis *entities.DynamicAnalysis) error

	// FindByID retrieves a dynamic analysis by its ID
	FindByID(ctx context.Context, id entities.DynamicAnalysisID) (*entities.DynamicAnalysis, error)

	// FindByContractID retrieves dynamic analyses by contract ID
	FindByContractID(ctx context.Context, contractID entities.ContractID) ([]*entities.DynamicAnalysis, error)

	// FindBySourceHash retrieves dynamic analysis by source hash
	FindBySourceHash(ctx context.Context, sourceHash string) (*entities.DynamicAnalysis, error)

	// FindLatestByContractID retrieves the latest dynamic analysis for a contract
	FindLatestByContractID(ctx context.Context, contractID entities.ContractID) (*entities.DynamicAnalysis, error)

	// UpdateAnalysis updates an existing dynamic analysis
	UpdateAnalysis(ctx context.Context, analysis *entities.DynamicAnalysis) error

	// RemoveAnalysis deletes a dynamic analysis
	RemoveAnalysis(ctx context.Context, id entities.DynamicAnalysisID) error

	// ListAnalyses retrieves dynamic analyses with pagination
	ListAnalyses(ctx context.Context, offset, limit int) ([]*entities.DynamicAnalysis, error)

	// ExistsBySourceHash checks if a dynamic analysis exists for the given source hash
	ExistsBySourceHash(ctx context.Context, sourceHash string) (bool, error)

	// CountAnalyses returns the total number of analyses
	CountAnalyses(ctx context.Context) (int64, error)
}
