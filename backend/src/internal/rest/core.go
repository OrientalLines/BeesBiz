package rest

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/orientallines/beesbiz/internal/database"
	"github.com/orientallines/beesbiz/internal/handlers"
)

// Server is a wrapper around fiber.App
type Server struct {
	app *fiber.App
	db  *database.DB
}

// NewServer creates a new Server
func NewServer(db *database.DB) *Server {

	return &Server{
		app: fiber.New(),
		db:  db,
	}
}

// SetupRoutes sets up the routes for the server
func (s *Server) SetupRoutes() {
	s.app.Use(requestid.New())
	// s.app.Use(logger.New(logger.Config{
	// 	Format: "[${time}] ${status} - ${method} ${path}\n",
	// }))
	s.app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/livez",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		ReadinessEndpoint: "/readyz",
	}))

	s.app.Get("/apiary/:id", handlers.GetApiary(s.db))
	s.app.Post("/apiary", handlers.CreateApiary(s.db))
	s.app.Put("/apiary", handlers.UpdateApiary(s.db))
	s.app.Delete("/apiary/:id", handlers.DeleteApiary(s.db))
	s.app.Get("/apiaries", handlers.GetAllApiaries(s.db))

	// Hive routes
	s.app.Get("/hives", handlers.GetAllHives(s.db))
	s.app.Post("/hive", handlers.CreateHive(s.db))
	s.app.Put("/hive", handlers.UpdateHive(s.db))
	s.app.Delete("/hive/:id", handlers.DeleteHive(s.db))
	s.app.Get("/apiary/:apiaryID/hives", handlers.GetAllHivesByApiaryID(s.db))

	// BeeCommunity routes
	s.app.Get("/bee-communities", handlers.GetAllBeeCommunities(s.db))
	s.app.Post("/bee-community", handlers.CreateBeeCommunity(s.db))
	s.app.Put("/bee-community", handlers.UpdateBeeCommunity(s.db))
	s.app.Delete("/bee-community/:id", handlers.DeleteBeeCommunity(s.db))
	s.app.Get("/hive/:hiveID/bee-communities", handlers.GetAllBeeCommunitiesByHiveID(s.db))

	// HoneyHarvest routes
	s.app.Get("/honey-harvest/:id", handlers.GetHoneyHarvest(s.db))
	s.app.Post("/honey-harvest", handlers.CreateHoneyHarvest(s.db))
	s.app.Put("/honey-harvest", handlers.UpdateHoneyHarvest(s.db))
	s.app.Delete("/honey-harvest/:id", handlers.DeleteHoneyHarvest(s.db))
	s.app.Get("/honey-harvests", handlers.GetAllHoneyHarvests(s.db))

	// Region routes
	s.app.Get("/region/:id", handlers.GetRegion(s.db))
	s.app.Post("/region", handlers.CreateRegion(s.db))
	s.app.Put("/region", handlers.UpdateRegion(s.db))
	s.app.Delete("/region/:id", handlers.DeleteRegion(s.db))
	s.app.Get("/regions", handlers.GetAllRegions(s.db))

	// AllowedRegion routes
	s.app.Get("/allowed-region/:id", handlers.GetAllowedRegion(s.db))
	s.app.Post("/allowed-region", handlers.CreateAllowedRegion(s.db))
	s.app.Put("/allowed-region", handlers.UpdateAllowedRegion(s.db))
	s.app.Delete("/allowed-region/:id", handlers.DeleteAllowedRegion(s.db))
	s.app.Get("/allowed-regions", handlers.GetAllAllowedRegions(s.db))

	// RegionApiary routes
	s.app.Get("/region-apiary/:id", handlers.GetRegionApiary(s.db))
	s.app.Post("/region-apiary", handlers.CreateRegionApiary(s.db))
	s.app.Put("/region-apiary", handlers.UpdateRegionApiary(s.db))
	s.app.Delete("/region-apiary/:id", handlers.DeleteRegionApiary(s.db))
	s.app.Get("/region-apiaries", handlers.GetAllRegionApiaries(s.db))

	// ObservationLog routes
	s.app.Get("/observation-log/:id", handlers.GetObservationLog(s.db))
	s.app.Post("/observation-log", handlers.CreateObservationLog(s.db))
	s.app.Put("/observation-log", handlers.UpdateObservationLog(s.db))
	s.app.Delete("/observation-log/:id", handlers.DeleteObservationLog(s.db))
	s.app.Get("/observation-logs", handlers.GetAllObservationLogs(s.db))

	// MaintenancePlan routes
	s.app.Get("/maintenance-plan/:id", handlers.GetMaintenancePlan(s.db))
	s.app.Post("/maintenance-plan", handlers.CreateMaintenancePlan(s.db))
	s.app.Put("/maintenance-plan", handlers.UpdateMaintenancePlan(s.db))
	s.app.Delete("/maintenance-plan/:id", handlers.DeleteMaintenancePlan(s.db))
	s.app.Get("/maintenance-plans", handlers.GetAllMaintenancePlans(s.db))

	// Incident routes
	s.app.Get("/incident/:id", handlers.GetIncident(s.db))
	s.app.Post("/incident", handlers.CreateIncident(s.db))
	s.app.Put("/incident", handlers.UpdateIncident(s.db))
	s.app.Delete("/incident/:id", handlers.DeleteIncident(s.db))
	s.app.Get("/incidents", handlers.GetAllIncidents(s.db))

	// VeterinaryPassport routes
	s.app.Get("/veterinary-passport/:id", handlers.GetVeterinaryPassport(s.db))
	s.app.Post("/veterinary-passport", handlers.CreateVeterinaryPassport(s.db))
	s.app.Put("/veterinary-passport", handlers.UpdateVeterinaryPassport(s.db))
	s.app.Delete("/veterinary-passport/:id", handlers.DeleteVeterinaryPassport(s.db))
	s.app.Get("/veterinary-passports", handlers.GetAllVeterinaryPassports(s.db))

	// VeterinaryRecord routes
	s.app.Get("/veterinary-record/:id", handlers.GetVeterinaryRecord(s.db))
	s.app.Post("/veterinary-record", handlers.CreateVeterinaryRecord(s.db))
	s.app.Put("/veterinary-record", handlers.UpdateVeterinaryRecord(s.db))
	s.app.Delete("/veterinary-record/:id", handlers.DeleteVeterinaryRecord(s.db))
	s.app.Get("/veterinary-records", handlers.GetAllVeterinaryRecords(s.db))

	// User routes
	s.app.Get("/user/:id", handlers.GetUser(s.db))
	s.app.Post("/user", handlers.CreateUser(s.db))
	s.app.Put("/user", handlers.UpdateUser(s.db))
	s.app.Delete("/user/:id", handlers.DeleteUser(s.db))
	s.app.Get("/users", handlers.GetAllUsers(s.db))

	// ProductionReport routes
	s.app.Get("/production-report/:id", handlers.GetProductionReport(s.db))
	s.app.Post("/production-report", handlers.CreateProductionReport(s.db))
	s.app.Put("/production-report", handlers.UpdateProductionReport(s.db))
	s.app.Delete("/production-report/:id", handlers.DeleteProductionReport(s.db))
	s.app.Get("/production-reports", handlers.GetAllProductionReports(s.db))

	// Sensor routes
	s.app.Get("/sensor/:id", handlers.GetSensor(s.db))
	s.app.Post("/sensor", handlers.CreateSensor(s.db))
	s.app.Put("/sensor", handlers.UpdateSensor(s.db))
	s.app.Delete("/sensor/:id", handlers.DeleteSensor(s.db))
	s.app.Get("/sensors", handlers.GetAllSensors(s.db))

	// SensorReading routes
	s.app.Get("/sensor-reading/:id", handlers.GetSensorReading(s.db))
	s.app.Post("/sensor-reading", handlers.CreateSensorReading(s.db))
	s.app.Put("/sensor-reading", handlers.UpdateSensorReading(s.db))
	s.app.Delete("/sensor-reading/:id", handlers.DeleteSensorReading(s.db))
	s.app.Get("/sensor-readings", handlers.GetAllSensorReadings(s.db))

	// WeatherData routes
	s.app.Get("/weather-data/:id", handlers.GetWeatherData(s.db))
	s.app.Post("/weather-data", handlers.CreateWeatherData(s.db))
	s.app.Put("/weather-data", handlers.UpdateWeatherData(s.db))
	s.app.Delete("/weather-data/:id", handlers.DeleteWeatherData(s.db))
	s.app.Get("/weather-data-all", handlers.GetAllWeatherData(s.db))

}

// Run starts the server
func (s *Server) Run(address string) error {
	return s.app.Listen(address)
}

// Shutdown shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
