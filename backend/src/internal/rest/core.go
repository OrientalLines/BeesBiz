package rest

import (
	"context"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/orientallines/beesbiz/internal/config"
	"github.com/orientallines/beesbiz/internal/database"
	"github.com/orientallines/beesbiz/internal/handlers"
	"github.com/orientallines/beesbiz/internal/rabbitmq"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

// Server is a wrapper around fiber.App
type Server struct {
	app    *fiber.App
	db     *database.DB
	rmq    *rabbitmq.RabbitMQ
	jwtKey []byte
}

// NewServer creates a new Server
func NewServer(db *database.DB, rmq *rabbitmq.RabbitMQ) *Server {
	return &Server{
		app:    fiber.New(),
		db:     db,
		rmq:    rmq,
		jwtKey: []byte(config.GlobalConfig.JwtSecret),
	}
}

// SetupRoutes sets up the routes for the server
func (s *Server) SetupRoutes() {
	s.app.Use(requestid.New())
	s.app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))
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

	prometheus := fiberprometheus.New("beesbiz")
	prometheus.RegisterAt(s.app, "/metrics")
	prometheus.SetSkipPaths([]string{"/ping", "/livez", "/readyz"})

	s.app.Use(prometheus.Middleware)

	s.app.Use(cors.New(cors.Config{
		AllowOrigins:     "*", // Your frontend URL
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowCredentials: false,
	}))

	auth := s.app.Group("/auth")

	auth.Post("/login", handlers.Login(s.db, s.jwtKey))
	auth.Post("/register", handlers.Register(s.db))

	api := s.app.Group("/api", jwtMiddleware(s.jwtKey))

	// Apiary routes
	apiary := api.Group("/apiary", roleMiddleware(types.Worker, types.Manager, types.Admin))

	apiary.Get("/:id", handlers.GetApiary(s.db))
	apiary.Post("/", handlers.CreateApiary(s.db))
	apiary.Put("/", handlers.UpdateApiary(s.db))
	apiary.Delete("/:id", handlers.DeleteApiary(s.db))
	apiary.Get("/", handlers.GetAllApiaries(s.db))

	// Hive routes
	hive := api.Group("/hive", roleMiddleware(types.Worker, types.Manager, types.Admin))

	hive.Get("/", handlers.GetAllHives(s.db))
	hive.Post("/", handlers.CreateHive(s.db))
	hive.Put("/", handlers.UpdateHive(s.db))
	hive.Delete("/:id", handlers.DeleteHive(s.db, s.rmq))
	hive.Get("/:apiaryID/hives", handlers.GetAllHivesByApiaryID(s.db))

	// BeeCommunity routes
	beeCommunity := api.Group("/bee-community", roleMiddleware(types.Worker, types.Manager, types.Admin))

	beeCommunity.Get("/", handlers.GetAllBeeCommunities(s.db))
	beeCommunity.Post("/", handlers.CreateBeeCommunity(s.db))
	beeCommunity.Put("/", handlers.UpdateBeeCommunity(s.db))
	beeCommunity.Delete("/:id", handlers.DeleteBeeCommunity(s.db))
	beeCommunity.Get("/:hiveID/bee-communities", handlers.GetAllBeeCommunitiesByHiveID(s.db))

	// HoneyHarvest routes
	honeyHarvest := api.Group("/honey-harvest", roleMiddleware(types.Worker, types.Manager, types.Admin))

	honeyHarvest.Get("/:id", handlers.GetHoneyHarvest(s.db))
	honeyHarvest.Post("/", handlers.CreateHoneyHarvest(s.db))
	honeyHarvest.Put("/", handlers.UpdateHoneyHarvest(s.db))
	honeyHarvest.Delete("/:id", handlers.DeleteHoneyHarvest(s.db))
	honeyHarvest.Get("/", handlers.GetAllHoneyHarvests(s.db))

	// Region routes
	region := api.Group("/region", roleMiddleware(types.Worker, types.Manager, types.Admin))

	region.Get("/:id", handlers.GetRegion(s.db))
	region.Post("/", handlers.CreateRegion(s.db))
	region.Put("/", handlers.UpdateRegion(s.db))
	region.Delete("/:id", handlers.DeleteRegion(s.db))
	region.Get("/", handlers.GetAllRegions(s.db))

	// AllowedRegion routes
	allowedRegion := api.Group("/allowed-region", roleMiddleware(types.Manager, types.Admin))

	allowedRegion.Get("/user/:id", handlers.GetAllowedRegionsForUser(s.db))
	allowedRegion.Post("/", handlers.CreateAllowedRegion(s.db))
	allowedRegion.Put("/", handlers.UpdateAllowedRegion(s.db))
	allowedRegion.Delete("/:id", handlers.DeleteAllowedRegion(s.db))
	allowedRegion.Get("/", handlers.GetAllAllowedRegions(s.db))

	// RegionApiary routes
	regionApiary := api.Group("/region-apiary", roleMiddleware(types.Manager, types.Admin))

	regionApiary.Get("/:id", handlers.GetRegionApiary(s.db))
	regionApiary.Post("/", handlers.CreateRegionApiary(s.db))
	regionApiary.Put("/", handlers.UpdateRegionApiary(s.db))
	regionApiary.Delete("/:id", handlers.DeleteRegionApiary(s.db))
	regionApiary.Get("/", handlers.GetAllRegionApiaries(s.db))

	// User routes
	user := api.Group("/user", roleMiddleware(types.Admin, types.Manager))

	user.Get("/:id", handlers.GetUser(s.db))
	user.Post("/", handlers.CreateUser(s.db))
	user.Put("/", handlers.UpdateUser(s.db))
	user.Delete("/:id", handlers.DeleteUser(s.db))
	user.Get("/", handlers.GetAllUsers(s.db))
	user.Get("/:id/allowed-regions", handlers.GetUserAllowedRegions(s.db))
	user.Put("/role", handlers.ModifyUserRole(s.db))
	user.Put("/allowed-regions", handlers.ModifyUserAllowedRegions(s.db))

	// WorkerGroup routes
	workerGroup := api.Group("/worker-group", roleMiddleware(types.Admin, types.Manager))

	workerGroup.Get("/", handlers.GetAllWorkerGroups(s.db))
	workerGroup.Get("/:id", handlers.GetWorkerGroup(s.db))
	workerGroup.Post("/", handlers.CreateWorkerGroup(s.db))
	workerGroup.Get("/manager/:manager_id", handlers.GetWorkerGroupsByManager(s.db))
	workerGroup.Post("/:group_id/members", handlers.AddGroupMember(s.db))
	workerGroup.Delete("/:group_id/members/:worker_id", handlers.RemoveGroupMember(s.db))
	workerGroup.Get("/:group_id/members", handlers.GetGroupMembers(s.db))
	workerGroup.Get("/worker/:worker_id/groups", handlers.GetWorkerGroups(s.db))
	workerGroup.Delete("/:id", handlers.DeleteWorkerGroup(s.db))
	workerGroup.Put("/:id", handlers.UpdateWorkerGroup(s.db))

	// ProductionReport routes
	productionReport := api.Group("/production-report", roleMiddleware(types.Manager, types.Worker, types.Admin))

	productionReport.Get("/:id", handlers.GetProductionReport(s.db))
	productionReport.Post("/", handlers.CreateProductionReport(s.db))
	productionReport.Put("/", handlers.UpdateProductionReport(s.db))
	productionReport.Delete("/:id", handlers.DeleteProductionReport(s.db))
	productionReport.Get("/", handlers.GetAllProductionReports(s.db))
	productionReport.Get("/curated/:id", handlers.GetCuratedProductionReportsByUser(s.db))

	// Sensor routes
	sensor := api.Group("/sensor", roleMiddleware(types.Admin, types.Manager, types.Worker))

	sensor.Get("/:id", handlers.GetSensor(s.db))
	sensor.Post("/", handlers.CreateSensor(s.db, s.rmq))
	sensor.Put("/", handlers.UpdateSensor(s.db))
	sensor.Delete("/:id", handlers.DeleteSensor(s.db, s.rmq))
	sensor.Get("/", handlers.GetAllSensors(s.db))

	// SensorReading routes
	sensorReading := api.Group("/sensor-reading", roleMiddleware(types.Admin, types.Manager, types.Worker))

	sensorReading.Get("/:id", handlers.GetSensorReading(s.db))
	sensorReading.Post("/", handlers.CreateSensorReading(s.db))
	sensorReading.Put("/", handlers.UpdateSensorReading(s.db))
	sensorReading.Delete("/:id", handlers.DeleteSensorReading(s.db))
	sensorReading.Get("/", handlers.GetAllSensorReadings(s.db))

	// WeatherData routes
	weatherData := api.Group("/weather-data", roleMiddleware(types.Admin, types.Manager, types.Worker))

	weatherData.Get("/:id", handlers.GetWeatherData(s.db))
	weatherData.Post("/", handlers.CreateWeatherData(s.db))
	weatherData.Put("/", handlers.UpdateWeatherData(s.db))
	weatherData.Delete("/:id", handlers.DeleteWeatherData(s.db))
	weatherData.Get("/", handlers.GetAllWeatherData(s.db))

	// Incident routes
	incident := api.Group("/incident", roleMiddleware(types.Worker, types.Manager, types.Admin))

	incident.Get("/:id", handlers.GetIncident(s.db))
	incident.Post("/", handlers.CreateIncident(s.db, s.rmq))
	incident.Put("/", handlers.UpdateIncident(s.db))
	incident.Delete("/:id", handlers.DeleteIncident(s.db))
	incident.Get("/", handlers.GetAllIncidents(s.db))
	incident.Put("/:id/status", handlers.UpdateIncidentStatus(s.db))

	// Observation routes
	observation := api.Group("/observation", roleMiddleware(types.Worker, types.Manager, types.Admin))

	observation.Get("/:id", handlers.GetObservationLog(s.db))
	observation.Post("/", handlers.CreateObservationLog(s.db))
	observation.Put("/", handlers.UpdateObservationLog(s.db))
	observation.Delete("/:id", handlers.DeleteObservationLog(s.db))
	observation.Get("/", handlers.GetAllObservationLogs(s.db))

	// Maintenance routes
	maintenance := api.Group("/maintenance", roleMiddleware(types.Worker, types.Manager, types.Admin))

	maintenance.Get("/:id", handlers.GetMaintenancePlan(s.db))
	maintenance.Post("/", handlers.CreateMaintenancePlan(s.db))
	maintenance.Put("/", handlers.UpdateMaintenancePlan(s.db))
	maintenance.Delete("/:id", handlers.DeleteMaintenancePlan(s.db))
	maintenance.Get("/", handlers.GetAllMaintenancePlans(s.db))
	maintenance.Put("/:id/status", handlers.UpdateMaintenancePlanStatus(s.db))

}

// Run starts the server
func (s *Server) Run(address string) error {
	return s.app.Listen(address)
}

// Shutdown shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
