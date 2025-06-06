package entity

import (
	"gorm.io/gorm"
)

type Contract struct {
	gorm.Model
	Address         string `gorm:"uniqueIndex:idx_address_network;not null"`
	Network         string `gorm:"uniqueIndex:idx_address_network;not null"`
	SourceCode      string `gorm:"type:text;not null"`
	ContractName    string `gorm:"not null"`
	CompilerVersion string `gorm:"not null"`
	SourceHash      string `gorm:"not null"`
}

// Business methods
func (c *Contract) IsValid() error {
	// Validation logic
	return nil
}
