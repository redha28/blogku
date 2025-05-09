package main

import (
	"log"
	"os"
	"path/filepath"

	_ "github.com/joho/godotenv/autoload"
	// "github.com/redha28/blogku/internal/handlers"

	"github.com/redha28/blogku/internals/middlewares"
	"github.com/redha28/blogku/internals/routes"
	"github.com/redha28/blogku/pkg"

	// "github.com/redha28/blogku/pkg/handlers"

	_ "github.com/redha28/blogku/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Blog CMS API
// @version 1.0
// @description This is a Blog CMS API server.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
func main() {
	// Initialize logger
	logPath := filepath.Join("logs", "app.log")
	logger, err := pkg.InitLogger(pkg.LevelDebug, logPath)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	pkg.Info("Starting Blog CMS application...")

	// Connect to MySQL
	mySql, err := pkg.Connect()
	if err != nil {
		pkg.Error("Unable to create database connection pool", err)
		os.Exit(1)
	}
	defer func() {
		pkg.Info("Closing DB connection")
		mySql.Close()
	}()

	// Connect to Redis
	pkg.Info("Connecting to Redis...")
	rdb := pkg.RedisConnect()

	// Initialize router
	pkg.Info("Initializing router...")
	router := routes.InitRouter(mySql, rdb)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Middleware
	router.Use(middlewares.LoggerMiddleware())
	router.Use(middlewares.RecoveryMiddleware())

	// Add CORS middleware
	// router.Use(handlers.)

	// Start server
	serverAddr := "localhost:8080"
	pkg.Info("Server starting on " + serverAddr)
	if err := router.Run(serverAddr); err != nil {
		pkg.Error("Failed to start server", err)
		os.Exit(1)
	}
}
