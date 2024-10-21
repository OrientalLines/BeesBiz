package server

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/orientallines/beesbiz/internal/grpc"
	"github.com/orientallines/beesbiz/internal/rest"
)

type Server struct {
	grpcServer *grpc.Server
	restServer *rest.Server
}

func NewServer(db *sqlx.DB) *Server {
	return &Server{
		grpcServer: grpc.NewServer(db),
		restServer: rest.NewServer(),
	}
}

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

func (s *Server) Shutdown(ctx context.Context) error {
	errChan := make(chan error, 2)

	go func() {
		errChan <- s.grpcServer.Shutdown(ctx)
	}()

	go func() {
		errChan <- s.restServer.Shutdown(ctx)
	}()

	// Wait for both servers to shut down
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}
