package entities

import (
	"time"

	"github.com/Azzurriii/slythr/internal/domain/valueobjects"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Contract struct {
	ID              ContractID     `gorm:"type:uuid;primaryKey"`
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

type ContractID uuid.UUID

func NewContract(
	address string,
	network string,
	sourceCode string,
	contractName string,
	compilerVersion string,
	sourceHash string,
) (*Contract, error) {
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
		ID:              ContractID(uuid.New()),
		Address:         addressVO.Value(),
		Network:         networkVO.Value(),
		SourceCode:      sourceCodeVO.Value(),
		ContractName:    contractNameVO.Value(),
		CompilerVersion: compilerVersionVO.Value(),
		SourceHash:      sourceHashVO.Value(),
	}, nil
}

func (c *Contract) GetID() ContractID {
	return c.ID
}

func (c *Contract) GetAddress() valueobjects.ContractAddress {
	address, _ := valueobjects.NewContractAddress(c.Address)
	return address
}

func (c *Contract) GetNetwork() valueobjects.Network {
	network, _ := valueobjects.NewNetwork(c.Network)
	return network
}

func (c *Contract) GetSourceCode() valueobjects.SourceCode {
	sourceCode, _ := valueobjects.NewSourceCode(c.SourceCode)
	return sourceCode
}

func (c *Contract) GetSourceHash() valueobjects.SourceHash {
	sourceHash, _ := valueobjects.NewSourceHash(c.SourceHash)
	return sourceHash
}

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

func (c *Contract) HasSourceCode() bool {
	sourceCode := c.GetSourceCode()
	return sourceCode.HasContent()
}

func (Contract) TableName() string {
	return "contracts"
}
