package gorm

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azzurriii/slythr/internal/domain/entities"
	domainerrors "github.com/Azzurriii/slythr/internal/domain/errors"
	"github.com/Azzurriii/slythr/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type contractRepository struct {
	db *gorm.DB
}

func NewContractRepository(db *gorm.DB) repository.ContractRepository {
	return &contractRepository{
		db: db,
	}
}

func (r *contractRepository) SaveContract(ctx context.Context, contract *entities.Contract) error {
	if err := r.db.WithContext(ctx).Create(contract).Error; err != nil {
		if r.isDuplicateKeyError(err) {
			return domainerrors.ErrContractAlreadyExists
		}
		return fmt.Errorf("failed to save contract: %w", err)
	}
	return nil
}

func (r *contractRepository) FindByID(ctx context.Context, id entities.ContractID) (*entities.Contract, error) {
	var contract entities.Contract
	err := r.db.WithContext(ctx).First(&contract, uuid.UUID(id)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrContractNotFound
		}
		return nil, fmt.Errorf("failed to find contract by ID: %w", err)
	}
	return &contract, nil
}

func (r *contractRepository) FindByAddressAndNetwork(ctx context.Context, address, network string) (*entities.Contract, error) {
	var contract entities.Contract
	err := r.db.WithContext(ctx).
		Where("address = ? AND network = ?", address, network).
		First(&contract).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrContractNotFound
		}
		return nil, fmt.Errorf("failed to find contract by address and network: %w", err)
	}
	return &contract, nil
}

func (r *contractRepository) FindBySourceHash(ctx context.Context, sourceHash string) ([]*entities.Contract, error) {
	var contracts []*entities.Contract
	err := r.db.WithContext(ctx).
		Where("source_hash = ?", sourceHash).
		Find(&contracts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find contracts by source hash: %w", err)
	}
	return contracts, nil
}

func (r *contractRepository) FindFirstBySourceHash(ctx context.Context, sourceHash string) (*entities.Contract, error) {
	var contract entities.Contract
	err := r.db.WithContext(ctx).
		Where("source_hash = ?", sourceHash).
		First(&contract).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrContractNotFound
		}
		return nil, fmt.Errorf("failed to find contract by source hash: %w", err)
	}
	return &contract, nil
}

func (r *contractRepository) ExistsBySourceHash(ctx context.Context, sourceHash string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Contract{}).
		Where("source_hash = ?", sourceHash).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check contract existence by source hash: %w", err)
	}

	return count > 0, nil
}

func (r *contractRepository) UpdateContract(ctx context.Context, contract *entities.Contract) error {
	result := r.db.WithContext(ctx).Save(contract)
	if result.Error != nil {
		return fmt.Errorf("failed to update contract: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domainerrors.ErrContractNotFound
	}

	return nil
}

func (r *contractRepository) RemoveContract(ctx context.Context, id entities.ContractID) error {
	result := r.db.WithContext(ctx).Delete(&entities.Contract{}, uuid.UUID(id))
	if result.Error != nil {
		return fmt.Errorf("failed to remove contract: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domainerrors.ErrContractNotFound
	}

	return nil
}

func (r *contractRepository) ExistsByAddressAndNetwork(ctx context.Context, address, network string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Contract{}).
		Where("address = ? AND network = ?", address, network).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check contract existence: %w", err)
	}

	return count > 0, nil
}

func (r *contractRepository) ListContracts(ctx context.Context, offset, limit int) ([]*entities.Contract, error) {
	var contracts []*entities.Contract
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&contracts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list contracts: %w", err)
	}
	return contracts, nil
}

func (r *contractRepository) CountContracts(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Contract{}).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count contracts: %w", err)
	}

	return count, nil
}

func (r *contractRepository) isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	return r.containsErrorCode(err.Error(), "23505") ||
		r.containsErrorCode(err.Error(), "duplicate key")
}

func (r *contractRepository) containsErrorCode(errMsg, code string) bool {
	if errMsg == "" || code == "" {
		return false
	}

	return len(errMsg) >= len(code) && r.indexOf(errMsg, code) >= 0
}

func (r *contractRepository) indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
