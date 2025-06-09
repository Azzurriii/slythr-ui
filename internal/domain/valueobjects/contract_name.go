package valueobjects

import (
	"strings"

	"github.com/Azzurriii/slythr-go-backend/internal/domain/errors"
)

type ContractName struct {
	value string
}

func NewContractName(name string) (ContractName, error) {
	if err := validateContractName(name); err != nil {
		return ContractName{}, err
	}

	return ContractName{value: strings.TrimSpace(name)}, nil
}

func (cn ContractName) String() string {
	return cn.value
}

func (cn ContractName) Value() string {
	return cn.value
}

func (cn ContractName) IsValid() bool {
	return validateContractName(cn.value) == nil
}

func (cn ContractName) IsZero() bool {
	return strings.TrimSpace(cn.value) == ""
}

func (cn ContractName) Length() int {
	return len(cn.value)
}

func validateContractName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return errors.ErrInvalidContractName
	}

	if len(trimmed) > 255 {
		return errors.ErrInvalidContractName
	}

	return nil
}
