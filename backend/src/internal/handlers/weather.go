package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/orientallines/beesbiz/internal/database"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

// WeatherData handlers
func GetWeatherData(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid weather data ID: %v", err)})
		}
		weatherData, err := db.GetWeatherData(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get weather data: %v", err)})
		}
		return c.JSON(weatherData)
	}
}

func CreateWeatherData(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var weatherData types.WeatherData
		if err := c.BodyParser(&weatherData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid weather data: %v", err)})
		}
		createdWeatherData, err := db.CreateWeatherData(weatherData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create weather data: %v", err)})
		}
		return c.JSON(createdWeatherData)
	}
}

func UpdateWeatherData(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var weatherData types.WeatherData
		if err := c.BodyParser(&weatherData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid weather data: %v", err)})
		}
		updatedWeatherData, err := db.UpdateWeatherData(weatherData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update weather data: %v", err)})
		}
		return c.JSON(updatedWeatherData)
	}
}

func DeleteWeatherData(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid weather data ID: %v", err)})
		}
		if err := db.DeleteWeatherData(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete weather data: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllWeatherData(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		weatherDataList, err := db.GetAllWeatherData()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all weather data: %v", err)})
		}
		return c.JSON(weatherDataList)
	}
}
