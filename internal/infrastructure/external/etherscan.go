package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	config "github.com/Azzurriii/slythr-go-backend/config"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/constants"
)

const (
	etherscanAPIBaseURL = "https://api.etherscan.io/v2/api"
	defaultTimeout      = 30 * time.Second
	maxResponseSize     = 50 * 1024 * 1024
)

var (
	stringBuilderPool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
)

type EtherscanService interface {
	GetContractSourceCode(ctx context.Context, address string, network string) (string, error)
	GetContractDetails(ctx context.Context, address string, network string) (*ContractInfo, error)
}

type EtherscanClient struct {
	apiKey     string
	httpClient *http.Client
}

type EtherscanError struct {
	Status  string
	Message string
	Body    string
}

func (e EtherscanError) Error() string {
	return fmt.Sprintf("etherscan API error: %s - %s", e.Message, e.Body)
}

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

type apiResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

// Simplified source structures for parsing
type sourceFile struct {
	Content string `json:"content"`
}

func NewEtherscanClient(cfg *config.EtherscanConfig) *EtherscanClient {
	return &EtherscanClient{
		apiKey: cfg.APIKey,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     30 * time.Second,
				DisableCompression:  false,
			},
		},
	}
}

func NewEtherscanClientWithHTTPClient(apiKey string, httpClient *http.Client) *EtherscanClient {
	return &EtherscanClient{
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (c *EtherscanClient) GetContractSourceCode(ctx context.Context, address string, network string) (string, error) {
	contractInfo, err := c.fetchContractInfo(ctx, address, network)
	if err != nil {
		return "", err
	}

	return c.extractMainSource(contractInfo.SourceCode, contractInfo.ContractName)
}

func (c *EtherscanClient) GetContractDetails(ctx context.Context, address string, network string) (*ContractInfo, error) {
	contractInfo, err := c.fetchContractInfo(ctx, address, network)
	if err != nil {
		return nil, err
	}

	processedSourceCode, err := c.extractMainSource(contractInfo.SourceCode, contractInfo.ContractName)
	if err != nil {
		return nil, err
	}

	contractInfo.SourceCode = processedSourceCode
	return &contractInfo, nil
}

func (c *EtherscanClient) fetchContractInfo(ctx context.Context, address string, network string) (ContractInfo, error) {
	chainID, ok := constants.GetChainID(network)
	if !ok {
		return ContractInfo{}, fmt.Errorf("unsupported network: %s", network)
	}

	url := c.buildAPIURL(address, chainID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ContractInfo{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "slythr-backend/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ContractInfo{}, fmt.Errorf("failed to make request to Etherscan API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ContractInfo{}, fmt.Errorf("etherscan API returned status %d", resp.StatusCode)
	}

	limitedReader := io.LimitReader(resp.Body, maxResponseSize)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return ContractInfo{}, fmt.Errorf("failed to read response body: %w", err)
	}

	return c.parseAPIResponse(body, address)
}

func (c *EtherscanClient) buildAPIURL(address, chainID string) string {
	builder := stringBuilderPool.Get().(*strings.Builder)
	defer func() {
		builder.Reset()
		stringBuilderPool.Put(builder)
	}()

	builder.WriteString(etherscanAPIBaseURL)
	builder.WriteString("?module=contract&action=getsourcecode&address=")
	builder.WriteString(address)
	builder.WriteString("&chainid=")
	builder.WriteString(chainID)
	builder.WriteString("&apikey=")
	builder.WriteString(c.apiKey)

	return builder.String()
}

func (c *EtherscanClient) parseAPIResponse(body []byte, address string) (ContractInfo, error) {
	var response apiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return ContractInfo{}, fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	if response.Status != "1" {
		return ContractInfo{}, EtherscanError{
			Status:  response.Status,
			Message: response.Message,
			Body:    string(body),
		}
	}

	var contractInfoResult []ContractInfo
	if err := json.Unmarshal(response.Result, &contractInfoResult); err != nil {
		return ContractInfo{}, EtherscanError{
			Status:  "0",
			Message: "Contract source code not found or not verified",
			Body:    string(response.Result),
		}
	}

	if len(contractInfoResult) == 0 {
		return ContractInfo{}, fmt.Errorf("no contract information found for address: %s", address)
	}

	return contractInfoResult[0], nil
}

// extractMainSource efficiently extracts only the main contract source
func (c *EtherscanClient) extractMainSource(sourceCode string, contractName string) (string, error) {
	if len(sourceCode) == 0 {
		return "", fmt.Errorf("empty source code")
	}

	// If not JSON format, return as-is (single file contract)
	if sourceCode[0] != '{' {
		return sourceCode, nil
	}

	// Handle double-encoded JSON
	if len(sourceCode) > 1 && sourceCode[0] == '{' && sourceCode[1] == '{' {
		var jsonString string
		if err := json.Unmarshal([]byte(sourceCode), &jsonString); err == nil {
			sourceCode = jsonString
		} else if len(sourceCode) >= 4 {
			// Try removing extra braces
			sourceCode = sourceCode[1 : len(sourceCode)-1]
		}
	}

	// Use json.RawMessage for efficient parsing without full unmarshaling
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(sourceCode), &raw); err != nil {
		return sourceCode, nil // Return as-is if can't parse
	}

	// Check if it has the standard format with "sources" field
	if sourcesRaw, ok := raw["sources"]; ok {
		return c.extractFromSources(sourcesRaw, contractName), nil
	}

	// Direct format (map of files)
	return c.extractFromRawMap(raw, contractName), nil
}

// extractFromSources handles the standard multi-file format
func (c *EtherscanClient) extractFromSources(sourcesRaw json.RawMessage, contractName string) string {
	var sources map[string]sourceFile
	if err := json.Unmarshal(sourcesRaw, &sources); err != nil {
		return ""
	}

	// Priority order for finding the main contract
	patterns := []func(string) bool{
		// 1. Exact match in contracts directory
		func(path string) bool {
			return strings.HasSuffix(path, "contracts/"+contractName+".sol")
		},
		// 2. Exact filename match
		func(path string) bool {
			return strings.HasSuffix(path, "/"+contractName+".sol") || path == contractName+".sol"
		},
		// 3. Any file in contracts directory (non-library)
		func(path string) bool {
			return strings.Contains(path, "contracts/") &&
				strings.HasSuffix(path, ".sol") &&
				!c.isLibraryPath(path)
		},
		// 4. Any non-library Solidity file
		func(path string) bool {
			return strings.HasSuffix(path, ".sol") && !c.isLibraryPath(path)
		},
	}

	// Try each pattern in order
	for _, pattern := range patterns {
		for path, file := range sources {
			if pattern(path) {
				return file.Content
			}
		}
	}

	// Fallback: return first available file
	for _, file := range sources {
		return file.Content
	}

	return ""
}

// extractFromRawMap handles direct map format
func (c *EtherscanClient) extractFromRawMap(raw map[string]json.RawMessage, contractName string) string {
	// Similar logic but working with RawMessage
	patterns := []func(string) bool{
		func(path string) bool {
			return strings.HasSuffix(path, "contracts/"+contractName+".sol")
		},
		func(path string) bool {
			return strings.HasSuffix(path, "/"+contractName+".sol") || path == contractName+".sol"
		},
		func(path string) bool {
			return strings.Contains(path, "contracts/") &&
				strings.HasSuffix(path, ".sol") &&
				!c.isLibraryPath(path)
		},
		func(path string) bool {
			return strings.HasSuffix(path, ".sol") && !c.isLibraryPath(path)
		},
	}

	for _, pattern := range patterns {
		for path, rawFile := range raw {
			if pattern(path) {
				var file sourceFile
				if err := json.Unmarshal(rawFile, &file); err == nil {
					return file.Content
				}
			}
		}
	}

	// Fallback
	for _, rawFile := range raw {
		var file sourceFile
		if err := json.Unmarshal(rawFile, &file); err == nil {
			return file.Content
		}
	}

	return ""
}

func (c *EtherscanClient) isLibraryPath(path string) bool {
	libraries := []string{
		"@openzeppelin",
		"@chainlink",
		"@uniswap",
		"node_modules",
		"@gnosis",
		"@aave",
		"@compound",
	}

	pathLower := strings.ToLower(path)
	for _, lib := range libraries {
		if strings.Contains(pathLower, lib) {
			return true
		}
	}
	return false
}
