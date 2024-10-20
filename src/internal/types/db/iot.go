package types

type Sensor struct {
	SensorID        int    `json:"sensor_id" db:"sensor_id"`
	HiveID          int    `json:"hive_id" db:"hive_id"`
	SensorType      string `json:"sensor_type" db:"sensor_type"`
	LastReading     int    `json:"last_reading" db:"last_reading"`
	LastReadingTime string `json:"last_reading_time" db:"last_reading_time"`
}

type SensorReading struct {
	ReadingID int    `json:"reading_id" db:"reading_id"`
	SensorID  int    `json:"sensor_id" db:"sensor_id"`
	Value     int    `json:"value" db:"value"`
	Timestamp string `json:"timestamp" db:"timestamp"`
}
