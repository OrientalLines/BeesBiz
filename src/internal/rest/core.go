package rest

import (
	"context"

	"github.com/gofiber/fiber/v3"
)

type Server struct {
	app *fiber.App
}

func NewServer() *Server {
	return &Server{
		app: fiber.New(),
	}
}

func (s *Server) SetupRoutes() {
	s.app.Get("/tikv/:key", func(c fiber.Ctx) error {
		// TODO: Implement getState logic
		return c.SendString("Get state placeholder")
	})
	s.app.Post("/tikv/:key", func(c fiber.Ctx) error {
		// TODO: Implement saveState logic
		return c.SendString("Save state placeholder")
	})
}

func (s *Server) Run(address string) error {
	return s.app.Listen(address)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
