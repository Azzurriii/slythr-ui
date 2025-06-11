package routes

import (
	config "github.com/Azzurriii/slythr-go-backend/config"
	contractHandlers "github.com/Azzurriii/slythr-go-backend/internal/application/handlers/contracts"
	staticAnalysisHandlers "github.com/Azzurriii/slythr-go-backend/internal/application/handlers/static_analysis"
	"github.com/Azzurriii/slythr-go-backend/internal/application/services"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/repository"
	"github.com/Azzurriii/slythr-go-backend/internal/infrastructure/external"
	"github.com/Azzurriii/slythr-go-backend/internal/interface/http/middleware"
	"github.com/gin-gonic/gin"
)

type Logger interface {
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

type RouterDependencies struct {
	ContractRepo    repository.ContractRepository
	EtherscanClient external.EtherscanService
	Logger          Logger
	Config          *config.Config
}

func SetupRouter(deps *RouterDependencies) *gin.Engine {
	gin.SetMode(deps.Config.Server.Env)

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	contractService := services.NewContractService(
		deps.ContractRepo,
		deps.EtherscanClient,
		deps.Logger,
	)

	contractHandler := contractHandlers.NewContractHandler(contractService)
	staticAnalysisHandler := staticAnalysisHandlers.NewStaticAnalysisHandler()

	setupAPIRoutes(r, contractHandler, staticAnalysisHandler)

	return r
}

func setupAPIRoutes(router *gin.Engine, contractHandler *contractHandlers.ContractHandler, staticAnalysisHandler *staticAnalysisHandlers.StaticAnalysisHandler) {
	apiV1 := router.Group("/api/v1")
	{
		setupContractRoutes(apiV1, contractHandler)
		setupStaticAnalysisRoutes(apiV1, staticAnalysisHandler)
	}
}

func setupContractRoutes(group *gin.RouterGroup, handler *contractHandlers.ContractHandler) {
	contracts := group.Group("/contracts")
	{
		contracts.GET("/:address", handler.GetContract)
		contracts.GET("/:address/source-code", handler.GetSourceCode)
	}
}

func setupStaticAnalysisRoutes(group *gin.RouterGroup, handler *staticAnalysisHandlers.StaticAnalysisHandler) {
	group.POST("/static-analysis", handler.AnalyzeContract)
}

func SetupRouterLegacy(repo repository.ContractRepository, cfg *config.Config, etherscanClient external.EtherscanService, logger Logger) *gin.Engine {
	deps := &RouterDependencies{
		ContractRepo:    repo,
		EtherscanClient: etherscanClient,
		Logger:          logger,
		Config:          cfg,
	}
	return SetupRouter(deps)
}
