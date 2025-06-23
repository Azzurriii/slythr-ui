package repository

import (
	"context"

	"github.com/Azzurriii/slythr/internal/domain/entities"
)

// GeneratedTestCasesRepository defines the contract for generated test cases data access
type GeneratedTestCasesRepository interface {
	// SaveTestCases creates new generated test cases
	SaveTestCases(ctx context.Context, testCases *entities.GeneratedTestCases) error

	// FindByID retrieves generated test cases by its ID
	FindByID(ctx context.Context, id entities.GeneratedTestCasesID) (*entities.GeneratedTestCases, error)

	// FindBySourceHash retrieves generated test cases by source hash
	FindBySourceHash(ctx context.Context, sourceHash string) (*entities.GeneratedTestCases, error)

	// UpdateTestCases updates existing generated test cases
	UpdateTestCases(ctx context.Context, testCases *entities.GeneratedTestCases) error

	// RemoveTestCases deletes generated test cases
	RemoveTestCases(ctx context.Context, id entities.GeneratedTestCasesID) error

	// ListTestCases retrieves generated test cases with pagination
	ListTestCases(ctx context.Context, offset, limit int) ([]*entities.GeneratedTestCases, error)

	// ExistsBySourceHash checks if generated test cases exist for the given source hash
	ExistsBySourceHash(ctx context.Context, sourceHash string) (bool, error)

	// CountTestCases returns the total number of generated test cases
	CountTestCases(ctx context.Context) (int64, error)
}
