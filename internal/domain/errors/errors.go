package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// Domain error codes
const (
	ErrCodeContractNotFound         = "CONTRACT_NOT_FOUND"
	ErrCodeInvalidAddress           = "INVALID_ADDRESS"
	ErrCodeInvalidNetwork           = "INVALID_NETWORK"
	ErrCodeEmptySourceCode          = "EMPTY_SOURCE_CODE"
	ErrCodeInvalidContractName      = "INVALID_CONTRACT_NAME"
	ErrCodeInvalidCompilerVersion   = "INVALID_COMPILER_VERSION"
	ErrCodeInvalidSourceHash        = "INVALID_SOURCE_HASH"
	ErrCodeContractAlreadyExists    = "CONTRACT_ALREADY_EXISTS"
	ErrCodeInvalidContractReference = "INVALID_CONTRACT_REFERENCE"
	ErrCodeStaticAnalysisNotFound   = "STATIC_ANALYSIS_NOT_FOUND"
	ErrCodeInvalidAnalysisResults   = "INVALID_ANALYSIS_RESULTS"
	ErrCodeDynamicAnalysisNotFound  = "DYNAMIC_ANALYSIS_NOT_FOUND"
	ErrCodeInvalidLLMResponse       = "INVALID_LLM_RESPONSE"
)

// Domain errors for contracts
var (
	ErrContractNotFound         = &DomainError{Code: ErrCodeContractNotFound, Message: "contract not found"}
	ErrInvalidAddress           = &DomainError{Code: ErrCodeInvalidAddress, Message: "invalid contract address"}
	ErrInvalidNetwork           = &DomainError{Code: ErrCodeInvalidNetwork, Message: "invalid network"}
	ErrEmptySourceCode          = &DomainError{Code: ErrCodeEmptySourceCode, Message: "empty source code"}
	ErrInvalidContractName      = &DomainError{Code: ErrCodeInvalidContractName, Message: "invalid contract name"}
	ErrInvalidCompilerVersion   = &DomainError{Code: ErrCodeInvalidCompilerVersion, Message: "invalid compiler version"}
	ErrInvalidSourceHash        = &DomainError{Code: ErrCodeInvalidSourceHash, Message: "invalid source hash"}
	ErrContractAlreadyExists    = &DomainError{Code: ErrCodeContractAlreadyExists, Message: "contract already exists"}
	ErrInvalidContractReference = &DomainError{Code: ErrCodeInvalidContractReference, Message: "invalid contract reference"}
)

// Domain errors for static analysis
var (
	ErrStaticAnalysisNotFound = &DomainError{Code: ErrCodeStaticAnalysisNotFound, Message: "static analysis not found"}
	ErrInvalidAnalysisResults = &DomainError{Code: ErrCodeInvalidAnalysisResults, Message: "invalid analysis results"}
)

// Domain errors for dynamic analysis
var (
	ErrDynamicAnalysisNotFound = &DomainError{Code: ErrCodeDynamicAnalysisNotFound, Message: "dynamic analysis not found"}
	ErrInvalidLLMResponse      = &DomainError{Code: ErrCodeInvalidLLMResponse, Message: "invalid LLM response"}
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NewErrorResponse creates a new error response from an error
func NewErrorResponse(err error) *ErrorResponse {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return &ErrorResponse{
			Error:   domainErr.Error(),
			Code:    domainErr.Code,
			Message: getHTTPMessage(domainErr),
			Details: getErrorDetails(domainErr),
		}
	}

	// Fallback for non-domain errors
	return &ErrorResponse{
		Error:   err.Error(),
		Code:    "INTERNAL_ERROR",
		Message: "Internal server error",
	}
}

// GetHTTPStatusCode returns appropriate HTTP status code for domain errors
func GetHTTPStatusCode(err error) int {
	var domainErr *DomainError
	if !errors.As(err, &domainErr) {
		return http.StatusInternalServerError
	}

	switch domainErr.Code {
	case ErrCodeContractNotFound, ErrCodeStaticAnalysisNotFound, ErrCodeDynamicAnalysisNotFound:
		return http.StatusNotFound
	case ErrCodeInvalidAddress, ErrCodeInvalidNetwork, ErrCodeEmptySourceCode,
		ErrCodeInvalidContractName, ErrCodeInvalidCompilerVersion, ErrCodeInvalidSourceHash,
		ErrCodeInvalidContractReference, ErrCodeInvalidAnalysisResults, ErrCodeInvalidLLMResponse:
		return http.StatusBadRequest
	case ErrCodeContractAlreadyExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// Helper functions

func getHTTPMessage(domainErr *DomainError) string {
	switch domainErr.Code {
	case ErrCodeContractNotFound:
		return "Contract not found"
	case ErrCodeInvalidAddress:
		return "Invalid contract address format"
	case ErrCodeInvalidNetwork:
		return "Invalid network"
	case ErrCodeEmptySourceCode:
		return "Contract source code is empty"
	case ErrCodeInvalidContractName:
		return "Invalid contract name"
	case ErrCodeInvalidCompilerVersion:
		return "Invalid compiler version"
	case ErrCodeInvalidSourceHash:
		return "Invalid source hash"
	case ErrCodeContractAlreadyExists:
		return "Contract already exists"
	case ErrCodeInvalidContractReference:
		return "Invalid contract reference"
	case ErrCodeStaticAnalysisNotFound:
		return "Static analysis not found"
	case ErrCodeInvalidAnalysisResults:
		return "Invalid analysis results"
	case ErrCodeDynamicAnalysisNotFound:
		return "Dynamic analysis not found"
	case ErrCodeInvalidLLMResponse:
		return "Invalid LLM response"
	default:
		return "Internal server error"
	}
}

func getErrorDetails(domainErr *DomainError) string {
	if domainErr.Err != nil {
		return domainErr.Err.Error()
	}
	return ""
}

// Error factory functions for better error creation

// NewContractNotFoundError creates a new contract not found error
func NewContractNotFoundError(address, network string) error {
	return &DomainError{
		Code:    ErrCodeContractNotFound,
		Message: fmt.Sprintf("contract not found: %s on %s", address, network),
	}
}

// NewInvalidAddressError creates a new invalid address error
func NewInvalidAddressError(address string) error {
	return &DomainError{
		Code:    ErrCodeInvalidAddress,
		Message: fmt.Sprintf("invalid contract address: %s", address),
	}
}

// NewContractAlreadyExistsError creates a new contract already exists error
func NewContractAlreadyExistsError(address, network string) error {
	return &DomainError{
		Code:    ErrCodeContractAlreadyExists,
		Message: fmt.Sprintf("contract already exists: %s on %s", address, network),
	}
}
