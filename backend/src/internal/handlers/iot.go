package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/orientallines/beesbiz/internal/database"
	"github.com/orientallines/beesbiz/internal/rabbitmq"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

// Sensor handlers

// GetSensor gets a sensor by ID
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

// CreateSensor creates a new sensor
func CreateSensor(db *database.DB, rmq *rabbitmq.RabbitMQ) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var sensor types.Sensor
		if err := c.BodyParser(&sensor); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid sensor data: %v", err)})
		}

		// Create sensor in database
		createdSensor, err := db.CreateSensor(sensor)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create sensor: %v", err)})
		}

		// Create message in format expected by IoT service
		iotMessage := types.HiveSensorMessage{
			HiveID:  createdSensor.HiveID,
			Sensors: []types.Sensor{createdSensor},
		}

		// Publish sensor creation message to RabbitMQ
		if err := rmq.PublishMessage(rabbitmq.HiveQueue, iotMessage); err != nil {
			zap.L().Error("Failed to publish sensor creation message",
				zap.Error(err),
				zap.Int("sensor_id", createdSensor.SensorID),
				zap.Int("hive_id", createdSensor.HiveID))
		}

		return c.JSON(createdSensor)
	}
}

// UpdateSensor updates a sensor
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

// DeleteSensor deletes a sensor
func DeleteSensor(db *database.DB, rmq *rabbitmq.RabbitMQ) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid sensor ID: %v", err)})
		}
		sensor, err := db.GetSensor(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get sensor: %v", err)})
		}
		if err := db.DeleteSensor(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete sensor: %v", err)})
		}
		deleteMsg := types.DeleteSensor{
			HiveID:   sensor.HiveID,
			SensorID: id,
		}
		if err := rmq.PublishMessage(rabbitmq.DeleteSensorQueue, deleteMsg); err != nil {
			zap.L().Error("Failed to publish sensor deletion message",
				zap.Error(err),
				zap.Int("sensor_id", id),
				zap.Int("hive_id", sensor.HiveID))
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// GetAllSensors gets all sensors
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

// GetSensorReading gets a sensor reading by ID
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

// CreateSensorReading creates a new sensor reading
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

// UpdateSensorReading updates a sensor reading
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

// DeleteSensorReading deletes a sensor reading
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

// GetAllSensorReadings gets all sensor readings
func GetAllSensorReadings(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		readings, err := db.GetAllSensorReadings()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all sensor readings: %v", err)})
		}
		return c.JSON(readings)
	}
}
