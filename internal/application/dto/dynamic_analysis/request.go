package contracts

import "errors"

var (
	ErrInvalidAddress = errors.New("invalid contract address")
)

// GetContractSourceCodeRequest represents the request for getting contract source code
type GetContractSourceCodeRequest struct {
	Address string `json:"address" validate:"required"`
	Network string `json:"network"`
}

// Validate validates the request
func (r *GetContractSourceCodeRequest) Validate() error {
	if r.Address == "" {
		return ErrInvalidAddress
	}

	if r.Network == "" {
		r.Network = "ethereum" // Default network
	}

	return nil
}

// RefreshContractSourceCodeRequest represents the request for refreshing contract source code
type RefreshContractSourceCodeRequest struct {
	Address string `json:"address" validate:"required"`
}

// Validate validates the request
func (r *RefreshContractSourceCodeRequest) Validate() error {
	if r.Address == "" {
		return ErrInvalidAddress
	}
	return nil
}
