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

type multiFileSource struct {
	Sources map[string]fileContent `json:"sources"`
}

type fileContent struct {
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

	return c.processSourceCode(contractInfo.SourceCode)
}

func (c *EtherscanClient) GetContractDetails(ctx context.Context, address string, network string) (*ContractInfo, error) {
	contractInfo, err := c.fetchContractInfo(ctx, address, network)
	if err != nil {
		return nil, err
	}

	processedSourceCode, err := c.processSourceCode(contractInfo.SourceCode)
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

func (c *EtherscanClient) processSourceCode(sourceCode string) (string, error) {
	if len(sourceCode) == 0 {
		return "", fmt.Errorf("empty source code")
	}

	if sourceCode[0] != '{' {
		return sourceCode, nil
	}

	if processedSource := c.parseMultiFileSource(sourceCode); processedSource != "" {
		return processedSource, nil
	}

	return sourceCode, nil
}

func (c *EtherscanClient) parseMultiFileSource(sourceCode string) string {
	var multiFile multiFileSource
	if err := json.Unmarshal([]byte(sourceCode), &multiFile); err == nil && len(multiFile.Sources) > 0 {
		return c.combineSourceFiles(multiFile.Sources)
	}

	// Try alternative format (direct map)
	var altFormat map[string]fileContent
	if err := json.Unmarshal([]byte(sourceCode), &altFormat); err == nil && len(altFormat) > 0 {
		return c.combineSourceFiles(altFormat)
	}

	return ""
}

func (c *EtherscanClient) combineSourceFiles(sources map[string]fileContent) string {
	if len(sources) == 0 {
		return ""
	}

	estimatedSize := 0
	for path, content := range sources {
		estimatedSize += len(path) + len(content.Content) + 20
	}

	builder := stringBuilderPool.Get().(*strings.Builder)
	defer func() {
		builder.Reset()
		stringBuilderPool.Put(builder)
	}()

	builder.Grow(estimatedSize)

	for path, fileContent := range sources {
		builder.WriteString("// File: ")
		builder.WriteString(path)
		builder.WriteString("\n")
		builder.WriteString(fileContent.Content)
		builder.WriteString("\n\n")
	}

	return builder.String()
}
