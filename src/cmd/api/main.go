package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	dapr "github.com/dapr/go-sdk/client"
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
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
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
	app.Get("/tikv/:key", getState("tikv-store"))
	app.Post("/tikv/:key", saveState("tikv-store"))
	app.Post("/publish", publishEvent)

	// Start server
	go func() {
		if err := app.Listen(":4040"); err != nil {
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

// func main() {
// 	// Initialize logger
// 	if err := config.Init(&config.GlobalConfig); err != nil {
// 		zap.S().Fatal("Failed to initialize logger: ", err)
// 	}
// 	defer zap.S().Sync()

// 	app := fiber.New()

// 	app.Get("/psql/:key", getState("postgresql-store"))
// 	app.Post("/psql/:key", saveState("postgresql-store"))
// 	app.Get("/tikv/:key", getState("tikv-store"))
// 	app.Post("/tikv/:key", saveState("tikv-store"))
// 	app.Post("/publish", publishEvent)

// 	port := "4040"
// 	log.Printf("Starting server on :%s", port)
// 	log.Fatal(app.Listen(":" + port))
// }

func getState(storeName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		client, err := dapr.NewClient()
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Failed to create Dapr client: %v", err))
		}
		defer client.Close()

		key := c.Params("key")
		state, err := client.GetState(c.Context(), storeName, key, nil)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		if state.Value == nil {
			return c.Status(404).SendString("Key not found")
		}
		return c.Send(state.Value)
	}
}

func saveState(storeName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		client, err := dapr.NewClient()
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Failed to create Dapr client: %v", err))
		}
		defer client.Close()

		key := c.Params("key")
		value := c.Body()
		err = client.SaveState(c.Context(), storeName, key, value, nil)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.SendString("Data saved successfully")
	}
}

func publishEvent(c *fiber.Ctx) error {
	client, err := dapr.NewClient()
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Failed to create Dapr client: %v", err))
	}
	defer client.Close()

	err = client.PublishEvent(c.Context(), "rabbitmq-pubsub", "topic-name", c.Body())
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.SendString("Message published")
}
