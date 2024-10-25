package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/orientallines/beesbiz/internal/database"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

// Region handlers
func GetRegion(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid region ID: %v", err)})
		}
		region, err := db.GetRegion(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get region: %v", err)})
		}
		return c.JSON(region)
	}
}

func CreateRegion(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var region types.Region
		if err := c.BodyParser(&region); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid region data: %v", err)})
		}
		createdRegion, err := db.CreateRegion(region)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create region: %v", err)})
		}
		return c.JSON(createdRegion)
	}
}

func UpdateRegion(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var region types.Region
		if err := c.BodyParser(&region); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid region data: %v", err)})
		}
		updatedRegion, err := db.UpdateRegion(region)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update region: %v", err)})
		}
		return c.JSON(updatedRegion)
	}
}

func DeleteRegion(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid region ID: %v", err)})
		}
		if err := db.DeleteRegion(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete region: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllRegions(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		regions, err := db.GetAllRegions()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all regions: %v", err)})
		}
		return c.JSON(regions)
	}
}

// AllowedRegion handlers
func CreateAllowedRegion(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var allowedRegion types.AllowedRegion
		if err := c.BodyParser(&allowedRegion); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid allowed region data: %v", err)})
		}

		createdAllowedRegion, err := db.CreateAllowedRegion(allowedRegion)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create allowed region: %v", err)})
		}
		return c.JSON(createdAllowedRegion)
	}
}

func UpdateAllowedRegion(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var allowedRegion types.AllowedRegion
		if err := c.BodyParser(&allowedRegion); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid allowed region data: %v", err)})
		}

		updatedAllowedRegion, err := db.UpdateAllowedRegion(allowedRegion)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update allowed region: %v", err)})
		}
		return c.JSON(updatedAllowedRegion)
	}
}

func DeleteAllowedRegion(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid allowed region ID: %v", err)})
		}
		if err := db.DeleteAllowedRegion(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete allowed region: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllAllowedRegions(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		allowedRegions, err := db.GetAllAllowedRegions()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all allowed regions: %v", err)})
		}
		return c.JSON(allowedRegions)
	}
}

// GetAllowedRegionsForUser handler
func GetAllowedRegionsForUser(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid allowed region ID: %v", err)})
		}
		allowedRegions, err := db.GetAllowedRegions(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get allowed region: %v", err)})
		}
		return c.JSON(allowedRegions)
	}
}

// RegionApiary handlers
func GetRegionApiary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid region apiary ID: %v", err)})
		}
		regionApiary, err := db.GetRegionApiary(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get region apiary: %v", err)})
		}
		return c.JSON(regionApiary)
	}
}

func CreateRegionApiary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var regionApiary types.RegionApiary
		if err := c.BodyParser(&regionApiary); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid region apiary data: %v", err)})
		}

		createdRegionApiary, err := db.CreateRegionApiary(regionApiary)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create region apiary: %v", err)})
		}
		return c.JSON(createdRegionApiary)
	}
}

func UpdateRegionApiary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var regionApiary types.RegionApiary
		if err := c.BodyParser(&regionApiary); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid region apiary data: %v", err)})
		}

		updatedRegionApiary, err := db.UpdateRegionApiary(regionApiary)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update region apiary: %v", err)})
		}
		return c.JSON(updatedRegionApiary)
	}
}

func DeleteRegionApiary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid region apiary ID: %v", err)})
		}
		if err := db.DeleteRegionApiary(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete region apiary: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllRegionApiaries(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		regionApiaries, err := db.GetAllRegionApiaries()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all region apiaries: %v", err)})
		}
		return c.JSON(regionApiaries)
	}
}
