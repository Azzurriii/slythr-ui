package valueobjects

import (
	"strings"

	"github.com/Azzurriii/slythr/internal/domain/errors"
)

type CompilerVersion struct {
	value string
}

func NewCompilerVersion(version string) (CompilerVersion, error) {
	if err := validateCompilerVersion(version); err != nil {
		return CompilerVersion{}, err
	}

	return CompilerVersion{value: strings.TrimSpace(version)}, nil
}

func (cv CompilerVersion) String() string {
	return cv.value
}

func (cv CompilerVersion) Value() string {
	return cv.value
}

func (cv CompilerVersion) IsValid() bool {
	return validateCompilerVersion(cv.value) == nil
}

func (cv CompilerVersion) IsZero() bool {
	return strings.TrimSpace(cv.value) == ""
}

func validateCompilerVersion(version string) error {
	trimmed := strings.TrimSpace(version)
	if trimmed == "" {
		return errors.ErrInvalidCompilerVersion
	}

	if len(trimmed) > 50 {
		return errors.ErrInvalidCompilerVersion
	}

	return nil
}
