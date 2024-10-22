package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/orientallines/beesbiz/internal/database"
)

func GetApiary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	}
}

func CreateApiary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	}
}
