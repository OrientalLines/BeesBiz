package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/orientallines/beesbiz/internal/database"
	"github.com/orientallines/beesbiz/internal/rabbitmq"
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
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid apiary data: %v", err)})
		}
		createdApiary, err := db.CreateApiary(apiary)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create apiary: %v", err)})
		}
		return c.JSON(createdApiary)
	}
}

func UpdateApiary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var apiary types.Apiary
		if err := c.BodyParser(&apiary); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid apiary data: %v", err)})
		}
		updatedApiary, err := db.UpdateApiary(apiary)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update apiary: %v", err)})
		}
		return c.JSON(updatedApiary)
	}
}

func DeleteApiary(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid apiary ID: %v", err)})
		}
		if err := db.DeleteApiary(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete apiary: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllApiaries(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiaries, err := db.GetAllApiaries()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all apiaries: %v", err)})
		}
		return c.JSON(apiaries)
	}
}

// Hive Handlers
func GetAllHives(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		hives, err := db.GetAllHives()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all hives: %v", err)})
		}
		return c.JSON(hives)
	}
}

func CreateHive(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var hive types.Hive
		if err := c.BodyParser(&hive); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid hive data: %v", err)})
		}
		createdHive, err := db.CreateHive(hive)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create hive: %v", err)})
		}
		return c.JSON(createdHive)
	}
}

func UpdateHive(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var hive types.Hive
		if err := c.BodyParser(&hive); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid hive data: %v", err)})
		}
		updatedHive, err := db.UpdateHive(hive)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update hive: %v", err)})
		}
		return c.JSON(updatedHive)
	}
}

func DeleteHive(db *database.DB, rmq *rabbitmq.RabbitMQ) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid hive ID: %v", err)})
		}

		// Get all sensors for this hive before deleting
		sensors, err := db.GetAllSensorsByHiveID(id)
		if err != nil {
			zap.L().Error("Failed to get sensors for hive",
				zap.Error(err),
				zap.Int("hive_id", id))
			// Continue with hive deletion even if we can't get sensors
		} else if len(sensors) > 0 {
			// Send delete messages for each sensor
			for _, sensor := range sensors {
				deleteMsg := types.DeleteSensor{
					HiveID:   id,
					SensorID: sensor.SensorID,
				}

				if err := rmq.PublishMessage(rabbitmq.DeleteSensorQueue, deleteMsg); err != nil {
					zap.L().Error("Failed to publish sensor deletion message",
						zap.Error(err),
						zap.Int("sensor_id", sensor.SensorID),
						zap.Int("hive_id", id))
					// Continue with other sensors even if one fails
				}
			}
		}

		// Delete the hive
		if err := db.DeleteHive(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to delete hive: %v", err),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllHivesByApiaryID(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiaryID, err := c.ParamsInt("apiaryID")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid apiary ID: %v", err)})
		}
		hives, err := db.GetAllHivesByApiaryID(apiaryID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get hives for the apiary: %v", err)})
		}
		return c.JSON(hives)
	}
}

// BeeCommunity Handlers
func GetAllBeeCommunities(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		communities, err := db.GetAllBeeCommunities()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all bee communities: %v", err)})
		}
		return c.JSON(communities)
	}
}

func CreateBeeCommunity(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var community types.BeeCommunity
		if err := c.BodyParser(&community); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid bee community data: %v", err)})
		}
		createdCommunity, err := db.CreateBeeCommunity(community)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create bee community: %v", err)})
		}
		return c.JSON(createdCommunity)
	}
}

func UpdateBeeCommunity(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var community types.BeeCommunity
		if err := c.BodyParser(&community); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid bee community data: %v", err)})
		}
		updatedCommunity, err := db.UpdateBeeCommunity(community)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update bee community: %v", err)})
		}
		return c.JSON(updatedCommunity)
	}
}

func DeleteBeeCommunity(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid bee community ID: %v", err)})
		}
		if err := db.DeleteBeeCommunity(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete bee community: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllBeeCommunitiesByHiveID(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		hiveID, err := c.ParamsInt("hiveID")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid hive ID: %v", err)})
		}
		communities, err := db.GetAllBeeCommunitiesByHiveID(hiveID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get bee communities for the hive: %v", err)})
		}
		return c.JSON(communities)
	}
}

// HoneyHarvest Handlers
func GetHoneyHarvest(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid honey harvest ID: %v", err)})
		}
		harvest, err := db.GetHoneyHarvest(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get honey harvest: %v", err)})
		}
		return c.JSON(harvest)
	}
}

func CreateHoneyHarvest(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var harvest types.HoneyHarvest
		if err := c.BodyParser(&harvest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid honey harvest data: %v", err)})
		}
		createdHarvest, err := db.CreateHoneyHarvest(harvest)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create honey harvest %v", err)})
		}
		return c.JSON(createdHarvest)
	}
}

func UpdateHoneyHarvest(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var harvest types.HoneyHarvest
		if err := c.BodyParser(&harvest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid honey harvest data: %v", err)})
		}
		updatedHarvest, err := db.UpdateHoneyHarvest(harvest)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update honey harvest: %v", err)})
		}
		return c.JSON(updatedHarvest)
	}
}

func DeleteHoneyHarvest(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid honey harvest ID: %v", err)})
		}
		if err := db.DeleteHoneyHarvest(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete honey harvest: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllHoneyHarvests(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		harvests, err := db.GetAllHoneyHarvests()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all honey harvests: %v", err)})
		}
		return c.JSON(harvests)
	}
}
