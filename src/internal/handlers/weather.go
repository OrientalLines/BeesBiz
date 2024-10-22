package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/orientallines/beesbiz/internal/database"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

// WeatherData handlers
func GetWeatherData(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid weather data ID"})
		}
		weatherData, err := db.GetWeatherData(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get weather data"})
		}
		return c.JSON(weatherData)
	}
}

func CreateWeatherData(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var weatherData types.WeatherData
		if err := c.BodyParser(&weatherData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid weather data"})
		}
		createdWeatherData, err := db.CreateWeatherData(weatherData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create weather data"})
		}
		return c.JSON(createdWeatherData)
	}
}

func UpdateWeatherData(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var weatherData types.WeatherData
		if err := c.BodyParser(&weatherData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid weather data"})
		}
		updatedWeatherData, err := db.UpdateWeatherData(weatherData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update weather data"})
		}
		return c.JSON(updatedWeatherData)
	}
}

func DeleteWeatherData(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid weather data ID"})
		}
		if err := db.DeleteWeatherData(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete weather data"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllWeatherData(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		weatherDataList, err := db.GetAllWeatherData()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get all weather data"})
		}
		return c.JSON(weatherDataList)
	}
}
