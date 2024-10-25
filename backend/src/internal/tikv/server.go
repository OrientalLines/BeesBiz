package tikv

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/orientallines/beesbiz/internal/database"
	types "github.com/orientallines/beesbiz/internal/types/db"
	"github.com/tikv/client-go/v2/rawkv"
	"go.uber.org/zap"
)

// Server handles caching with TiKV
type Server struct {
	client *rawkv.Client
	db     *database.DB
	ctx    context.Context
	cancel context.CancelFunc
	ErrCh  chan error // Changed to public for external access
}

// NewServer initializes a new TiKV caching server
func NewServer(client *rawkv.Client, db *database.DB) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		client: client,
		db:     db,
		ctx:    ctx,
		cancel: cancel,
		ErrCh:  make(chan error, 100), // Buffer for error handling
	}
}

// Run starts the caching server
func (s *Server) Run() error {
	go s.cacheSensorReadings()
	go s.cacheProductionReports()

	// Return nil to indicate successful startup
	return nil
}

// Shutdown gracefully shuts down the caching server
func (s *Server) Shutdown() error {
	s.cancel()
	close(s.ErrCh)
	return s.client.Close()
}

// errorHandler handles errors from the caching server
func (s *Server) errorHandler() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case err := <-s.ErrCh:
			zap.L().Error("TiKV cache error", zap.Error(err))
		}
	}
}

// TODO: Test this
// cacheSensorReadings caches the sensor_reading data
func (s *Server) cacheSensorReadings() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			// Fetch latest sensor readings from database
			readings, err := s.db.GetLatestSensorReadings(100) // example limit
			if err != nil {
				s.ErrCh <- fmt.Errorf("failed to fetch sensor readings: %w", err)
				continue
			}

			keys := make([][]byte, len(readings))
			values := make([][]byte, len(readings))
			ttls := make([]uint64, len(readings))

			for i, reading := range readings {
				key := fmt.Sprintf("sensor_reading:%d", reading.ReadingID)
				keys[i] = []byte(key)

				value, err := sonic.Marshal(reading)
				if err != nil {
					s.ErrCh <- fmt.Errorf("failed to marshal reading %d: %w", reading.ReadingID, err)
					continue
				}
				values[i] = value
				ttls[i] = 3600 // 1 hour TTL
			}

			if err := s.client.BatchPutWithTTL(s.ctx, keys, values, ttls); err != nil {
				s.ErrCh <- fmt.Errorf("failed to batch put sensor readings: %w", err)
			}
		}
	}
}

// TODO: Test this
// cacheProductionReports caches production report data
func (s *Server) cacheProductionReports() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			// Fetch latest production reports from database
			reports, err := s.db.GetRecentProductionReports(50) // example limit
			if err != nil {
				s.ErrCh <- fmt.Errorf("failed to fetch production reports: %w", err)
				continue
			}

			keys := make([][]byte, len(reports))
			values := make([][]byte, len(reports))

			for i, report := range reports {
				key := fmt.Sprintf("production_report:%d", report.ReportID)
				keys[i] = []byte(key)

				value, err := sonic.Marshal(report)
				if err != nil {
					s.ErrCh <- fmt.Errorf("failed to marshal report %d: %w", report.ReportID, err)
					continue
				}
				values[i] = value
			}
			// TODO: Implement BatchPutWithTTL
			if err := s.client.BatchPut(s.ctx, keys, values); err != nil {
				s.ErrCh <- fmt.Errorf("failed to batch put production reports: %w", err)
			}
		}
	}
}

// TODO: Implement this
// GetProductionReports retrieves production reports by IDs
func (s *Server) GetProductionReports(ids []int) ([]types.ProductionReport, error) {
	return nil, nil
}

// TODO: Test this and implement it in the rest server
// GetSensorReadings retrieves sensor readings by IDs
func (s *Server) GetSensorReadings(ids []int) ([]types.SensorReading, error) {
	keys := make([][]byte, len(ids))
	for i, id := range ids {
		keys[i] = []byte(fmt.Sprintf("sensor_reading:%d", id))
	}
	// TODO: Should consider TTL
	values, err := s.client.BatchGet(s.ctx, keys)
	if err != nil {
		return nil, fmt.Errorf("failed to batch get sensor readings: %w", err)
	}

	readings := make([]types.SensorReading, 0, len(values))
	for i, value := range values {
		if value == nil {
			continue
		}
		var reading types.SensorReading
		if err := sonic.Unmarshal(value, &reading); err != nil {
			s.ErrCh <- fmt.Errorf("failed to unmarshal reading for id %d: %w", ids[i], err)
			continue
		}
		readings = append(readings, reading)
	}

	return readings, nil
}

// TODO: Test this and implement it in the rest server
// DeleteSensorReadings deletes sensor readings by IDs
func (s *Server) DeleteSensorReadings(ids []int) error {
	keys := make([][]byte, len(ids))
	for i, id := range ids {
		keys[i] = []byte(fmt.Sprintf("sensor_reading:%d", id))
	}
	// TODO: Should consider TTL
	if err := s.client.BatchDelete(s.ctx, keys); err != nil {
		return fmt.Errorf("failed to batch delete sensor readings: %w", err)
	}
	return nil
}
