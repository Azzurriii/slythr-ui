package errors

import (
	"errors"
	"net/http"
)

var (
	ErrContractNotFound = errors.New("contract not found")
	ErrInvalidAddress   = errors.New("invalid contract address")
	ErrEmptySourceCode  = errors.New("empty source code")
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewErrorResponse creates a new error response from an error
func NewErrorResponse(err error) *ErrorResponse {
	var code int
	var message string

	switch {
	case errors.Is(err, ErrContractNotFound):
		code = http.StatusNotFound
		message = "Contract not found"
	case errors.Is(err, ErrInvalidAddress):
		code = http.StatusBadRequest
		message = "Invalid contract address"
	case errors.Is(err, ErrEmptySourceCode):
		code = http.StatusNotFound
		message = "Contract source code is empty"
	default:
		code = http.StatusInternalServerError
		message = "Internal server error"
	}

	return &ErrorResponse{
		Error:   err.Error(),
		Code:    code,
		Message: message,
	}
}
