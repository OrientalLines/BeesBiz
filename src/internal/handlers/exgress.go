package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/orientallines/beesbiz/internal/database"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

// Production Report Handlers

func GetProductionReport(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid production report ID"})
		}
		report, err := db.GetProductionReport(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get production report"})
		}
		return c.JSON(report)
	}
}

func CreateProductionReport(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var report types.ProductionReport
		if err := c.BodyParser(&report); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid production report data"})
		}
		createdReport, err := db.CreateProductionReport(report)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create production report"})
		}
		return c.JSON(createdReport)
	}
}

func UpdateProductionReport(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var report types.ProductionReport
		if err := c.BodyParser(&report); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid production report data"})
		}
		updatedReport, err := db.UpdateProductionReport(report)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update production report"})
		}
		return c.JSON(updatedReport)
	}
}

func DeleteProductionReport(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid production report ID"})
		}
		if err := db.DeleteProductionReport(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete production report"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllProductionReports(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		reports, err := db.GetAllProductionReports()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get all production reports"})
		}
		return c.JSON(reports)
	}
}
