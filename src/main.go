package main

import (
	"app/src/config"
	"app/src/database"
	"app/src/middleware"
	"app/src/router"
	"app/src/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"gorm.io/gorm"
)

// @title Nutribox API documentation
// @version 1.0.0
// @license.name MIT
// @license.url https://github.com/indrayyana/go-fiber-boilerplate/blob/main/LICENSE
// @host localhost:5000
// @BasePath /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Example Value: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Log environment variables for debugging
	fmt.Printf("Environment variables:\n")
	fmt.Printf("APP_ENV: %s\n", os.Getenv("APP_ENV"))
	fmt.Printf("APP_HOST: %s\n", os.Getenv("APP_HOST"))
	fmt.Printf("APP_PORT: %s\n", os.Getenv("APP_PORT"))
	fmt.Printf("DB_HOST: %s\n", os.Getenv("DB_HOST"))
	fmt.Printf("DB_PORT: %s\n", os.Getenv("DB_PORT"))

	// Debug config values
	fmt.Printf("\nConfig values from viper:\n")
	fmt.Printf("IsProd: %v\n", config.IsProd)
	fmt.Printf("AppHost: %s\n", config.AppHost)
	fmt.Printf("AppPort: %d\n", config.AppPort)
	fmt.Printf("DBHost: %s\n", config.DBHost)
	fmt.Printf("DBPort: %d\n", config.DBPort)

	// Give PostgreSQL a moment to fully initialize in Docker environment
	if os.Getenv("APP_ENV") == "prod" {
		utils.Log.Info("Running in production mode, waiting for database to initialize...")
		time.Sleep(5 * time.Second)
	}

	app := setupFiberApp()

	// Setup database connection
	utils.Log.Info("Setting up database connection...")
	db := setupDatabase()
	defer closeDatabase(db)

	// Setup routes
	utils.Log.Info("Setting up API routes...")
	setupRoutes(app, db)

	port := os.Getenv("PORT")
	if port == "" {
		// Jika tidak ada (saat development lokal), gunakan dari file config Anda.
		port = fmt.Sprintf("%d", config.AppPort)
	}

	// PENTING: Gunakan "0.0.0.0" untuk host agar bisa diakses di dalam container.
	address := fmt.Sprintf("0.0.0.0:%s", port)

	utils.Log.Infof("Starting server on %s", address)
	utils.Log.Infof("Starting server on %s", address)

	// Start server and handle graceful shutdown
	serverErrors := make(chan error, 1)
	go startServer(app, address, serverErrors)
	handleGracefulShutdown(ctx, app, serverErrors)
}

func setupFiberApp() *fiber.App {
	utils.Log.Info("Setting up Fiber app...")
	app := fiber.New(config.FiberConfig())

	// Middleware setup
	app.Use("/v1/auth", middleware.LimiterConfig())
	app.Use(middleware.LoggerConfig())
	app.Use(middleware.APILoggerConfig())
	app.Use(middleware.RequestBodyLoggerConfig())
	app.Use(helmet.New())
	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(middleware.RecoverConfig())

	app.Static("/uploads", "./uploads")

	return app
}

func setupDatabase() *gorm.DB {
	db := database.Connect(config.DBHost, config.DBName)
	return db
}

func setupRoutes(app *fiber.App, db *gorm.DB) {
	router.Routes(app, db)
	app.Use(utils.NotFoundHandler)
}

func startServer(app *fiber.App, address string, errs chan<- error) {
	log.Printf("Starting server on %s", address)

	// Listen for incoming connections
	if err := app.Listen(address); err != nil {
		utils.Log.Errorf("Failed to start server on %s: %v", address, err)
		fmt.Printf("Error starting server: %v\n", err)

		// Detect specific error for address already in use
		if errors.Is(err, syscall.EADDRINUSE) {
			utils.Log.Error("Port is already in use, try another port or wait a moment")
		}

		errs <- fmt.Errorf("failed to start server on %s (config: host=%s port=%d): %w",
			address, config.AppHost, config.AppPort, err)
	}
}

func closeDatabase(db *gorm.DB) {
	sqlDB, errDB := db.DB()
	if errDB != nil {
		utils.Log.Errorf("Error getting database instance: %v", errDB)
		return
	}

	if err := sqlDB.Close(); err != nil {
		utils.Log.Errorf("Error closing database connection: %v", err)
	} else {
		utils.Log.Info("Database connection closed successfully")
	}
}

func handleGracefulShutdown(ctx context.Context, app *fiber.App, serverErrors <-chan error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		utils.Log.Fatalf("Server error: %v", err)
	case <-quit:
		utils.Log.Info("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			utils.Log.Fatalf("Error during server shutdown: %v", err)
		}
	case <-ctx.Done():
		utils.Log.Info("Server exiting due to context cancellation")
	}

	utils.Log.Info("Server exited")
}
