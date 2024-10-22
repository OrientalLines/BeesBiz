package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/orientallines/beesbiz/internal/database"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

// VeterinaryPassport handlers
func GetVeterinaryPassport(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid passport ID: %v", err)})
		}
		passport, err := db.GetVeterinaryPassport(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get veterinary passport: %v", err)})
		}
		return c.JSON(passport)
	}
}

func CreateVeterinaryPassport(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var passport types.VeterinaryPassport
		if err := c.BodyParser(&passport); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid passport data: %v", err)})
		}
		createdPassport, err := db.CreateVeterinaryPassport(passport)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create veterinary passport: %v", err)})
		}
		return c.JSON(createdPassport)
	}
}

func UpdateVeterinaryPassport(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var passport types.VeterinaryPassport
		if err := c.BodyParser(&passport); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid passport data: %v", err)})
		}
		updatedPassport, err := db.UpdateVeterinaryPassport(passport)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update veterinary passport: %v", err)})
		}
		return c.JSON(updatedPassport)
	}
}

func DeleteVeterinaryPassport(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid passport ID: %v", err)})
		}
		if err := db.DeleteVeterinaryPassport(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete veterinary passport: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllVeterinaryPassports(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		passports, err := db.GetAllVeterinaryPassports()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all veterinary passports: %v", err)})
		}
		return c.JSON(passports)
	}
}

// VeterinaryRecord handlers
func GetVeterinaryRecord(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid record ID: %v", err)})
		}
		record, err := db.GetVeterinaryRecord(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get veterinary record: %v", err)})
		}
		return c.JSON(record)
	}
}

func CreateVeterinaryRecord(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var record types.VeterinaryRecord
		if err := c.BodyParser(&record); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid record data: %v", err)})
		}
		createdRecord, err := db.CreateVeterinaryRecord(record)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create veterinary record: %v", err)})
		}
		return c.JSON(createdRecord)
	}
}

func UpdateVeterinaryRecord(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var record types.VeterinaryRecord
		if err := c.BodyParser(&record); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid record data: %v", err)})
		}
		updatedRecord, err := db.UpdateVeterinaryRecord(record)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update veterinary record: %v", err)})
		}
		return c.JSON(updatedRecord)
	}
}

func DeleteVeterinaryRecord(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid record ID: %v", err)})
		}
		if err := db.DeleteVeterinaryRecord(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete veterinary record: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllVeterinaryRecords(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		records, err := db.GetAllVeterinaryRecords()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all veterinary records: %v", err)})
		}
		return c.JSON(records)
	}
}
