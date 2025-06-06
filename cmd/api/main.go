package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/Azzurriii/slythr-go-backend/config"
	_ "github.com/Azzurriii/slythr-go-backend/docs"
	database "github.com/Azzurriii/slythr-go-backend/internal/infrastructure/database"
	routes "github.com/Azzurriii/slythr-go-backend/internal/interface/http/routes"
	server "github.com/Azzurriii/slythr-go-backend/internal/interface/server"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Slyther Go Backend API
// @version 1.0
// @description This is a sample server for Slyther Go Backend.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// init dbs
	_ = database.InitDatabases(database.NewPostgresConfig(), database.RedisConfig(cfg.Redis))

	// Initialize PostgreSQL
	db := database.GetPostgres()
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get DB connection: %v", err)
	}
	defer sqlDb.Close()

	// Initialize Redis
	redisClient := database.GetRedis()
	defer redisClient.Close()

	// Setup router
	router := routes.SetupRouter(db, cfg)

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Use the server abstraction
	srv := server.NewServer(router)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		fmt.Println("Shutting down server...")

		// Create shutdown context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Shutdown services gracefully
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}

		redisClient.Close()
		sqlDb.Close()
		fmt.Println("Server gracefully stopped")
	}()

	// Start server
	port := cfg.Server.Port
	if err := srv.Start(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
