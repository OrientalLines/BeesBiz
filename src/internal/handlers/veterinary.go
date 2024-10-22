package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/orientallines/beesbiz/internal/database"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

// VeterinaryPassport handlers
func GetVeterinaryPassport(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid passport ID"})
		}
		passport, err := db.GetVeterinaryPassport(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get veterinary passport"})
		}
		return c.JSON(passport)
	}
}

func CreateVeterinaryPassport(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var passport types.VeterinaryPassport
		if err := c.BodyParser(&passport); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid passport data"})
		}
		createdPassport, err := db.CreateVeterinaryPassport(passport)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create veterinary passport"})
		}
		return c.JSON(createdPassport)
	}
}

func UpdateVeterinaryPassport(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var passport types.VeterinaryPassport
		if err := c.BodyParser(&passport); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid passport data"})
		}
		updatedPassport, err := db.UpdateVeterinaryPassport(passport)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update veterinary passport"})
		}
		return c.JSON(updatedPassport)
	}
}

func DeleteVeterinaryPassport(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid passport ID"})
		}
		if err := db.DeleteVeterinaryPassport(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete veterinary passport"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllVeterinaryPassports(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		passports, err := db.GetAllVeterinaryPassports()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get all veterinary passports"})
		}
		return c.JSON(passports)
	}
}

// VeterinaryRecord handlers
func GetVeterinaryRecord(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid record ID"})
		}
		record, err := db.GetVeterinaryRecord(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get veterinary record"})
		}
		return c.JSON(record)
	}
}

func CreateVeterinaryRecord(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var record types.VeterinaryRecord
		if err := c.BodyParser(&record); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid record data"})
		}
		createdRecord, err := db.CreateVeterinaryRecord(record)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create veterinary record"})
		}
		return c.JSON(createdRecord)
	}
}

func UpdateVeterinaryRecord(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var record types.VeterinaryRecord
		if err := c.BodyParser(&record); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid record data"})
		}
		updatedRecord, err := db.UpdateVeterinaryRecord(record)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update veterinary record"})
		}
		return c.JSON(updatedRecord)
	}
}

func DeleteVeterinaryRecord(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid record ID"})
		}
		if err := db.DeleteVeterinaryRecord(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete veterinary record"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllVeterinaryRecords(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		records, err := db.GetAllVeterinaryRecords()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get all veterinary records"})
		}
		return c.JSON(records)
	}
}
