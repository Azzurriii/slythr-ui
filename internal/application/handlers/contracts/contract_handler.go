package contracts

import (
	"errors"
	"net/http"

	contractsDTO "github.com/Azzurriii/slythr-go-backend/internal/application/dto/contracts"
	"github.com/Azzurriii/slythr-go-backend/internal/application/services"
	apperrors "github.com/Azzurriii/slythr-go-backend/internal/domain/errors"
	"github.com/gin-gonic/gin"
)

type ContractHandler struct {
	contractService *services.ContractService
}

func NewContractHandler(contractService *services.ContractService) *ContractHandler {
	return &ContractHandler{
		contractService: contractService,
	}
}

// GetSourceCode godoc
// @Summary Get contract source code
// @Description Get the source code of a smart contract from Etherscan by its address and save it to the database
// @Tags contracts
// @Accept json
// @Produce json
// @Param address path string true "Contract Address"
// @Param network query string false "Network Name" default(ethereum)
// @Router /contracts/{address}/source-code [get]
func (h *ContractHandler) GetSourceCode(c *gin.Context) {
	address := c.Param("address")
	network := c.DefaultQuery("network", services.NetworkEthereum) // Default to Ethereum

	if address == "" {
		c.JSON(http.StatusBadRequest, apperrors.NewErrorResponse(apperrors.ErrInvalidAddress))
		return
	}

	req := &contractsDTO.GetContractSourceCodeRequest{
		Address: address,
		Network: network,
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, apperrors.NewErrorResponse(err))
		return
	}

	response, err := h.contractService.GetAndSaveContractSourceCode(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrContractNotFound):
			c.JSON(http.StatusNotFound, apperrors.NewErrorResponse(err))
		case errors.Is(err, apperrors.ErrInvalidAddress):
			c.JSON(http.StatusBadRequest, apperrors.NewErrorResponse(err))
		default:
			c.JSON(http.StatusInternalServerError, apperrors.NewErrorResponse(err))
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshSourceCode godoc
// @Summary Refresh contract source code
// @Description Forces a refresh of the source code of a smart contract from Etherscan by its address
// @Tags contracts
// @Accept json
// @Produce json
// @Param address path string true "Contract Address"
// @Param network query string false "Network Name" default(ethereum)
// @Router /contracts/{address}/source-code [post]
func (h *ContractHandler) RefreshSourceCode(c *gin.Context) {
	address := c.Param("address")
	network := c.DefaultQuery("network", services.NetworkEthereum) // Default to Ethereum

	req := &contractsDTO.RefreshContractSourceCodeRequest{
		Address: address,
		Network: network,
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, apperrors.NewErrorResponse(err))
		return
	}

	response, err := h.contractService.RefreshContractSourceCode(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrContractNotFound):
			c.JSON(http.StatusNotFound, apperrors.NewErrorResponse(err))
		case errors.Is(err, apperrors.ErrInvalidAddress):
			c.JSON(http.StatusBadRequest, apperrors.NewErrorResponse(err))
		default:
			c.JSON(http.StatusInternalServerError, apperrors.NewErrorResponse(err))
		}
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
// @Param address path string true "Contract Address"
// @Param network query string false "Network Name" default(ethereum)
// @Router /contracts/{address} [get]
func (h *ContractHandler) GetContract(c *gin.Context) {
	address := c.Param("address")
	network := c.DefaultQuery("network", services.NetworkEthereum)

	if address == "" {
		c.JSON(http.StatusBadRequest, apperrors.NewErrorResponse(apperrors.ErrInvalidAddress))
		return
	}

	response, err := h.contractService.GetContractByAddressAndNetwork(c.Request.Context(), address, network)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrContractNotFound):
			c.JSON(http.StatusNotFound, apperrors.NewErrorResponse(err))
		default:
			c.JSON(http.StatusInternalServerError, apperrors.NewErrorResponse(err))
		}
		return
	}

	c.JSON(http.StatusOK, response)
}
