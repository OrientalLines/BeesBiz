package server

import (
	"context"

	"github.com/orientallines/beesbiz/internal/database"
	"github.com/orientallines/beesbiz/internal/grpc"
	"github.com/orientallines/beesbiz/internal/rest"
)

type Server struct {
	grpcServer *grpc.Server
	restServer *rest.Server
}

// NewServer creates a new Server
func NewServer(db *database.DB) *Server {
	return &Server{
		grpcServer: grpc.NewServer(db),
		restServer: rest.NewServer(db),
	}
}

// SetupAndRun sets up the servers and runs them
func (s *Server) SetupAndRun(grpcAddress, restAddress string) error {
	s.grpcServer.RegisterServices()
	s.restServer.SetupRoutes()

	errChan := make(chan error, 2)

	go func() {
		errChan <- s.grpcServer.Run(grpcAddress)
	}()

	go func() {
		errChan <- s.restServer.Run(restAddress)
	}()

	// Return the first error that occurs
	return <-errChan
}

// Shutdown shuts down the servers
func (s *Server) Shutdown(ctx context.Context) error {
	errChan := make(chan error, 2)

	go func() {
		errChan <- s.grpcServer.Shutdown(ctx)
	}()

	go func() {
		errChan <- s.restServer.Shutdown(ctx)
	}()

	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}
