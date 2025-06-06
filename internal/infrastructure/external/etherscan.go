package infrastructures

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	config "github.com/Azzurriii/slythr-go-backend/config"
)

const (
	etherscanAPIBaseURL = "https://api.etherscan.io/v2/api"
	defaultTimeout      = 30 * time.Second
)

var networkToChainID = map[string]string{
	"ethereum":  "1",
	"polygon":   "137",
	"bsc":       "56",
	"base":      "8453",
	"arbitrum":  "42161",
	"avalanche": "43114",
	"optimism":  "10",
	"gnosis":    "100",
	"fantom":    "250",
	"celo":      "42220",
}

// EtherscanClient handles interactions with the Etherscan API
type EtherscanClient struct {
	apiKey     string
	httpClient HTTPClient
}

// HTTPClient interface for easier testing and mocking
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewEtherscanClient creates a new Etherscan API client
func NewEtherscanClient(cfg *config.EtherscanConfig) *EtherscanClient {
	return &EtherscanClient{
		apiKey: cfg.APIKey,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// NewEtherscanClientWithHTTPClient creates a new client with custom HTTP client
func NewEtherscanClientWithHTTPClient(apiKey string, httpClient HTTPClient) *EtherscanClient {
	return &EtherscanClient{
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

// EtherscanError represents an error from the Etherscan API
type EtherscanError struct {
	Status  string
	Message string
	Body    string
}

func (e EtherscanError) Error() string {
	return fmt.Sprintf("etherscan API error: %s - %s", e.Message, e.Body)
}

// ContractSourceResponse represents the response structure from Etherscan API
type ContractSourceResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

// ContractInfo contains detailed information about a smart contract
type ContractInfo struct {
	SourceCode           string `json:"SourceCode"`
	ABI                  string `json:"ABI"`
	ContractName         string `json:"ContractName"`
	CompilerVersion      string `json:"CompilerVersion"`
	Runs                 string `json:"Runs"`
	ConstructorArguments string `json:"ConstructorArguments"`
	EVMVersion           string `json:"EVMVersion"`
	Library              string `json:"Library"`
	LicenseType          string `json:"LicenseType"`
	Proxy                string `json:"Proxy"`
	Implementation       string `json:"Implementation"`
	SwarmSource          string `json:"SwarmSource"`
}

// MultiFileSource represents the structure for multi-file contracts
type MultiFileSource struct {
	Sources map[string]FileContent `json:"sources"`
}

// FileContent represents individual file content
type FileContent struct {
	Content string `json:"content"`
}

// GetContractSourceCode retrieves and processes the source code for a given contract address
func (c *EtherscanClient) GetContractSourceCode(ctx context.Context, address string, network string) (string, error) {
	contractInfo, err := c.fetchContractInfo(ctx, address, network)
	if err != nil {
		return "", err
	}

	return c.processSourceCode(contractInfo.SourceCode)
}

// GetContractDetails retrieves comprehensive contract information
func (c *EtherscanClient) GetContractDetails(ctx context.Context, address string, network string) (*ContractInfo, error) {
	contractInfo, err := c.fetchContractInfo(ctx, address, network)
	if err != nil {
		return nil, err
	}

	processedSourceCode, err := c.processSourceCode(contractInfo.SourceCode)
	if err != nil {
		return nil, err
	}

	// Update the source code with processed version
	contractInfo.SourceCode = processedSourceCode
	return &contractInfo, nil
}

// fetchContractInfo makes the API request and returns the contract information
func (c *EtherscanClient) fetchContractInfo(ctx context.Context, address string, network string) (ContractInfo, error) {
	chainID, ok := networkToChainID[network]
	if !ok {
		return ContractInfo{}, fmt.Errorf("unsupported network: %s", network)
	}
	reqURL := fmt.Sprintf("%s?module=contract&action=getsourcecode&address=%s&chainid=%s&apikey=%s",
		etherscanAPIBaseURL, address, chainID, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return ContractInfo{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ContractInfo{}, fmt.Errorf("failed to make request to Etherscan API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ContractInfo{}, fmt.Errorf("etherscan API returned non-OK status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ContractInfo{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var rawResponse ContractSourceResponse
	if err := json.Unmarshal(body, &rawResponse); err != nil {
		return ContractInfo{}, fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	if rawResponse.Status != "1" {
		return ContractInfo{}, EtherscanError{
			Status:  rawResponse.Status,
			Message: rawResponse.Message,
			Body:    string(body),
		}
	}

	var contractInfoResult []ContractInfo
	if err := json.Unmarshal(rawResponse.Result, &contractInfoResult); err != nil {
		// This can happen if result is a string (e.g., "Contract source code not verified")
		return ContractInfo{}, EtherscanError{
			Status:  "0",
			Message: "Contract source code not found or not verified",
			Body:    string(rawResponse.Result),
		}
	}

	if len(contractInfoResult) == 0 {
		return ContractInfo{}, fmt.Errorf("no contract information found for address: %s", address)
	}

	return contractInfoResult[0], nil
}

// processSourceCode handles multi-file source code processing
func (c *EtherscanClient) processSourceCode(sourceCode string) (string, error) {
	if len(sourceCode) == 0 {
		return "", fmt.Errorf("empty source code")
	}

	// Check if it's a single file (not JSON)
	if sourceCode[0] != '{' {
		return sourceCode, nil
	}

	// Try to parse as multi-file source
	processedSource, err := c.parseMultiFileSource(sourceCode)
	if err != nil {
		// If parsing fails, return original source code
		return sourceCode, nil
	}

	if processedSource != "" {
		return processedSource, nil
	}

	return sourceCode, nil
}

// parseMultiFileSource attempts to parse and combine multi-file sources
func (c *EtherscanClient) parseMultiFileSource(sourceCode string) (string, error) {
	// Try standard multi-file format first
	var multiFileSource MultiFileSource
	if err := json.Unmarshal([]byte(sourceCode), &multiFileSource); err == nil {
		if len(multiFileSource.Sources) > 0 {
			return c.combineMultiFileSources(multiFileSource.Sources), nil
		}
	}

	// Try alternative format
	var altFormat map[string]FileContent
	if err := json.Unmarshal([]byte(sourceCode), &altFormat); err == nil {
		if len(altFormat) > 0 {
			return c.combineAlternativeFormat(altFormat), nil
		}
	}

	return "", fmt.Errorf("failed to parse multi-file source")
}

// combineMultiFileSources combines multiple source files into a single string
func (c *EtherscanClient) combineMultiFileSources(sources map[string]FileContent) string {
	var combined string
	for path, fileContent := range sources {
		combined += fmt.Sprintf("// File: %s\n%s\n\n", path, fileContent.Content)
	}
	return combined
}

// combineAlternativeFormat combines alternative format sources
func (c *EtherscanClient) combineAlternativeFormat(sources map[string]FileContent) string {
	var combined string
	for fileName, fileContent := range sources {
		combined += fmt.Sprintf("// File: %s\n%s\n\n", fileName, fileContent.Content)
	}
	return combined
}
