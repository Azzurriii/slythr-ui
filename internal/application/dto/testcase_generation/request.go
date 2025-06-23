package testcase_generation

type TestCaseGenerateRequest struct {
	SourceCode    string `json:"source_code" binding:"required" validate:"required,min=1"`
	TestFramework string `json:"test_framework" binding:"required" validate:"required" example:"hardhat"`
	TestLanguage  string `json:"test_language" binding:"required" validate:"required" example:"javascript"`
}
