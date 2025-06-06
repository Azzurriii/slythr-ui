package routes

import (
	config "github.com/Azzurriii/slythr-go-backend/config"
	contractHandlers "github.com/Azzurriii/slythr-go-backend/internal/application/handlers/contracts"
	"github.com/Azzurriii/slythr-go-backend/internal/application/services"
	"github.com/Azzurriii/slythr-go-backend/internal/interface/http/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Route defines the structure for dynamic routing
type Route struct {
	Method      string
	Path        string
	HandlerFunc gin.HandlerFunc
}

// Controller defines the structure for a controller with routes
type Controller struct {
	Routes []Route
}

// SetupRouter dynamically sets up routes
func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.Server.Env)

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Initialize services
	contractService := services.NewContractServiceWithDefaults(cfg, db)

	// Initialize handlers
	contractHandler := contractHandlers.NewContractHandler(contractService)

	// API v1 group
	apiV1 := r.Group("/api/v1")
	{
		// Contract routes
		contractRoutes := apiV1.Group("/contracts")
		{
			contractRoutes.GET("/:address/source-code", contractHandler.GetSourceCode)
		}
	}
	return r
}
