package tikv

// import (
// 	"context"
// 	"log"
// 	"net"
//
//
//
// 

// 	"google.golang.org/grpc"

// 	pb "github.com/orientallines/beesbiz/internal/tikv/proto"
// )

// type server struct {
// 	pb.UnimplementedTiKVServiceServer
// }

// func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
// 	// Implement Get logic
// 	return &pb.GetResponse{Value: []byte("Value for " + req.Key)}, nil
// }

// func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
// 	// Implement Set logic
// 	return &pb.SetResponse{Success: true}, nil
// }

// func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
// 	// Implement Delete logic
// 	return &pb.DeleteResponse{Success: true}, nil
// }

// func StartGRPCServer() {
// 	lis, err := net.Listen("tcp", ":50051")
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}
// 	s := grpc.NewServer()
// 	pb.RegisterTiKVServiceServer(s, &server{})
// 	log.Printf("gRPC server listening at %v", lis.Addr())
// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: %v", err)
// 	}
// }
