package entities

import (
	"time"

	"github.com/Azzurriii/slythr-go-backend/internal/domain/valueobjects"
)

// DynamicAnalysis represents a dynamic analysis result with LLM response in the domain
type DynamicAnalysis struct {
	ID          DynamicAnalysisID `gorm:"primaryKey" json:"id"`
	ContractID  ContractID        `gorm:"not null;index" json:"contract_id"`
	SourceHash  string            `gorm:"not null;size:64;index" json:"source_hash"`
	LLMResponse string            `gorm:"type:text;not null" json:"llm_response"`
	CreatedAt   time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time         `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Contract *Contract `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
}

// DynamicAnalysisID represents the unique identifier for a dynamic analysis
type DynamicAnalysisID uint

// NewDynamicAnalysis creates a new dynamic analysis with validation
func NewDynamicAnalysis(contractID uint, sourceHash, llmResponse string) (*DynamicAnalysis, error) {
	// Validate source hash
	sourceHashVO, err := valueobjects.NewSourceHash(sourceHash)
	if err != nil {
		return nil, err
	}

	// Validate LLM response
	llmResponseVO, err := valueobjects.NewLLMResponse(llmResponse)
	if err != nil {
		return nil, err
	}

	return &DynamicAnalysis{
		ContractID:  ContractID(contractID),
		SourceHash:  sourceHashVO.Value(),
		LLMResponse: llmResponseVO.Value(),
	}, nil
}

// GetID returns the dynamic analysis ID
func (d *DynamicAnalysis) GetID() DynamicAnalysisID {
	return d.ID
}

// GetContractID returns the contract ID
func (d *DynamicAnalysis) GetContractID() ContractID {
	return d.ContractID
}

// GetSourceHash returns the source hash as value object
func (d *DynamicAnalysis) GetSourceHash() valueobjects.SourceHash {
	sourceHash, _ := valueobjects.NewSourceHash(d.SourceHash)
	return sourceHash
}

// GetLLMResponse returns the LLM response as value object
func (d *DynamicAnalysis) GetLLMResponse() valueobjects.LLMResponse {
	llmResponse, _ := valueobjects.NewLLMResponse(d.LLMResponse)
	return llmResponse
}

// IsValid checks if the dynamic analysis is valid
func (d *DynamicAnalysis) IsValid() bool {
	if d.ContractID == 0 {
		return false
	}

	sourceHash, err := valueobjects.NewSourceHash(d.SourceHash)
	if err != nil {
		return false
	}

	llmResponse, err := valueobjects.NewLLMResponse(d.LLMResponse)
	if err != nil {
		return false
	}

	return sourceHash.IsValid() && llmResponse.IsValid()
}

// HasResponse checks if the analysis has a valid LLM response
func (d *DynamicAnalysis) HasResponse() bool {
	llmResponse := d.GetLLMResponse()
	return llmResponse.HasContent()
}

// TableName returns the table name for GORM
func (DynamicAnalysis) TableName() string {
	return "dynamic_analyses"
}
