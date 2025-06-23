package entities

import (
	"encoding/json"
	"time"

	"github.com/Azzurriii/slythr/internal/domain/valueobjects"
	"github.com/google/uuid"
)

type GeneratedTestCases struct {
	ID                         GeneratedTestCasesID `gorm:"type:uuid;primaryKey" json:"id"`
	SourceHash                 string               `gorm:"not null;size:64;index" json:"source_hash"`
	TestCode                   string               `gorm:"type:text" json:"test_code,omitempty"`
	TestFramework              string               `gorm:"size:50" json:"test_framework"`
	TestLanguage               string               `gorm:"size:50" json:"test_language"`
	FileName                   string               `gorm:"size:255" json:"file_name"`
	WarningsAndRecommendations string               `gorm:"type:jsonb" json:"warnings_and_recommendations,omitempty"` // Store as JSON array
	CreatedAt                  time.Time            `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt                  time.Time            `gorm:"autoUpdateTime" json:"updated_at"`
}

type GeneratedTestCasesID uuid.UUID

func NewGeneratedTestCases(
	sourceHash string,
	testCode string,
	testFramework string,
	testLanguage string,
	fileName string,
	warnings []string,
) (*GeneratedTestCases, error) {
	sourceHashVO, err := valueobjects.NewSourceHash(sourceHash)
	if err != nil {
		return nil, err
	}

	// Convert warnings to JSON
	var warningsJSON string
	if len(warnings) > 0 {
		warningsBytes, err := json.Marshal(warnings)
		if err != nil {
			return nil, err
		}
		warningsJSON = string(warningsBytes)
	} else {
		warningsJSON = "[]"
	}

	// Generate new UUID for ID
	newID := uuid.New()

	return &GeneratedTestCases{
		ID:                         GeneratedTestCasesID(newID),
		SourceHash:                 sourceHashVO.Value(),
		TestCode:                   testCode,
		TestFramework:              testFramework,
		TestLanguage:               testLanguage,
		FileName:                   fileName,
		WarningsAndRecommendations: warningsJSON,
	}, nil
}

func (g *GeneratedTestCases) GetID() GeneratedTestCasesID {
	return g.ID
}

func (g *GeneratedTestCases) GetSourceHash() valueobjects.SourceHash {
	sourceHash, _ := valueobjects.NewSourceHash(g.SourceHash)
	return sourceHash
}

func (g *GeneratedTestCases) GetWarnings() []string {
	var warnings []string
	if g.WarningsAndRecommendations != "" {
		json.Unmarshal([]byte(g.WarningsAndRecommendations), &warnings)
	}
	return warnings
}

func (g *GeneratedTestCases) IsValid() bool {
	if g.ID == GeneratedTestCasesID(uuid.Nil) {
		return false
	}

	sourceHash, err := valueobjects.NewSourceHash(g.SourceHash)
	if err != nil {
		return false
	}

	return sourceHash.IsValid() && g.TestFramework != "" && g.TestLanguage != ""
}

func (g *GeneratedTestCases) HasTestCode() bool {
	return g.TestCode != ""
}

func (id GeneratedTestCasesID) String() string {
	return uuid.UUID(id).String()
}

func (GeneratedTestCases) TableName() string {
	return "generated_test_cases"
}
