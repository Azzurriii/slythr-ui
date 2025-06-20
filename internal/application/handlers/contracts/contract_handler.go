package contracts

import (
	"net/http"

	contractsDTO "github.com/Azzurriii/slythr/internal/application/dto/contracts"
	"github.com/Azzurriii/slythr/internal/application/services"
	"github.com/Azzurriii/slythr/internal/domain/constants"
	domainerrors "github.com/Azzurriii/slythr/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

type ContractHandler struct {
	contractService services.ContractServiceInterface
}

func NewContractHandler(contractService services.ContractServiceInterface) *ContractHandler {
	return &ContractHandler{
		contractService: contractService,
	}
}

// GetSourceCode godoc
// @Summary Fetch contract source code
// @Description Get the source code of a smart contract from Etherscan by its address and save it to the database
// @Tags contracts
// @Accept json
// @Produce json
// @Param address path string true "Contract Address" minlength(42) maxlength(42)
// @Param network query string false "Network Name" default(ethereum) Enums(ethereum,polygon,bsc,base,arbitrum,avalanche,optimism,gnosis,fantom,celo)
// @Router /contracts/{address}/source-code [get]
func (h *ContractHandler) GetSourceCode(c *gin.Context) {
	req, err := h.buildGetSourceCodeRequest(c)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	response, err := h.contractService.FetchContractSourceCode(c.Request.Context(), req)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetContract godoc
// @Summary Get contract from database
// @Description Get a cached smart contract from the database by its address and network
// @Tags contracts
// @Accept json
// @Produce json
// @Param address path string true "Contract Address" minlength(42) maxlength(42)
// @Param network query string false "Network Name" default(ethereum) Enums(ethereum,polygon,bsc,base,arbitrum,avalanche,optimism,gnosis,fantom,celo)
// @Router /contracts/{address} [get]
func (h *ContractHandler) GetContract(c *gin.Context) {
	address := c.Param("address")
	network := c.DefaultQuery("network", constants.NetworkEthereum)

	if err := h.validateAddressAndNetwork(address, network); err != nil {
		h.respondWithError(c, err)
		return
	}

	response, err := h.contractService.GetContractByAddressAndNetwork(c.Request.Context(), address, network)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Helpers

func (h *ContractHandler) buildGetSourceCodeRequest(c *gin.Context) (*contractsDTO.GetContractSourceCodeRequest, error) {
	address := c.Param("address")
	network := c.DefaultQuery("network", constants.NetworkEthereum)

	req := &contractsDTO.GetContractSourceCodeRequest{
		Address: address,
		Network: network,
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	return req, nil
}

func (h *ContractHandler) validateAddressAndNetwork(address, network string) error {
	if address == "" || len(address) != 42 {
		return domainerrors.ErrInvalidAddress
	}

	if !constants.IsValidNetwork(network) {
		return domainerrors.ErrInvalidNetwork
	}

	return nil
}

// Handles error responses in a consistent way
func (h *ContractHandler) respondWithError(c *gin.Context, err error) {
	statusCode := domainerrors.GetHTTPStatusCode(err)
	errorResponse := domainerrors.NewErrorResponse(err)
	c.JSON(statusCode, errorResponse)
}
