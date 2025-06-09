package contracts

import "time"

type GetContractSourceCodeResponse struct {
	Address    string     `json:"address"`
	SourceCode string     `json:"source_code"`
	SourceHash string     `json:"source_hash"`
	Network    string     `json:"network"`
	CachedAt   *time.Time `json:"cached_at,omitempty"`
}

type ContractResponse struct {
	Address         string    `json:"address"`
	Network         string    `json:"network"`
	SourceCode      string    `json:"source_code"`
	ContractName    string    `json:"contract_name"`
	CompilerVersion string    `json:"compiler_version"`
	SourceHash      string    `json:"source_hash"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
