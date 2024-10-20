package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orientallines/beesbiz/internal/config"
	"github.com/orientallines/beesbiz/internal/database"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	if err := config.Init(&config.GlobalConfig); err != nil {
		zap.S().Fatal("Failed to initialize logger: ", err)
	}

	// Load configuration
	if err := config.LoadConfig(); err != nil {
		zap.S().Fatal("Failed to load configuration: ", err)
	}

	// Set up database connection
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.GlobalConfig.PostgresUser,
		config.GlobalConfig.PostgresPassword,
		config.GlobalConfig.PostgresHost,
		config.GlobalConfig.PostgresPort,
		config.GlobalConfig.PostgresDB)

	if err := database.InitDB(dbURL); err != nil {
		zap.S().Fatal("Failed to initialize database: ", err)
	}
	defer func() {
		database.DB.Close()
	}()

	// Set up Fiber app
	app := fiber.New()

	// TODO: Add routes
	// routes.SetupRoutes(app)

	// Start server
	go func() {
		if err := app.Listen(":8080"); err != nil {
			zap.S().Fatal("Failed to start server: ", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := app.ShutdownWithContext(ctx); err != nil {
		zap.S().Fatal("Server forced to shutdown: ", err)
	}

	zap.L().Info("Server exiting")
}
