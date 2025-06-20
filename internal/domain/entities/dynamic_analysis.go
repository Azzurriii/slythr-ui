package entities

import (
	"time"

	"github.com/Azzurriii/slythr/internal/domain/valueobjects"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DynamicAnalysis struct {
	ID          DynamicAnalysisID `gorm:"type:uuid;primaryKey" json:"id"`
	SourceHash  string            `gorm:"not null;size:64;index" json:"source_hash"`
	LLMResponse string            `gorm:"type:text;not null" json:"llm_response"`
	CreatedAt   time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"deleted_at,omitempty"`
}

type DynamicAnalysisID uuid.UUID

func NewDynamicAnalysis(sourceHash, llmResponse string) (*DynamicAnalysis, error) {
	sourceHashVO, err := valueobjects.NewSourceHash(sourceHash)
	if err != nil {
		return nil, err
	}

	llmResponseVO, err := valueobjects.NewLLMResponse(llmResponse)
	if err != nil {
		return nil, err
	}

	return &DynamicAnalysis{
		ID:          DynamicAnalysisID(uuid.New()),
		SourceHash:  sourceHashVO.Value(),
		LLMResponse: llmResponseVO.Value(),
	}, nil
}

func (d *DynamicAnalysis) GetID() DynamicAnalysisID {
	return d.ID
}

func (d *DynamicAnalysis) GetSourceHash() valueobjects.SourceHash {
	sourceHash, _ := valueobjects.NewSourceHash(d.SourceHash)
	return sourceHash
}

func (d *DynamicAnalysis) GetLLMResponse() valueobjects.LLMResponse {
	llmResponse, _ := valueobjects.NewLLMResponse(d.LLMResponse)
	return llmResponse
}

func (d *DynamicAnalysis) IsValid() bool {
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

func (d *DynamicAnalysis) HasResponse() bool {
	llmResponse := d.GetLLMResponse()
	return llmResponse.HasContent()
}

func (DynamicAnalysis) TableName() string {
	return "dynamic_analysis"
}
