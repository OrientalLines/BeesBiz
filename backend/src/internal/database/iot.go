package database

import (
	"fmt"

	types "github.com/orientallines/beesbiz/internal/types/db"
	"go.uber.org/zap"
)

func (db *DB) GetSensor(id int) (types.Sensor, error) {
	var sensor types.Sensor
	err := db.Get(&sensor, "SELECT * FROM sensor WHERE sensor_id = $1", id)
	if err != nil {
		zap.S().Error("Error getting sensor: ", err)
		return types.Sensor{}, fmt.Errorf("error getting sensor: %w", err)
	}
	return sensor, nil
}

func (db *DB) CreateSensor(sensor types.Sensor) (types.Sensor, error) {
	var createdSensor types.Sensor
	err := db.Get(&createdSensor, "INSERT INTO sensor (hive_id, sensor_type, last_reading, last_reading_time) VALUES ($1, $2, $3, $4) RETURNING *", sensor.HiveID, sensor.SensorType, sensor.LastReading, sensor.LastReadingTime)
	if err != nil {
		zap.S().Error("Error creating sensor: ", err)
		return types.Sensor{}, fmt.Errorf("error creating sensor: %w", err)
	}
	return createdSensor, nil
}

func (db *DB) UpdateSensor(sensor types.Sensor) (types.Sensor, error) {
	var updatedSensor types.Sensor
	err := db.Get(&updatedSensor, "UPDATE sensor SET hive_id = $1, sensor_type = $2, last_reading = $3, last_reading_time = $4 WHERE sensor_id = $5 RETURNING *", sensor.HiveID, sensor.SensorType, sensor.LastReading, sensor.LastReadingTime, sensor.SensorID)
	if err != nil {
		zap.S().Error("Error updating sensor: ", err)
		return types.Sensor{}, fmt.Errorf("error updating sensor: %w", err)
	}
	return updatedSensor, nil
}

func (db *DB) DeleteSensor(id int) error {
	_, err := db.Exec("DELETE FROM sensor WHERE sensor_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting sensor: ", err)
		return fmt.Errorf("error deleting sensor: %w", err)
	}
	return nil
}

func (db *DB) GetAllSensors() ([]types.Sensor, error) {
	var sensors []types.Sensor
	err := db.Select(&sensors, "SELECT * FROM sensor")
	if err != nil {
		zap.S().Error("Error getting all sensors: ", err)
		return []types.Sensor{}, fmt.Errorf("error getting all sensors: %w", err)
	}
	return sensors, nil
}

func (db *DB) GetSensorReading(id int) (types.SensorReading, error) {
	var reading types.SensorReading
	err := db.Get(&reading, "SELECT * FROM sensor_reading WHERE reading_id = $1", id)
	if err != nil {
		zap.S().Error("Error getting sensor reading: ", err)
		return types.SensorReading{}, fmt.Errorf("error getting sensor reading: %w", err)
	}
	return reading, nil
}

func (db *DB) CreateSensorReading(reading types.SensorReading) (types.SensorReading, error) {
	var createdReading types.SensorReading
	err := db.Get(&createdReading, "INSERT INTO sensor_reading (sensor_id, value, timestamp) VALUES ($1, $2, $3) RETURNING *", reading.SensorID, reading.Value, reading.Timestamp)
	if err != nil {
		zap.S().Error("Error creating sensor reading: ", err)
		return types.SensorReading{}, fmt.Errorf("error creating sensor reading: %w", err)
	}
	return createdReading, nil
}

func (db *DB) UpdateSensorReading(reading types.SensorReading) (types.SensorReading, error) {
	var updatedReading types.SensorReading
	err := db.Get(&updatedReading, "UPDATE sensor_reading SET sensor_id = $1, value = $2, timestamp = $3 WHERE reading_id = $4 RETURNING *", reading.SensorID, reading.Value, reading.Timestamp, reading.ReadingID)
	if err != nil {
		zap.S().Error("Error updating sensor reading: ", err)
		return types.SensorReading{}, fmt.Errorf("error updating sensor reading: %w", err)
	}
	return updatedReading, nil
}

func (db *DB) DeleteSensorReading(id int) error {
	_, err := db.Exec("DELETE FROM sensor_reading WHERE reading_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting sensor reading: ", err)
		return fmt.Errorf("error deleting sensor reading: %w", err)
	}
	return nil
}

func (db *DB) GetAllSensorReadings() ([]types.SensorReading, error) {
	var readings []types.SensorReading
	err := db.Select(&readings, "SELECT * FROM sensor_reading")
	if err != nil {
		zap.S().Error("Error getting all sensor readings: ", err)
		return []types.SensorReading{}, fmt.Errorf("error getting all sensor readings: %w", err)
	}
	return readings, nil
}
