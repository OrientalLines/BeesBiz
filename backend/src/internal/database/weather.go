package database

import (
	"fmt"

	"go.uber.org/zap"

	types "github.com/orientallines/beesbiz/internal/types/db"
)

// GetWeatherData gets weather data by ID
//
// It returns the weather data and an error
func (db *DB) GetWeatherData(id int) (types.WeatherData, error) {
	var weather types.WeatherData
	err := db.Get(&weather, "SELECT * FROM weather_data WHERE weather_id = $1", id)
	if err != nil {
		zap.S().Error("Error getting weather data: ", err)
		return types.WeatherData{}, fmt.Errorf("error getting weather data: %w", err)
	}
	return weather, nil
}

// CreateWeatherData creates a new weather data
//
// It returns the created weather data and an error
func (db *DB) CreateWeatherData(weather types.WeatherData) (types.WeatherData, error) {
	var createdWeather types.WeatherData
	err := db.Get(&createdWeather, "INSERT INTO weather_data (region_id, date, temperature, humidity, wind_speed, precipitation) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *", weather.RegionID, weather.Date, weather.Temperature, weather.Humidity, weather.WindSpeed, weather.Precipitation)
	if err != nil {
		zap.S().Error("Error creating weather data: ", err)
		return types.WeatherData{}, fmt.Errorf("error creating weather data: %w", err)
	}
	return createdWeather, nil
}

// UpdateWeatherData updates a weather data
//
// It returns the updated weather data and an error
func (db *DB) UpdateWeatherData(weather types.WeatherData) (types.WeatherData, error) {
	var updatedWeather types.WeatherData
	err := db.Get(&updatedWeather, "UPDATE weather_data SET region_id = $1, date = $2, temperature = $3, humidity = $4, wind_speed = $5, precipitation = $6 WHERE weather_id = $7 RETURNING *", weather.RegionID, weather.Date, weather.Temperature, weather.Humidity, weather.WindSpeed, weather.Precipitation, weather.WeatherID)
	if err != nil {
		zap.S().Error("Error updating weather data: ", err)
		return types.WeatherData{}, fmt.Errorf("error updating weather data: %w", err)
	}
	return updatedWeather, nil
}

// DeleteWeatherData deletes a weather data
//
// It returns an error
func (db *DB) DeleteWeatherData(id int) error {
	_, err := db.Exec("DELETE FROM weather_data WHERE weather_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting weather data: ", err)
		return fmt.Errorf("error deleting weather data: %w", err)
	}
	return nil
}

// GetAllWeatherData gets all weather data
// 
// It returns a list of weather data and an error
func (db *DB) GetAllWeatherData() ([]types.WeatherData, error) {
	var weatherDataList []types.WeatherData
	err := db.Select(&weatherDataList, "SELECT * FROM weather_data")
	if err != nil {
		zap.S().Error("Error getting all weather data: ", err)
		return []types.WeatherData{}, fmt.Errorf("error getting all weather data: %w", err)
	}
	return weatherDataList, nil
}
