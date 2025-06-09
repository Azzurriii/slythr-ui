package entities

import (
	"time"

	"github.com/Azzurriii/slythr-go-backend/internal/domain/valueobjects"
	"gorm.io/gorm"
)

// Contract represents a smart contract entity in the domain
type Contract struct {
	ID              ContractID     `gorm:"primaryKey"`
	Address         string         `gorm:"uniqueIndex:idx_address_network;not null;size:42"`
	Network         string         `gorm:"uniqueIndex:idx_address_network;not null;size:20"`
	SourceCode      string         `gorm:"type:text;not null"`
	ContractName    string         `gorm:"not null;size:255"`
	CompilerVersion string         `gorm:"not null;size:50"`
	SourceHash      string         `gorm:"not null;size:64;index"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

// ContractID represents the unique identifier for a contract
type ContractID uint

// NewContract creates a new contract with validation
func NewContract(
	address string,
	network string,
	sourceCode string,
	contractName string,
	compilerVersion string,
	sourceHash string,
) (*Contract, error) {
	// Validate using value objects
	addressVO, err := valueobjects.NewContractAddress(address)
	if err != nil {
		return nil, err
	}

	networkVO, err := valueobjects.NewNetwork(network)
	if err != nil {
		return nil, err
	}

	sourceCodeVO, err := valueobjects.NewSourceCode(sourceCode)
	if err != nil {
		return nil, err
	}

	contractNameVO, err := valueobjects.NewContractName(contractName)
	if err != nil {
		return nil, err
	}

	compilerVersionVO, err := valueobjects.NewCompilerVersion(compilerVersion)
	if err != nil {
		return nil, err
	}

	sourceHashVO, err := valueobjects.NewSourceHash(sourceHash)
	if err != nil {
		return nil, err
	}

	return &Contract{
		Address:         addressVO.Value(),
		Network:         networkVO.Value(),
		SourceCode:      sourceCodeVO.Value(),
		ContractName:    contractNameVO.Value(),
		CompilerVersion: compilerVersionVO.Value(),
		SourceHash:      sourceHashVO.Value(),
	}, nil
}

// GetID returns the contract ID
func (c *Contract) GetID() ContractID {
	return c.ID
}

// GetAddress returns the contract address as value object
func (c *Contract) GetAddress() valueobjects.ContractAddress {
	address, _ := valueobjects.NewContractAddress(c.Address)
	return address
}

// GetNetwork returns the network as value object
func (c *Contract) GetNetwork() valueobjects.Network {
	network, _ := valueobjects.NewNetwork(c.Network)
	return network
}

// GetSourceCode returns the source code as value object
func (c *Contract) GetSourceCode() valueobjects.SourceCode {
	sourceCode, _ := valueobjects.NewSourceCode(c.SourceCode)
	return sourceCode
}

// GetSourceHash returns the source hash as value object
func (c *Contract) GetSourceHash() valueobjects.SourceHash {
	sourceHash, _ := valueobjects.NewSourceHash(c.SourceHash)
	return sourceHash
}

// IsValid checks if the contract is valid
func (c *Contract) IsValid() bool {
	address, err := valueobjects.NewContractAddress(c.Address)
	if err != nil {
		return false
	}

	network, err := valueobjects.NewNetwork(c.Network)
	if err != nil {
		return false
	}

	sourceCode, err := valueobjects.NewSourceCode(c.SourceCode)
	if err != nil {
		return false
	}

	contractName, err := valueobjects.NewContractName(c.ContractName)
	if err != nil {
		return false
	}

	compilerVersion, err := valueobjects.NewCompilerVersion(c.CompilerVersion)
	if err != nil {
		return false
	}

	sourceHash, err := valueobjects.NewSourceHash(c.SourceHash)
	if err != nil {
		return false
	}

	return address.IsValid() &&
		network.IsValid() &&
		sourceCode.IsValid() &&
		contractName.IsValid() &&
		compilerVersion.IsValid() &&
		sourceHash.IsValid()
}

// HasSourceCode checks if the contract has valid source code
func (c *Contract) HasSourceCode() bool {
	sourceCode := c.GetSourceCode()
	return sourceCode.HasContent()
}

// TableName returns the table name for GORM
func (Contract) TableName() string {
	return "contracts"
}
