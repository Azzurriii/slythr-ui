package valueobjects

import (
	"encoding/json"
	"strings"

	"github.com/Azzurriii/slythr-go-backend/internal/domain/errors"
)

type AnalysisResults struct {
	value string
}

func NewAnalysisResults(results string) (AnalysisResults, error) {
	if err := validateAnalysisResults(results); err != nil {
		return AnalysisResults{}, err
	}

	return AnalysisResults{value: results}, nil
}

func (ar AnalysisResults) String() string {
	return ar.value
}

func (ar AnalysisResults) Value() string {
	return ar.value
}

func (ar AnalysisResults) IsValid() bool {
	return validateAnalysisResults(ar.value) == nil
}

func (ar AnalysisResults) IsZero() bool {
	return strings.TrimSpace(ar.value) == ""
}

func (ar AnalysisResults) HasResults() bool {
	trimmed := strings.TrimSpace(ar.value)
	return trimmed != "" && trimmed != "{}" && trimmed != "null"
}

func (ar AnalysisResults) IsValidJSON() bool {
	var temp interface{}
	return json.Unmarshal([]byte(ar.value), &temp) == nil
}

func validateAnalysisResults(results string) error {
	if strings.TrimSpace(results) == "" {
		return errors.ErrInvalidAnalysisResults
	}

	var temp interface{}
	if err := json.Unmarshal([]byte(results), &temp); err != nil {
		return errors.ErrInvalidAnalysisResults
	}

	return nil
}
