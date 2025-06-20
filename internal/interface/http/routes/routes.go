package routes

import (
	config "github.com/Azzurriii/slythr/config"
	contractHandlers "github.com/Azzurriii/slythr/internal/application/handlers/contracts"
	dynamicAnalysisHandlers "github.com/Azzurriii/slythr/internal/application/handlers/dynamic_analysis"
	staticAnalysisHandlers "github.com/Azzurriii/slythr/internal/application/handlers/static_analysis"
	"github.com/Azzurriii/slythr/internal/application/services"
	"github.com/Azzurriii/slythr/internal/domain/repository"
	"github.com/Azzurriii/slythr/internal/infrastructure/external"
	"github.com/Azzurriii/slythr/internal/interface/http/middleware"
	"github.com/gin-gonic/gin"
)

type Logger interface {
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

type RouterDependencies struct {
	ContractRepo        repository.ContractRepository
	DynamicAnalysisRepo repository.DynamicAnalysisRepository
	StaticAnalysisRepo  repository.StaticAnalysisRepository
	EtherscanClient     external.EtherscanService
	Logger              Logger
	Config              *config.Config
}

func SetupRouter(deps *RouterDependencies) *gin.Engine {
	gin.SetMode(deps.Config.Server.Env)

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	contractService := services.NewContractService(
		deps.ContractRepo,
		deps.EtherscanClient,
	)

	dynamicAnalysisService, err := services.NewDynamicAnalysisService(
		deps.DynamicAnalysisRepo,
		deps.ContractRepo,
		nil, // Use default options
	)
	if err != nil {
		panic("Failed to create dynamic analysis service: " + err.Error())
	}

	staticAnalysisService, err := services.NewStaticAnalysisService(
		deps.StaticAnalysisRepo,
		nil, // Use default options
	)
	if err != nil {
		panic("Failed to create static analysis service: " + err.Error())
	}

	// Create handlers
	contractHandler := contractHandlers.NewContractHandler(contractService)
	staticAnalysisHandler := staticAnalysisHandlers.NewStaticAnalysisHandler(staticAnalysisService)
	dynamicAnalysisHandler := dynamicAnalysisHandlers.NewDynamicAnalysisHandler(dynamicAnalysisService)

	setupAPIRoutes(r, contractHandler, staticAnalysisHandler, dynamicAnalysisHandler)

	return r
}

func setupAPIRoutes(router *gin.Engine, contractHandler *contractHandlers.ContractHandler, staticAnalysisHandler *staticAnalysisHandlers.StaticAnalysisHandler, dynamicAnalysisHandler *dynamicAnalysisHandlers.DynamicAnalysisHandler) {
	apiV1 := router.Group("/api/v1")
	{
		setupContractRoutes(apiV1, contractHandler)
		setupStaticAnalysisRoutes(apiV1, staticAnalysisHandler)
		setupDynamicAnalysisRoutes(apiV1, dynamicAnalysisHandler)
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

func setupDynamicAnalysisRoutes(group *gin.RouterGroup, handler *dynamicAnalysisHandlers.DynamicAnalysisHandler) {
	group.POST("/dynamic-analysis", handler.AnalyzeContract)
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
