package gorm

import (
	"context"
	"errors"

	"github.com/Azzurriii/slythr-go-backend/internal/domain/entities"
	domainerrors "github.com/Azzurriii/slythr-go-backend/internal/domain/errors"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type staticAnalysisRepository struct {
	db *gorm.DB
}

// NewStaticAnalysisRepository creates a new static analysis repository
func NewStaticAnalysisRepository(db *gorm.DB) repository.StaticAnalysisRepository {
	return &staticAnalysisRepository{
		db: db,
	}
}

// SaveAnalysis saves a static analysis to the database
func (r *staticAnalysisRepository) SaveAnalysis(ctx context.Context, analysis *entities.StaticAnalysis) error {
	if analysis == nil {
		return errors.New("analysis cannot be nil")
	}

	if !analysis.IsValid() {
		return errors.New("invalid analysis data")
	}

	// Use upsert to handle duplicates
	result := r.db.WithContext(ctx).
		Where("source_hash = ?", analysis.SourceHash).
		Assign("slither_output = ?", analysis.SlitherOutput).
		FirstOrCreate(analysis)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// FindBySourceHash finds a static analysis by source hash
func (r *staticAnalysisRepository) FindBySourceHash(ctx context.Context, sourceHash string) (*entities.StaticAnalysis, error) {
	if sourceHash == "" {
		return nil, errors.New("source hash cannot be empty")
	}

	var analysis entities.StaticAnalysis
	result := r.db.WithContext(ctx).
		Where("source_hash = ?", sourceHash).
		First(&analysis)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrStaticAnalysisNotFound
		}
		return nil, result.Error
	}

	return &analysis, nil
}

// ExistsBySourceHash checks if a static analysis exists by source hash
func (r *staticAnalysisRepository) ExistsBySourceHash(ctx context.Context, sourceHash string) (bool, error) {
	if sourceHash == "" {
		return false, errors.New("source hash cannot be empty")
	}

	var count int64
	result := r.db.WithContext(ctx).
		Model(&entities.StaticAnalysis{}).
		Where("source_hash = ?", sourceHash).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

// FindByID finds a static analysis by ID
func (r *staticAnalysisRepository) FindByID(ctx context.Context, id entities.StaticAnalysisID) (*entities.StaticAnalysis, error) {
	if id == entities.StaticAnalysisID(uuid.Nil) {
		return nil, errors.New("invalid analysis ID")
	}

	var analysis entities.StaticAnalysis
	result := r.db.WithContext(ctx).First(&analysis, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrStaticAnalysisNotFound
		}
		return nil, result.Error
	}

	return &analysis, nil
}

// FindByContractID finds static analyses by contract ID
func (r *staticAnalysisRepository) FindByContractID(ctx context.Context, contractID entities.ContractID) ([]*entities.StaticAnalysis, error) {
	// Since we don't have contract_id in static_analysis table, return empty slice
	return []*entities.StaticAnalysis{}, nil
}

// FindLatestByContractID finds the latest static analysis for a contract
func (r *staticAnalysisRepository) FindLatestByContractID(ctx context.Context, contractID entities.ContractID) (*entities.StaticAnalysis, error) {
	// Since we don't have contract_id in static_analysis table, return not found
	return nil, domainerrors.ErrStaticAnalysisNotFound
}

// UpdateAnalysis updates an existing static analysis
func (r *staticAnalysisRepository) UpdateAnalysis(ctx context.Context, analysis *entities.StaticAnalysis) error {
	if analysis == nil {
		return errors.New("analysis cannot be nil")
	}

	if !analysis.IsValid() {
		return errors.New("invalid analysis data")
	}

	result := r.db.WithContext(ctx).Save(analysis)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domainerrors.ErrStaticAnalysisNotFound
	}

	return nil
}

// RemoveAnalysis deletes a static analysis by ID
func (r *staticAnalysisRepository) RemoveAnalysis(ctx context.Context, id entities.StaticAnalysisID) error {
	if id == entities.StaticAnalysisID(uuid.Nil) {
		return errors.New("invalid analysis ID")
	}

	result := r.db.WithContext(ctx).Delete(&entities.StaticAnalysis{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domainerrors.ErrStaticAnalysisNotFound
	}

	return nil
}

// ListAnalyses retrieves static analyses with pagination
func (r *staticAnalysisRepository) ListAnalyses(ctx context.Context, offset, limit int) ([]*entities.StaticAnalysis, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	var analyses []*entities.StaticAnalysis
	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&analyses)

	if result.Error != nil {
		return nil, result.Error
	}

	return analyses, nil
}

// CountAnalyses returns the total number of static analyses
func (r *staticAnalysisRepository) CountAnalyses(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&entities.StaticAnalysis{}).
		Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}
