package valueobjects

import (
	"strings"

	"github.com/Azzurriii/slythr/internal/domain/errors"
)

type ContractAddress struct {
	value string
}

func NewContractAddress(address string) (ContractAddress, error) {
	if err := validateAddress(address); err != nil {
		return ContractAddress{}, err
	}

	return ContractAddress{value: address}, nil
}

func (ca ContractAddress) String() string {
	return ca.value
}

func (ca ContractAddress) Value() string {
	return ca.value
}

func (ca ContractAddress) IsValid() bool {
	return validateAddress(ca.value) == nil
}

func (ca ContractAddress) IsZero() bool {
	return ca.value == ""
}

func (ca ContractAddress) Equals(other ContractAddress) bool {
	return strings.EqualFold(ca.value, other.value)
}

func validateAddress(address string) error {
	if address == "" {
		return errors.ErrInvalidAddress
	}

	if !strings.HasPrefix(address, "0x") {
		return errors.ErrInvalidAddress
	}

	if len(address) != 42 {
		return errors.ErrInvalidAddress
	}

	for i := 2; i < len(address); i++ {
		c := address[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return errors.ErrInvalidAddress
		}
	}

	return nil
}
