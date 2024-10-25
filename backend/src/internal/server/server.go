package server

import (
	"context"
	"fmt"

	"github.com/orientallines/beesbiz/internal/database"
	"github.com/orientallines/beesbiz/internal/grpc"
	"github.com/orientallines/beesbiz/internal/rabbitmq"
	"github.com/orientallines/beesbiz/internal/rest"
)

type Server struct {
	grpcServer   *grpc.Server
	restServer   *rest.Server
	rabbitServer *rabbitmq.Server
	// tikvServer   *tikv.Server
}

// NewServer creates a new Server
// func NewServer(db *database.DB, rmq *rabbitmq.RabbitMQ, tikvClient *tikv.TiKV) (*Server, error) {
func NewServer(db *database.DB, rmq *rabbitmq.RabbitMQ) (*Server, error) {
	rabbitServer, err := rabbitmq.NewServer(rmq, db)
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ server: %w", err)
	}

	// tikvServer := tikv.NewServer(tikvClient.GetClient(), db)

	return &Server{
		grpcServer:   grpc.NewServer(db),
		restServer:   rest.NewServer(db),
		rabbitServer: rabbitServer,
		// tikvServer:   tikvServer,
	}, nil
}

// SetupAndRun sets up the servers and runs them
func (s *Server) SetupAndRun(grpcAddress, restAddress string) error {
	s.grpcServer.RegisterServices()
	s.restServer.SetupRoutes()

	errChan := make(chan error, 4)

	go func() {
		errChan <- s.grpcServer.Run(grpcAddress)
	}()

	go func() {
		errChan <- s.restServer.Run(restAddress)
	}()

	go func() {
		errChan <- s.rabbitServer.Run()
	}()

	// Start TiKV server and handle its errors
	// go func() {
	// 	if err := s.tikvServer.Run(); err != nil {
	// 		errChan <- fmt.Errorf("tikv server startup error: %w", err)
	// 		return
	// 	}

	// 	// Forward TiKV errors to main error channel
	// 	for err := range s.tikvServer.ErrCh {
	// 		if err != nil {
	// 			zap.L().Error("TiKV error", zap.Error(err))
	// 			// Only send fatal errors to errChan
	// 			if isFatalError(err) {
	// 				errChan <- fmt.Errorf("tikv error: %w", err)
	// 				return
	// 			}
	// 		}
	// 	}
	// }()

	// Return the first error that occurs
	return <-errChan
}

// Shutdown shuts down the servers
func (s *Server) Shutdown(ctx context.Context) error {
	errChan := make(chan error, 4)

	go func() {
		errChan <- s.grpcServer.Shutdown(ctx)
	}()

	go func() {
		errChan <- s.restServer.Shutdown(ctx)
	}()

	go func() {
		errChan <- s.rabbitServer.Shutdown(ctx)
	}()

	// go func() {
	// 	errChan <- s.tikvServer.Shutdown()
	// }()

	for i := 0; i < 4; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

// Helper function to determine if an error is fatal
func isFatalError(err error) bool {
	// TODO: @KXRXH
	// Add your fatal error conditions here
	// For example: connection loss, context cancellation, etc.
	return true
}
