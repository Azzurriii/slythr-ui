package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/Azzurriii/slythr-go-backend/config"
	_ "github.com/Azzurriii/slythr-go-backend/docs"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/entities"
	database "github.com/Azzurriii/slythr-go-backend/internal/infrastructure/database"
	"github.com/Azzurriii/slythr-go-backend/internal/infrastructure/external"
	gormRepo "github.com/Azzurriii/slythr-go-backend/internal/infrastructure/persistence/gorm"
	routes "github.com/Azzurriii/slythr-go-backend/internal/interface/http/routes"
	server "github.com/Azzurriii/slythr-go-backend/internal/interface/server"
	logger "github.com/Azzurriii/slythr-go-backend/pkg/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Slyther Go Backend API
// @version 1.0
// @description This is a server for Slyther Go Backend.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Init needed configs
	dbConfig := database.NewDatabaseConfig(cfg)
	redisConfig := database.NewRedisConfig(cfg)

	connectionManager, err := database.NewConnectionManager(dbConfig, redisConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database connections: %v", err)
	}
	defer connectionManager.Close()

	db := connectionManager.GetPostgres()

	if err := db.AutoMigrate(
		&entities.Contract{},
		&entities.StaticAnalysis{},
		&entities.DynamicAnalysis{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	contractRepo := gormRepo.NewContractRepository(db)
	dynamicAnalysisRepo := gormRepo.NewDynamicAnalysisRepository(db)
	staticAnalysisRepo := gormRepo.NewStaticAnalysisRepository(db)

	etherscanClient := external.NewEtherscanClient(&cfg.Etherscan)

	routerDependencies := &routes.RouterDependencies{
		ContractRepo:        contractRepo,
		DynamicAnalysisRepo: dynamicAnalysisRepo,
		StaticAnalysisRepo:  staticAnalysisRepo,
		EtherscanClient:     etherscanClient,
		Logger:              logger.Default,
		Config:              cfg,
	}

	router := routes.SetupRouter(routerDependencies)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := server.NewServer(router, logger.Default)

	// Shutdown channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		logger.Info("Shutting down server...")

		// Shutdown context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Shutdown server gracefully
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}

		// Close database connections gracefully
		if err := connectionManager.Close(); err != nil {
			log.Printf("Failed to close database connections: %v", err)
		}

		logger.Info("Server gracefully stopped")
	}()

	port := cfg.Server.Port
	if err := srv.Start(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
