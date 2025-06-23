package testcase_generation

import "time"

type TestCaseGenerateResponse struct {
	Success                    bool      `json:"success" example:"true"`
	Message                    string    `json:"message,omitempty" example:""`
	TestCode                   string    `json:"test_code,omitempty"`
	TestFramework              string    `json:"test_framework" example:"hardhat"`
	TestLanguage               string    `json:"test_language" example:"javascript"`
	FileName                   string    `json:"file_name" example:"MyContract.test.js"`
	SourceHash                 string    `json:"source_hash" example:"1234567890"`
	WarningsAndRecommendations []string  `json:"warnings_and_recommendations" example:"[]"`
	GeneratedAt                time.Time `json:"generated_at" example:"2025-06-17T16:19:00.579573+07:00"`
}
