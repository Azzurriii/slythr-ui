package valueobjects

import (
	"strings"

	"github.com/Azzurriii/slythr/internal/domain/errors"
)

type SourceHash struct {
	value string
}

func NewSourceHash(hash string) (SourceHash, error) {
	if err := validateSourceHash(hash); err != nil {
		return SourceHash{}, err
	}

	return SourceHash{value: strings.ToLower(hash)}, nil
}

func (sh SourceHash) String() string {
	return sh.value
}

func (sh SourceHash) Value() string {
	return sh.value
}

func (sh SourceHash) IsValid() bool {
	return validateSourceHash(sh.value) == nil
}

func (sh SourceHash) IsZero() bool {
	return sh.value == ""
}

func (sh SourceHash) Equals(other SourceHash) bool {
	return sh.value == other.value
}

func validateSourceHash(hash string) error {
	if hash == "" {
		return errors.ErrInvalidSourceHash
	}

	if len(hash) != 64 {
		return errors.ErrInvalidSourceHash
	}

	for i := 0; i < len(hash); i++ {
		c := hash[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return errors.ErrInvalidSourceHash
		}
	}

	return nil
}
