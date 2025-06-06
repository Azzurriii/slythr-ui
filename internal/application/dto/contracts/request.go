package contracts

import (
	"strings"

	errors "github.com/Azzurriii/slythr-go-backend/internal/domain/errors"
)

var (
	ErrInvalidAddress = errors.ErrInvalidAddress
)

// GetContractSourceCodeRequest represents the request for getting contract source code
type GetContractSourceCodeRequest struct {
	Address string `json:"address" binding:"required"`
	Network string `json:"network" binding:"required"`
}

// Validate validates the request
func (r *GetContractSourceCodeRequest) Validate() error {
	if r.Address == "" {
		return errors.ErrInvalidAddress
	}

	if !strings.HasPrefix(r.Address, "0x") {
		return errors.ErrInvalidAddress
	}

	if len(r.Address) != 42 {
		return errors.ErrInvalidAddress
	}

	return nil
}

// RefreshContractSourceCodeRequest represents the request for refreshing contract source code
type RefreshContractSourceCodeRequest struct {
	Address string `json:"address" binding:"required"`
	Network string `json:"network" binding:"required"`
}

// Validate validates the request
func (r *RefreshContractSourceCodeRequest) Validate() error {
	if r.Address == "" {
		return errors.ErrInvalidAddress
	}

	if !strings.HasPrefix(r.Address, "0x") {
		return errors.ErrInvalidAddress
	}

	if len(r.Address) != 42 {
		return errors.ErrInvalidAddress
	}

	return nil
}
