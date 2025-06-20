package valueobjects

import (
	"strings"

	"github.com/Azzurriii/slythr/internal/domain/constants"
	"github.com/Azzurriii/slythr/internal/domain/errors"
)

type Network struct {
	value string
}

func NewNetwork(network string) (Network, error) {
	if err := validateNetwork(network); err != nil {
		return Network{}, err
	}

	return Network{value: strings.ToLower(network)}, nil
}

func (n Network) String() string {
	return n.value
}

func (n Network) Value() string {
	return n.value
}

func (n Network) IsValid() bool {
	return validateNetwork(n.value) == nil
}

func (n Network) IsZero() bool {
	return n.value == ""
}

func (n Network) Equals(other Network) bool {
	return n.value == other.value
}

func (n Network) GetChainID() (string, bool) {
	return constants.GetChainID(n.value)
}

func validateNetwork(network string) error {
	if network == "" {
		return errors.ErrInvalidNetwork
	}

	if !constants.IsValidNetwork(strings.ToLower(network)) {
		return errors.ErrInvalidNetwork
	}

	return nil
}
