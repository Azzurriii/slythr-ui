package entity

import "time"

type DynamicAnalysis struct {
	ID          int       `db:"id" json:"id"`
	ContractID  int       `db:"contract_id" json:"contract_id"`
	SourceHash  string    `db:"source_hash" json:"source_hash"`
	LLMResponse string    `db:"llm_response" json:"llm_response"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
