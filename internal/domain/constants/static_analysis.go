package entity

import "time"

type StaticAnalysis struct {
	ID         int       `db:"id" json:"id"`
	ContractID int       `db:"contract_id" json:"contract_id"`
	SourceHash string    `db:"source_hash" json:"source_hash"`
	Results    string    `db:"results" json:"results"` // JSONB
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
