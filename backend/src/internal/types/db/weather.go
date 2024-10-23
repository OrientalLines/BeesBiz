package types

import "github.com/guregu/null"

type WeatherData struct {
	WeatherID     int       `json:"weather_id,omitempty" db:"weather_id"`
	RegionID      int       `json:"region_id" db:"region_id"`
	Date          null.Time `json:"date" db:"date"`
	Temperature   float32   `json:"temperature" db:"temperature"`
	Humidity      float32   `json:"humidity" db:"humidity"`
	WindSpeed     float32   `json:"wind_speed" db:"wind_speed"`
	Precipitation float32   `json:"precipitation" db:"precipitation"`
}
