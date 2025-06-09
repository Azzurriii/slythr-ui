package valueobjects

import (
	"strings"

	"github.com/Azzurriii/slythr-go-backend/internal/domain/errors"
)

type SourceCode struct {
	value string
}

func NewSourceCode(sourceCode string) (SourceCode, error) {
	if err := validateSourceCode(sourceCode); err != nil {
		return SourceCode{}, err
	}

	return SourceCode{value: sourceCode}, nil
}

func (sc SourceCode) String() string {
	return sc.value
}

func (sc SourceCode) Value() string {
	return sc.value
}

func (sc SourceCode) IsValid() bool {
	return validateSourceCode(sc.value) == nil
}

func (sc SourceCode) IsZero() bool {
	return strings.TrimSpace(sc.value) == ""
}

func (sc SourceCode) Length() int {
	return len(sc.value)
}

func (sc SourceCode) HasContent() bool {
	return !sc.IsZero()
}

func validateSourceCode(sourceCode string) error {
	if strings.TrimSpace(sourceCode) == "" {
		return errors.ErrEmptySourceCode
	}

	return nil
}
