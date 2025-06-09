package entities

import (
	"time"

	"github.com/Azzurriii/slythr-go-backend/internal/domain/valueobjects"
)

// StaticAnalysis represents a static analysis result in the domain
type StaticAnalysis struct {
	ID         StaticAnalysisID `gorm:"primaryKey" json:"id"`
	ContractID ContractID       `gorm:"not null;index" json:"contract_id"`
	SourceHash string           `gorm:"not null;size:64;index" json:"source_hash"`
	Results    string           `gorm:"type:jsonb;not null" json:"results"`
	CreatedAt  time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time        `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Contract *Contract `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
}

// StaticAnalysisID represents the unique identifier for a static analysis
type StaticAnalysisID uint

// NewStaticAnalysis creates a new static analysis with validation
func NewStaticAnalysis(contractID uint, sourceHash, results string) (*StaticAnalysis, error) {
	// Validate source hash
	sourceHashVO, err := valueobjects.NewSourceHash(sourceHash)
	if err != nil {
		return nil, err
	}

	// Validate analysis results
	resultsVO, err := valueobjects.NewAnalysisResults(results)
	if err != nil {
		return nil, err
	}

	return &StaticAnalysis{
		ContractID: ContractID(contractID),
		SourceHash: sourceHashVO.Value(),
		Results:    resultsVO.Value(),
	}, nil
}

// GetID returns the static analysis ID
func (s *StaticAnalysis) GetID() StaticAnalysisID {
	return s.ID
}

// GetContractID returns the contract ID
func (s *StaticAnalysis) GetContractID() ContractID {
	return s.ContractID
}

// GetSourceHash returns the source hash as value object
func (s *StaticAnalysis) GetSourceHash() valueobjects.SourceHash {
	sourceHash, _ := valueobjects.NewSourceHash(s.SourceHash)
	return sourceHash
}

// GetResults returns the analysis results as value object
func (s *StaticAnalysis) GetResults() valueobjects.AnalysisResults {
	results, _ := valueobjects.NewAnalysisResults(s.Results)
	return results
}

// IsValid checks if the static analysis is valid
func (s *StaticAnalysis) IsValid() bool {
	if s.ContractID == 0 {
		return false
	}

	sourceHash, err := valueobjects.NewSourceHash(s.SourceHash)
	if err != nil {
		return false
	}

	results, err := valueobjects.NewAnalysisResults(s.Results)
	if err != nil {
		return false
	}

	return sourceHash.IsValid() && results.IsValid()
}

// HasResults checks if the analysis has valid results
func (s *StaticAnalysis) HasResults() bool {
	results := s.GetResults()
	return results.HasResults()
}

// TableName returns the table name for GORM
func (StaticAnalysis) TableName() string {
	return "static_analyses"
}
