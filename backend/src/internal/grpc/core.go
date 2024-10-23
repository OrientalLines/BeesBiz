package grpc

import (
	"context"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"database/sql"

	"github.com/orientallines/beesbiz/internal/database"
	pb "github.com/orientallines/beesbiz/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	server *grpc.Server
	db     *database.DB
	pb.UnimplementedBeeManagementServiceServer
}

func NewServer(db *database.DB) *Server {
	return &Server{
		server: grpc.NewServer(),
		db:     db,
	}
}

func (s *Server) Run(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	return s.server.Serve(lis)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.server.GracefulStop()
	return nil
}

func (s *Server) RegisterServices() {
	pb.RegisterBeeManagementServiceServer(s.server, s)
}

func (s *Server) GetTotalHoneyHarvested(ctx context.Context, req *pb.GetTotalHoneyHarvestedRequest) (*pb.GetTotalHoneyHarvestedResponse, error) {
	var totalHoney float64
	err := s.db.GetContext(ctx, &totalHoney, "SELECT get_total_honey_harvested($1, $2, $3)", req.HiveId, req.StartDate, req.EndDate)
	if err != nil {
		zap.S().Error("Error getting total honey harvested: ", err)
		return nil, err
	}
	return &pb.GetTotalHoneyHarvestedResponse{TotalHoney: totalHoney}, nil
}

func (s *Server) AddObservation(ctx context.Context, req *pb.AddObservationRequest) (*emptypb.Empty, error) {
	_, err := s.db.ExecContext(ctx, "CALL add_observation($1, $2, $3, $4)", req.HiveId, req.ObservationDate, req.Description, req.Recommendations)
	if err != nil {
		zap.S().Error("Error adding observation: ", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetCommunityHealthStatus(ctx context.Context, req *pb.GetCommunityHealthStatusRequest) (*pb.GetCommunityHealthStatusResponse, error) {
	var healthStatus string
	err := s.db.GetContext(ctx, &healthStatus, "SELECT get_community_health_status($1)", req.CommunityId)
	if err != nil {
		zap.S().Error("Error getting community health status: ", err)
		return nil, err
	}
	return &pb.GetCommunityHealthStatusResponse{HealthStatus: healthStatus}, nil
}

func (s *Server) UpdateHiveStatus(ctx context.Context, req *pb.UpdateHiveStatusRequest) (*emptypb.Empty, error) {
	_, err := s.db.ExecContext(ctx, "CALL update_hive_status($1, $2)", req.HiveId, req.NewStatus)
	if err != nil {
		zap.S().Error("Error updating hive status: ", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetAvgTemperature(ctx context.Context, req *pb.GetAvgTemperatureRequest) (*pb.GetAvgTemperatureResponse, error) {
	var avgTemp float64
	err := s.db.GetContext(ctx, &avgTemp, "SELECT get_avg_temperature($1, $2)", req.RegionId, req.Days)
	if err != nil {
		zap.S().Error("Error getting average temperature: ", err)
		return nil, err
	}
	return &pb.GetAvgTemperatureResponse{AvgTemperature: avgTemp}, nil
}

func (s *Server) AssignMaintenancePlan(ctx context.Context, req *pb.AssignMaintenancePlanRequest) (*emptypb.Empty, error) {
	_, err := s.db.ExecContext(ctx, "CALL assign_maintenance_plan($1, $2)", req.PlanId, req.UserId)
	if err != nil {
		zap.S().Error("Error assigning maintenance plan: ", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) HasRegionAccess(ctx context.Context, req *pb.HasRegionAccessRequest) (*pb.HasRegionAccessResponse, error) {
	var hasAccess bool
	err := s.db.GetContext(ctx, &hasAccess, "SELECT has_region_access($1, $2)", req.UserId, req.RegionId)
	if err != nil {
		zap.S().Error("Error checking region access: ", err)
		return nil, err
	}
	return &pb.HasRegionAccessResponse{HasAccess: hasAccess}, nil
}

func (s *Server) RegisterIncident(ctx context.Context, req *pb.RegisterIncidentRequest) (*emptypb.Empty, error) {
	_, err := s.db.ExecContext(ctx, "CALL register_incident($1, $2, $3, $4)", req.HiveId, req.IncidentDate, req.Description, req.Severity)
	if err != nil {
		zap.S().Error("Error registering incident: ", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetLatestSensorReading(ctx context.Context, req *pb.GetLatestSensorReadingRequest) (*pb.GetLatestSensorReadingResponse, error) {
	var value []byte
	var timestamp time.Time
	err := s.db.QueryRowContext(ctx, "SELECT * FROM get_latest_sensor_reading($1, $2)", req.HiveId, req.SensorType).
		Scan(&value, &timestamp)
	if err != nil {
		zap.S().Error("Error getting latest sensor reading: ", err)
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "No sensor reading found for the specified hive and sensor type.")
		}
		// For other errors, return a generic Unknown error
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &pb.GetLatestSensorReadingResponse{
		Value:     value,
		Timestamp: timestamp.Format(time.RFC3339),
	}, nil
}

func (s *Server) CreateProductionReport(ctx context.Context, req *pb.CreateProductionReportRequest) (*emptypb.Empty, error) {
	_, err := s.db.ExecContext(ctx, "CALL create_production_report($1, $2, $3)", req.ApiaryId, req.StartDate, req.EndDate)
	if err != nil {
		zap.S().Error("Error creating production report: ", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
