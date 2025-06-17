package gorm

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azzurriii/slythr-go-backend/internal/domain/entities"
	domainerrors "github.com/Azzurriii/slythr-go-backend/internal/domain/errors"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type dynamicAnalysisRepository struct {
	db *gorm.DB
}

func NewDynamicAnalysisRepository(db *gorm.DB) repository.DynamicAnalysisRepository {
	return &dynamicAnalysisRepository{
		db: db,
	}
}

func (r *dynamicAnalysisRepository) SaveAnalysis(ctx context.Context, analysis *entities.DynamicAnalysis) error {
	if err := r.db.WithContext(ctx).Create(analysis).Error; err != nil {
		if r.isDuplicateKeyError(err) {
			return domainerrors.ErrAnalysisAlreadyExists
		}
		return fmt.Errorf("failed to save dynamic analysis: %w", err)
	}
	return nil
}

func (r *dynamicAnalysisRepository) FindByID(ctx context.Context, id entities.DynamicAnalysisID) (*entities.DynamicAnalysis, error) {
	var analysis entities.DynamicAnalysis
	err := r.db.WithContext(ctx).First(&analysis, uuid.UUID(id)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrAnalysisNotFound
		}
		return nil, fmt.Errorf("failed to find dynamic analysis by ID: %w", err)
	}
	return &analysis, nil
}

func (r *dynamicAnalysisRepository) FindByContractID(ctx context.Context, contractID entities.ContractID) ([]*entities.DynamicAnalysis, error) {
	var analyses []*entities.DynamicAnalysis
	err := r.db.WithContext(ctx).
		Where("contract_id = ?", uuid.UUID(contractID)).
		Order("created_at DESC").
		Find(&analyses).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find dynamic analyses by contract ID: %w", err)
	}
	return analyses, nil
}

func (r *dynamicAnalysisRepository) FindBySourceHash(ctx context.Context, sourceHash string) (*entities.DynamicAnalysis, error) {
	var analysis entities.DynamicAnalysis
	err := r.db.WithContext(ctx).
		Where("source_hash = ?", sourceHash).
		Order("created_at DESC").
		First(&analysis).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrAnalysisNotFound
		}
		return nil, fmt.Errorf("failed to find dynamic analysis by source hash: %w", err)
	}
	return &analysis, nil
}

func (r *dynamicAnalysisRepository) FindLatestByContractID(ctx context.Context, contractID entities.ContractID) (*entities.DynamicAnalysis, error) {
	var analysis entities.DynamicAnalysis
	err := r.db.WithContext(ctx).
		Where("contract_id = ?", uuid.UUID(contractID)).
		Order("created_at DESC").
		First(&analysis).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrAnalysisNotFound
		}
		return nil, fmt.Errorf("failed to find latest dynamic analysis: %w", err)
	}
	return &analysis, nil
}

func (r *dynamicAnalysisRepository) UpdateAnalysis(ctx context.Context, analysis *entities.DynamicAnalysis) error {
	result := r.db.WithContext(ctx).Save(analysis)
	if result.Error != nil {
		return fmt.Errorf("failed to update dynamic analysis: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domainerrors.ErrAnalysisNotFound
	}

	return nil
}

func (r *dynamicAnalysisRepository) RemoveAnalysis(ctx context.Context, id entities.DynamicAnalysisID) error {
	result := r.db.WithContext(ctx).Delete(&entities.DynamicAnalysis{}, uuid.UUID(id))
	if result.Error != nil {
		return fmt.Errorf("failed to remove dynamic analysis: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domainerrors.ErrAnalysisNotFound
	}

	return nil
}

func (r *dynamicAnalysisRepository) ListAnalyses(ctx context.Context, offset, limit int) ([]*entities.DynamicAnalysis, error) {
	var analyses []*entities.DynamicAnalysis
	err := r.db.WithContext(ctx).
		Preload("Contract").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&analyses).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list dynamic analyses: %w", err)
	}
	return analyses, nil
}

func (r *dynamicAnalysisRepository) ExistsBySourceHash(ctx context.Context, sourceHash string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.DynamicAnalysis{}).
		Where("source_hash = ?", sourceHash).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check dynamic analysis existence by source hash: %w", err)
	}

	return count > 0, nil
}

func (r *dynamicAnalysisRepository) CountAnalyses(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.DynamicAnalysis{}).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count dynamic analyses: %w", err)
	}

	return count, nil
}

func (r *dynamicAnalysisRepository) isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	return r.containsErrorCode(err.Error(), "23505") ||
		r.containsErrorCode(err.Error(), "duplicate key")
}

func (r *dynamicAnalysisRepository) containsErrorCode(errMsg, code string) bool {
	if errMsg == "" || code == "" {
		return false
	}

	return len(errMsg) >= len(code) && r.indexOf(errMsg, code) >= 0
}

func (r *dynamicAnalysisRepository) indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
