package repository

import (
	"context"

	"github.com/Azzurriii/slythr-go-backend/internal/domain/entities"
)

// ContractRepository defines the contract for contract data access
type ContractRepository interface {
	// SaveContract creates a new contract
	SaveContract(ctx context.Context, contract *entities.Contract) error

	// FindByID retrieves a contract by its ID
	FindByID(ctx context.Context, id entities.ContractID) (*entities.Contract, error)

	// FindByAddressAndNetwork retrieves a contract by address and network
	FindByAddressAndNetwork(ctx context.Context, address, network string) (*entities.Contract, error)

	// FindBySourceHash retrieves contracts by source hash
	FindBySourceHash(ctx context.Context, sourceHash string) ([]*entities.Contract, error)

	// UpdateContract updates an existing contract
	UpdateContract(ctx context.Context, contract *entities.Contract) error

	// RemoveContract soft deletes a contract
	RemoveContract(ctx context.Context, id entities.ContractID) error

	// ExistsByAddress checks if a contract exists by address and network
	ExistsByAddressAndNetwork(ctx context.Context, address, network string) (bool, error)

	// ListContracts retrieves contracts with pagination
	ListContracts(ctx context.Context, offset, limit int) ([]*entities.Contract, error)

	// CountContracts returns the total number of contracts
	CountContracts(ctx context.Context) (int64, error)
}
