package types

import "time"

type WeatherData struct {
	WeatherID     int       `json:"weather_id" db:"weather_id"`
	RegionID      int       `json:"region_id" db:"region_id"`
	Date          time.Time `json:"date" db:"date"`
	Temperature   int       `json:"temperature" db:"temperature"`
	Humidity      int       `json:"humidity" db:"humidity"`
	WindSpeed     int       `json:"wind_speed" db:"wind_speed"`
	Precipitation int       `json:"precipitation" db:"precipitation"`
}
