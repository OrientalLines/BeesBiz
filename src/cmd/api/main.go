package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/orientallines/beesbiz/internal/config"
	"github.com/orientallines/beesbiz/internal/database"
	"github.com/orientallines/beesbiz/internal/server"
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

	// Create the database URL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.GlobalConfig.PostgresUser,
		config.GlobalConfig.PostgresPassword,
		config.GlobalConfig.PostgresHost,
		config.GlobalConfig.PostgresPort,
		config.GlobalConfig.PostgresDB)

	// Connect to the database
	db, err := database.New(dbURL)
	if err != nil {
		zap.S().Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	// Initialize the database schema
	if err := db.InitSchema("scripts", []string{"definition.sql", "ddl.sql", "indexes.sql", "functions.sql", "trigger.sql"}); err != nil {
		zap.S().Fatal("Failed to initialize database schema: ", err)
	}

	// Create the server
	srv := server.NewServer(db)

	// Start servers
	go func() {
		if err := srv.SetupAndRun(":50051", ":4040"); err != nil {
			zap.S().Fatal("Failed to start servers: ", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.S().Fatal("Server forced to shutdown: ", err)
	}

	zap.L().Info("Server exiting")
}
