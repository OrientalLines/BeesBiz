package rest

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Server is a wrapper around fiber.App
type Server struct {
	app *fiber.App
	db  *sqlx.DB
}

// NewServer creates a new Server
func NewServer(db *sqlx.DB) *Server {
	return &Server{
		app: fiber.New(),
		db:  db,
	}
}

// SetupRoutes sets up the routes for the server
func (s *Server) SetupRoutes() {
	s.app.Get("/tikv/:key", func(c *fiber.Ctx) error {
		// TODO: Implement getState logic
		return c.SendString("Get state placeholder")
	})
	s.app.Post("/tikv/:key", func(c *fiber.Ctx) error {
		// TODO: Implement saveState logic
		return c.SendString("Save state placeholder")
	})
}

// Run starts the server
func (s *Server) Run(address string) error {
	return s.app.Listen(address)
}

// Shutdown shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
