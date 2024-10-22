package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/orientallines/beesbiz/internal/database"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

// Sensor handlers
func GetSensor(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sensor ID"})
		}
		sensor, err := db.GetSensor(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get sensor"})
		}
		return c.JSON(sensor)
	}
}

func CreateSensor(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var sensor types.Sensor
		if err := c.BodyParser(&sensor); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sensor data"})
		}
		createdSensor, err := db.CreateSensor(sensor)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create sensor"})
		}
		return c.JSON(createdSensor)
	}
}

func UpdateSensor(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var sensor types.Sensor
		if err := c.BodyParser(&sensor); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sensor data"})
		}
		updatedSensor, err := db.UpdateSensor(sensor)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update sensor"})
		}
		return c.JSON(updatedSensor)
	}
}

func DeleteSensor(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sensor ID"})
		}
		if err := db.DeleteSensor(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete sensor"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllSensors(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sensors, err := db.GetAllSensors()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get all sensors"})
		}
		return c.JSON(sensors)
	}
}

// SensorReading handlers
func GetSensorReading(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sensor reading ID"})
		}
		reading, err := db.GetSensorReading(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get sensor reading"})
		}
		return c.JSON(reading)
	}
}

func CreateSensorReading(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reading types.SensorReading
		if err := c.BodyParser(&reading); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sensor reading data"})
		}
		createdReading, err := db.CreateSensorReading(reading)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create sensor reading"})
		}
		return c.JSON(createdReading)
	}
}

func UpdateSensorReading(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reading types.SensorReading
		if err := c.BodyParser(&reading); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sensor reading data"})
		}
		updatedReading, err := db.UpdateSensorReading(reading)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update sensor reading"})
		}
		return c.JSON(updatedReading)
	}
}

func DeleteSensorReading(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sensor reading ID"})
		}
		if err := db.DeleteSensorReading(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete sensor reading"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllSensorReadings(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		readings, err := db.GetAllSensorReadings()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get all sensor readings"})
		}
		return c.JSON(readings)
	}
}
