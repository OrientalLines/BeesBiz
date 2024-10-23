package types

import "github.com/guregu/null"

type Sensor struct {
	SensorID        int       `json:"sensor_id,omitempty" db:"sensor_id"`
	HiveID          int       `json:"hive_id" db:"hive_id"`
	SensorType      string    `json:"sensor_type" db:"sensor_type"`
	LastReading     []byte    `json:"last_reading" db:"last_reading"`
	LastReadingTime null.Time `json:"last_reading_time" db:"last_reading_time"`
}

type SensorReading struct {
	ReadingID int       `json:"reading_id,omitempty" db:"reading_id"`
	SensorID  int       `json:"sensor_id" db:"sensor_id"`
	Value     []byte    `json:"value" db:"value"`
	Timestamp null.Time `json:"timestamp" db:"timestamp"`
}
