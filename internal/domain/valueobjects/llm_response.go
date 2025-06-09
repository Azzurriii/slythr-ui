package valueobjects

import (
	"strings"

	"github.com/Azzurriii/slythr-go-backend/internal/domain/errors"
)

type LLMResponse struct {
	value string
}

func NewLLMResponse(response string) (LLMResponse, error) {
	if err := validateLLMResponse(response); err != nil {
		return LLMResponse{}, err
	}

	return LLMResponse{value: response}, nil
}

func (lr LLMResponse) String() string {
	return lr.value
}

func (lr LLMResponse) Value() string {
	return lr.value
}

func (lr LLMResponse) IsValid() bool {
	return validateLLMResponse(lr.value) == nil
}

func (lr LLMResponse) IsZero() bool {
	return strings.TrimSpace(lr.value) == ""
}

func (lr LLMResponse) HasContent() bool {
	return !lr.IsZero()
}

func (lr LLMResponse) Length() int {
	return len(lr.value)
}

func validateLLMResponse(response string) error {
	if strings.TrimSpace(response) == "" {
		return errors.ErrInvalidLLMResponse
	}

	return nil
}
