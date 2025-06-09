package repository

import (
	"context"

	"github.com/Azzurriii/slythr-go-backend/internal/domain/entities"
)

// StaticAnalysisRepository defines the contract for static analysis data access
type StaticAnalysisRepository interface {
	// SaveAnalysis creates a new static analysis
	SaveAnalysis(ctx context.Context, analysis *entities.StaticAnalysis) error

	// FindByID retrieves a static analysis by its ID
	FindByID(ctx context.Context, id entities.StaticAnalysisID) (*entities.StaticAnalysis, error)

	// FindByContractID retrieves static analyses by contract ID
	FindByContractID(ctx context.Context, contractID entities.ContractID) ([]*entities.StaticAnalysis, error)

	// FindBySourceHash retrieves static analysis by source hash
	FindBySourceHash(ctx context.Context, sourceHash string) (*entities.StaticAnalysis, error)

	// FindLatestByContractID retrieves the latest static analysis for a contract
	FindLatestByContractID(ctx context.Context, contractID entities.ContractID) (*entities.StaticAnalysis, error)

	// UpdateAnalysis updates an existing static analysis
	UpdateAnalysis(ctx context.Context, analysis *entities.StaticAnalysis) error

	// RemoveAnalysis deletes a static analysis
	RemoveAnalysis(ctx context.Context, id entities.StaticAnalysisID) error

	// ListAnalyses retrieves static analyses with pagination
	ListAnalyses(ctx context.Context, offset, limit int) ([]*entities.StaticAnalysis, error)

	// ExistsBySourceHash checks if a static analysis exists for the given source hash
	ExistsBySourceHash(ctx context.Context, sourceHash string) (bool, error)

	// CountAnalyses returns the total number of analyses
	CountAnalyses(ctx context.Context) (int64, error)
}
