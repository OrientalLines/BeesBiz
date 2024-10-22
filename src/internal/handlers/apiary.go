package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/orientallines/beesbiz/internal/database"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

func GetApiary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid apiary ID"})
		}
		apiary, err := db.GetApiary(id)
		return c.JSON(apiary)
	}
}

func CreateApiary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var apiary types.Apiary
		if err := c.BodyParser(&apiary); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid apiary data"})
		}
		createdApiary, err := db.CreateApiary(apiary)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create apiary"})
		}
		return c.JSON(createdApiary)
	}
}

func UpdateApiary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var apiary types.Apiary
		if err := c.BodyParser(&apiary); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid apiary data"})
		}
		updatedApiary, err := db.UpdateApiary(apiary)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update apiary"})
		}
		return c.JSON(updatedApiary)
	}
}

func DeleteApiary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid apiary ID"})
		}
		if err := db.DeleteApiary(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete apiary"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllApiaries(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiaries, err := db.GetAllApiaries()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get all apiaries"})
		}
		return c.JSON(apiaries)
	}
}

// Hive Handlers
func GetAllHives(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		hives, err := db.GetAllHives()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get all hives"})
		}
		return c.JSON(hives)
	}
}

func CreateHive(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var hive types.Hive
		if err := c.BodyParser(&hive); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid hive data"})
		}
		createdHive, err := db.CreateHive(hive)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create hive"})
		}
		return c.JSON(createdHive)
	}
}

func UpdateHive(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var hive types.Hive
		if err := c.BodyParser(&hive); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid hive data"})
		}
		updatedHive, err := db.UpdateHive(hive)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update hive"})
		}
		return c.JSON(updatedHive)
	}
}

func DeleteHive(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid hive ID"})
		}
		if err := db.DeleteHive(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete hive"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllHivesByApiaryID(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiaryID, err := c.ParamsInt("apiaryID")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid apiary ID"})
		}
		hives, err := db.GetAllHivesByApiaryID(apiaryID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get hives for the apiary"})
		}
		return c.JSON(hives)
	}
}

// BeeCommunity Handlers
func GetAllBeeCommunities(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		communities, err := db.GetAllBeeCommunities()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get all bee communities"})
		}
		return c.JSON(communities)
	}
}

func CreateBeeCommunity(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var community types.BeeCommunity
		if err := c.BodyParser(&community); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid bee community data"})
		}
		createdCommunity, err := db.CreateBeeCommunity(community)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create bee community"})
		}
		return c.JSON(createdCommunity)
	}
}

func UpdateBeeCommunity(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var community types.BeeCommunity
		if err := c.BodyParser(&community); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid bee community data"})
		}
		updatedCommunity, err := db.UpdateBeeCommunity(community)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update bee community"})
		}
		return c.JSON(updatedCommunity)
	}
}

func DeleteBeeCommunity(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid bee community ID"})
		}
		if err := db.DeleteBeeCommunity(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete bee community"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllBeeCommunitiesByHiveID(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		hiveID, err := c.ParamsInt("hiveID")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid hive ID"})
		}
		communities, err := db.GetAllBeeCommunitiesByHiveID(hiveID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get bee communities for the hive"})
		}
		return c.JSON(communities)
	}
}

// HoneyHarvest Handlers
func GetHoneyHarvest(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid honey harvest ID"})
		}
		harvest, err := db.GetHoneyHarvest(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get honey harvest"})
		}
		return c.JSON(harvest)
	}
}

func CreateHoneyHarvest(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var harvest types.HoneyHarvest
		if err := c.BodyParser(&harvest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid honey harvest data"})
		}
		createdHarvest, err := db.CreateHoneyHarvest(harvest)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create honey harvest"})
		}
		return c.JSON(createdHarvest)
	}
}

func UpdateHoneyHarvest(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var harvest types.HoneyHarvest
		if err := c.BodyParser(&harvest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid honey harvest data"})
		}
		updatedHarvest, err := db.UpdateHoneyHarvest(harvest)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update honey harvest"})
		}
		return c.JSON(updatedHarvest)
	}
}

func DeleteHoneyHarvest(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid honey harvest ID"})
		}
		if err := db.DeleteHoneyHarvest(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete honey harvest"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllHoneyHarvests(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		harvests, err := db.GetAllHoneyHarvests()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get all honey harvests"})
		}
		return c.JSON(harvests)
	}
}
