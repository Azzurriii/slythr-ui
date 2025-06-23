package gorm

import (
	"context"
	"errors"

	"github.com/Azzurriii/slythr/internal/domain/entities"
	domainerrors "github.com/Azzurriii/slythr/internal/domain/errors"
	"github.com/Azzurriii/slythr/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type generatedTestCasesRepository struct {
	db *gorm.DB
}

// NewGeneratedTestCasesRepository creates a new generated test cases repository
func NewGeneratedTestCasesRepository(db *gorm.DB) repository.GeneratedTestCasesRepository {
	return &generatedTestCasesRepository{
		db: db,
	}
}

// SaveTestCases saves generated test cases to the database
func (r *generatedTestCasesRepository) SaveTestCases(ctx context.Context, testCases *entities.GeneratedTestCases) error {
	if testCases == nil {
		return errors.New("test cases cannot be nil")
	}

	if !testCases.IsValid() {
		return errors.New("invalid test cases data")
	}

	result := r.db.WithContext(ctx).
		Where("source_hash = ?", testCases.SourceHash).
		Assign("test_code = ?", testCases.TestCode).
		Assign("test_framework = ?", testCases.TestFramework).
		Assign("test_language = ?", testCases.TestLanguage).
		Assign("file_name = ?", testCases.FileName).
		Assign("warnings_and_recommendations = ?", testCases.WarningsAndRecommendations).
		Assign("updated_at = ?", testCases.UpdatedAt).
		FirstOrCreate(testCases)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// FindBySourceHash finds generated test cases by source hash
func (r *generatedTestCasesRepository) FindBySourceHash(ctx context.Context, sourceHash string) (*entities.GeneratedTestCases, error) {
	if sourceHash == "" {
		return nil, errors.New("source hash cannot be empty")
	}

	var testCases entities.GeneratedTestCases
	result := r.db.WithContext(ctx).
		Where("source_hash = ?", sourceHash).
		First(&testCases)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrStaticAnalysisNotFound // Reuse existing error for now
		}
		return nil, result.Error
	}

	return &testCases, nil
}

// ExistsBySourceHash checks if generated test cases exist by source hash
func (r *generatedTestCasesRepository) ExistsBySourceHash(ctx context.Context, sourceHash string) (bool, error) {
	if sourceHash == "" {
		return false, errors.New("source hash cannot be empty")
	}

	var count int64
	result := r.db.WithContext(ctx).
		Model(&entities.GeneratedTestCases{}).
		Where("source_hash = ?", sourceHash).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

// FindByID finds generated test cases by ID
func (r *generatedTestCasesRepository) FindByID(ctx context.Context, id entities.GeneratedTestCasesID) (*entities.GeneratedTestCases, error) {
	if id == entities.GeneratedTestCasesID(uuid.Nil) {
		return nil, errors.New("invalid test cases ID")
	}

	var testCases entities.GeneratedTestCases
	result := r.db.WithContext(ctx).First(&testCases, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrStaticAnalysisNotFound
		}
		return nil, result.Error
	}

	return &testCases, nil
}

// UpdateTestCases updates existing generated test cases
func (r *generatedTestCasesRepository) UpdateTestCases(ctx context.Context, testCases *entities.GeneratedTestCases) error {
	if testCases == nil {
		return errors.New("test cases cannot be nil")
	}

	if !testCases.IsValid() {
		return errors.New("invalid test cases data")
	}

	result := r.db.WithContext(ctx).Save(testCases)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domainerrors.ErrStaticAnalysisNotFound
	}

	return nil
}

// RemoveTestCases deletes generated test cases
func (r *generatedTestCasesRepository) RemoveTestCases(ctx context.Context, id entities.GeneratedTestCasesID) error {
	if id == entities.GeneratedTestCasesID(uuid.Nil) {
		return errors.New("invalid test cases ID")
	}

	result := r.db.WithContext(ctx).Delete(&entities.GeneratedTestCases{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domainerrors.ErrStaticAnalysisNotFound
	}

	return nil
}

// ListTestCases retrieves generated test cases with pagination
func (r *generatedTestCasesRepository) ListTestCases(ctx context.Context, offset, limit int) ([]*entities.GeneratedTestCases, error) {
	var testCases []*entities.GeneratedTestCases
	result := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&testCases)

	if result.Error != nil {
		return nil, result.Error
	}

	return testCases, nil
}

// CountTestCases returns the total number of generated test cases
func (r *generatedTestCasesRepository) CountTestCases(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&entities.GeneratedTestCases{}).
		Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}
