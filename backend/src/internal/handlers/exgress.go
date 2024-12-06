package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/orientallines/beesbiz/internal/database"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

// Production Report Handlers

// GetProductionReport gets a production report by ID
func GetProductionReport(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid production report ID: %v", err)})
		}
		report, err := db.GetProductionReport(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get production report: %v", err)})
		}
		return c.JSON(report)
	}
}

// CreateProductionReport creates a new production report
func CreateProductionReport(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var report types.ProductionReport
		if err := c.BodyParser(&report); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid production report data: %v", err)})
		}
		createdReport, err := db.CreateProductionReport(report)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create production report: %v", err)})
		}
		return c.JSON(createdReport)
	}
}

// UpdateProductionReport updates a production report
func UpdateProductionReport(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var report types.ProductionReport
		if err := c.BodyParser(&report); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid production report data: %v", err)})
		}
		updatedReport, err := db.UpdateProductionReport(report)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update production report: %v", err)})
		}
		return c.JSON(updatedReport)
	}
}

// DeleteProductionReport deletes a production report
func DeleteProductionReport(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid production report ID: %v", err)})
		}
		if err := db.DeleteProductionReport(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete production report: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// GetAllProductionReports gets all production reports
func GetAllProductionReports(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		reports, err := db.GetAllProductionReports()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all production reports: %v", err)})
		}
		return c.JSON(reports)
	}
}

func GetCuratedProductionReportsByUser(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid user ID: %v", err)})
		}
		reports, err := db.GetCuratedProductionReportsByUser(userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get production reports by user: %v", err)})
		}
		return c.JSON(reports)
	}
}
