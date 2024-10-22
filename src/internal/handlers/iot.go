package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/orientallines/beesbiz/internal/database"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

// Sensor handlers
func GetSensor(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid sensor ID: %v", err)})
		}
		sensor, err := db.GetSensor(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get sensor: %v", err)})
		}
		return c.JSON(sensor)
	}
}

func CreateSensor(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var sensor types.Sensor
		if err := c.BodyParser(&sensor); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid sensor data: %v", err)})
		}
		createdSensor, err := db.CreateSensor(sensor)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create sensor: %v", err)})
		}
		return c.JSON(createdSensor)
	}
}

func UpdateSensor(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var sensor types.Sensor
		if err := c.BodyParser(&sensor); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid sensor data: %v", err)})
		}
		updatedSensor, err := db.UpdateSensor(sensor)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update sensor: %v", err)})
		}
		return c.JSON(updatedSensor)
	}
}

func DeleteSensor(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid sensor ID: %v", err)})
		}
		if err := db.DeleteSensor(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete sensor: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllSensors(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sensors, err := db.GetAllSensors()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all sensors: %v", err)})
		}
		return c.JSON(sensors)
	}
}

// SensorReading handlers
func GetSensorReading(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid sensor reading ID: %v", err)})
		}
		reading, err := db.GetSensorReading(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get sensor reading: %v", err)})
		}
		return c.JSON(reading)
	}
}

func CreateSensorReading(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reading types.SensorReading
		if err := c.BodyParser(&reading); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid sensor reading data: %v", err)})
		}
		createdReading, err := db.CreateSensorReading(reading)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create sensor reading: %v", err)})
		}
		return c.JSON(createdReading)
	}
}

func UpdateSensorReading(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reading types.SensorReading
		if err := c.BodyParser(&reading); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid sensor reading data: %v", err)})
		}
		updatedReading, err := db.UpdateSensorReading(reading)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update sensor reading: %v", err)})
		}
		return c.JSON(updatedReading)
	}
}

func DeleteSensorReading(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid sensor reading ID: %v", err)})
		}
		if err := db.DeleteSensorReading(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete sensor reading: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllSensorReadings(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		readings, err := db.GetAllSensorReadings()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all sensor readings: %v", err)})
		}
		return c.JSON(readings)
	}
}
