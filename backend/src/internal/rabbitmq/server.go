package rabbitmq

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"go.uber.org/zap"

	"github.com/orientallines/beesbiz/internal/database"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

type Server struct {
	rmq *RabbitMQ
	db  *database.DB
}

const (
	HiveQueue          = "hive_queue"
	SensorQueue        = "sensor_queue"
	SensorReadingQueue = "sensor_reading_queue"
	DeleteSensorQueue  = "sensor_delete_queue"
)

// NewServer creates a new RabbitMQ server
func NewServer(rmq *RabbitMQ, db *database.DB) (*Server, error) {
	server := &Server{
		rmq: rmq,
		db:  db,
	}

	// Declare queues
	ch := rmq.GetChannel()
	for _, queue := range []string{SensorQueue, SensorReadingQueue} {
		_, err := ch.QueueDeclare(
			queue, // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return nil, fmt.Errorf("failed to declare queue %s: %w", queue, err)
		}
	}

	return server, nil
}

// Run starts the RabbitMQ server
func (s *Server) Run() error {
	// Start consumers
	go s.consumeSensorData()
	go s.consumeSensorReadingData()

	return nil
}

// Shutdown closes the RabbitMQ connection and channel
func (s *Server) Shutdown(ctx context.Context) error {
	return s.rmq.Close()
}

// TODO: Test this
func (s *Server) consumeSensorData() {
	ch := s.rmq.GetChannel()
	msgs, err := ch.Consume(
		SensorQueue, // queue
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		zap.L().Error("Failed to register sensor consumer", zap.Error(err))
		return
	}

	for msg := range msgs {
		var sensor types.Sensor
		if err := sonic.Unmarshal(msg.Body, &sensor); err != nil {
			zap.L().Error("Failed to unmarshal sensor data", zap.Error(err))
			msg.Nack(false, false) // Negative acknowledgement, requeue
			continue
		}

		// If sensor has an ID, update it; otherwise create new
		var err error
		if sensor.SensorID != 0 {
			_, err = s.db.UpdateSensor(sensor)
		} else {
			_, err = s.db.CreateSensor(sensor)
		}

		if err != nil {
			zap.L().Error("Failed to save sensor data", zap.Error(err))
			msg.Nack(false, false)
			continue
		}

		msg.Ack(false)
		zap.L().Info("Successfully processed sensor data",
			zap.Int("sensor_id", sensor.SensorID),
			zap.Int("hive_id", sensor.HiveID))
	}
}

// TODO: Test this
func (s *Server) consumeSensorReadingData() {
	ch := s.rmq.GetChannel()
	msgs, err := ch.Consume(
		SensorReadingQueue, // queue
		"",                 // consumer
		false,              // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		zap.L().Error("Failed to register sensor reading consumer", zap.Error(err))
		return
	}

	for msg := range msgs {
		var reading types.SensorReading
		if err := sonic.Unmarshal(msg.Body, &reading); err != nil {
			zap.L().Error("Failed to unmarshal sensor reading data", zap.Error(err))
			msg.Nack(false, false)
			continue
		}

		// Create new sensor reading
		_, err = s.db.CreateSensorReading(reading)
		if err != nil {
			zap.L().Error("Failed to save sensor reading data", zap.Error(err))
			msg.Nack(false, false)
			continue
		}

		// Update the sensor's last reading
		sensor, err := s.db.GetSensor(reading.SensorID)
		if err != nil {
			zap.L().Error("Failed to get sensor for updating last reading",
				zap.Error(err),
				zap.Int("sensor_id", reading.SensorID))
			msg.Nack(false, false)
			continue
		}

		sensor.LastReading = reading.Value
		sensor.LastReadingTime = reading.Timestamp

		_, err = s.db.UpdateSensor(sensor)
		if err != nil {
			zap.L().Error("Failed to update sensor's last reading",
				zap.Error(err),
				zap.Int("sensor_id", reading.SensorID))
			msg.Nack(false, false)
			continue
		}

		msg.Ack(false)
		zap.L().Info("Successfully processed sensor reading",
			zap.Int("sensor_id", reading.SensorID),
			zap.String("sensor_value", string(reading.Value)))
	}
}
