package entities

import (
	"time"

	"github.com/Azzurriii/slythr-go-backend/internal/domain/valueobjects"
	"github.com/google/uuid"
)

type StaticAnalysis struct {
	ID            StaticAnalysisID `gorm:"type:uuid;primaryKey" json:"id"`
	SourceHash    string           `gorm:"not null;size:64;index" json:"source_hash"`
	SlitherOutput string           `gorm:"type:jsonb;not null" json:"slither_output"`
	CreatedAt     time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
}

type StaticAnalysisID uuid.UUID

func NewStaticAnalysis(sourceHash, slitherOutput string) (*StaticAnalysis, error) {
	// Validate source hash
	sourceHashVO, err := valueobjects.NewSourceHash(sourceHash)
	if err != nil {
		return nil, err
	}

	// Validate analysis results
	slitherOutputVO, err := valueobjects.NewAnalysisResults(slitherOutput)
	if err != nil {
		return nil, err
	}

	return &StaticAnalysis{
		ID:            StaticAnalysisID(uuid.New()),
		SourceHash:    sourceHashVO.Value(),
		SlitherOutput: slitherOutputVO.Value(),
	}, nil
}

func (s *StaticAnalysis) GetID() StaticAnalysisID {
	return s.ID
}

func (s *StaticAnalysis) GetSourceHash() valueobjects.SourceHash {
	sourceHash, _ := valueobjects.NewSourceHash(s.SourceHash)
	return sourceHash
}

func (s *StaticAnalysis) GetResponse() valueobjects.AnalysisResults {
	results, _ := valueobjects.NewAnalysisResults(s.SlitherOutput)
	return results
}

func (s *StaticAnalysis) IsValid() bool {
	if s.ID == StaticAnalysisID(uuid.Nil) {
		return false
	}

	sourceHash, err := valueobjects.NewSourceHash(s.SourceHash)
	if err != nil {
		return false
	}

	results, err := valueobjects.NewAnalysisResults(s.SlitherOutput)
	if err != nil {
		return false
	}

	return sourceHash.IsValid() && results.IsValid()
}

func (s *StaticAnalysis) HasResults() bool {
	results := s.GetResponse()
	return results.HasResults()
}

func (id StaticAnalysisID) String() string {
	return uuid.UUID(id).String()
}

func (StaticAnalysis) TableName() string {
	return "static_analysis"
}
