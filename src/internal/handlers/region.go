package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/orientallines/beesbiz/internal/database"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

// Region handlers
func GetRegion(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid region ID"})
		}
		region, err := db.GetRegion(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get region"})
		}
		return c.JSON(region)
	}
}

func CreateRegion(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var region types.Region
		if err := c.BodyParser(&region); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid region data"})
		}
		createdRegion, err := db.CreateRegion(region)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create region"})
		}
		return c.JSON(createdRegion)
	}
}

func UpdateRegion(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var region types.Region
		if err := c.BodyParser(&region); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid region data"})
		}
		updatedRegion, err := db.UpdateRegion(region)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update region"})
		}
		return c.JSON(updatedRegion)
	}
}

func DeleteRegion(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid region ID"})
		}
		if err := db.DeleteRegion(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete region"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllRegions(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		regions, err := db.GetAllRegions()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get all regions"})
		}
		return c.JSON(regions)
	}
}

// AllowedRegion handlers
func CreateAllowedRegion(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var allowedRegion types.AllowedRegion
		if err := c.BodyParser(&allowedRegion); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid allowed region data"})
		}
		createdAllowedRegion, err := db.CreateAllowedRegion(allowedRegion)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create allowed region"})
		}
		return c.JSON(createdAllowedRegion)
	}
}

func UpdateAllowedRegion(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var allowedRegion types.AllowedRegion
		if err := c.BodyParser(&allowedRegion); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid allowed region data"})
		}
		updatedAllowedRegion, err := db.UpdateAllowedRegion(allowedRegion)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update allowed region"})
		}
		return c.JSON(updatedAllowedRegion)
	}
}

func DeleteAllowedRegion(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid allowed region ID"})
		}
		if err := db.DeleteAllowedRegion(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete allowed region"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllAllowedRegions(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		allowedRegions, err := db.GetAllAllowedRegions()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get all allowed regions"})
		}
		return c.JSON(allowedRegions)
	}
}

// RegionApiary handlers
func GetRegionApiary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid region apiary ID"})
		}
		regionApiary, err := db.GetRegionApiary(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get region apiary"})
		}
		return c.JSON(regionApiary)
	}
}

func CreateRegionApiary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var regionApiary types.RegionApiary
		if err := c.BodyParser(&regionApiary); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid region apiary data"})
		}
		createdRegionApiary, err := db.CreateRegionApiary(regionApiary)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create region apiary"})
		}
		return c.JSON(createdRegionApiary)
	}
}

func UpdateRegionApiary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var regionApiary types.RegionApiary
		if err := c.BodyParser(&regionApiary); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid region apiary data"})
		}
		updatedRegionApiary, err := db.UpdateRegionApiary(regionApiary)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update region apiary"})
		}
		return c.JSON(updatedRegionApiary)
	}
}

func DeleteRegionApiary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid region apiary ID"})
		}
		if err := db.DeleteRegionApiary(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete region apiary"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllRegionApiaries(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		regionApiaries, err := db.GetAllRegionApiaries()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get all region apiaries"})
		}
		return c.JSON(regionApiaries)
	}
}
