package rest

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/orientallines/beesbiz/internal/config"
	"github.com/orientallines/beesbiz/internal/database"
	"github.com/orientallines/beesbiz/internal/handlers"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

// Server is a wrapper around fiber.App
type Server struct {
	app    *fiber.App
	db     *database.DB
	jwtKey []byte
}

// NewServer creates a new Server
func NewServer(db *database.DB) *Server {
	return &Server{
		app:    fiber.New(),
		db:     db,
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
	hive.Delete("/:id", handlers.DeleteHive(s.db))
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
	region := api.Group("/region", roleMiddleware(types.Manager, types.Admin))

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

	// ProductionReport routes
	productionReport := api.Group("/production-report", roleMiddleware(types.Manager, types.Worker, types.Admin))

	productionReport.Get("/:id", handlers.GetProductionReport(s.db))
	productionReport.Post("/", handlers.CreateProductionReport(s.db))
	productionReport.Put("/", handlers.UpdateProductionReport(s.db))
	productionReport.Delete("/:id", handlers.DeleteProductionReport(s.db))
	productionReport.Get("/", handlers.GetAllProductionReports(s.db))

	// Sensor routes
	sensor := api.Group("/sensor", roleMiddleware(types.Admin, types.Manager, types.Worker))

	sensor.Get("/:id", handlers.GetSensor(s.db))
	sensor.Post("/", handlers.CreateSensor(s.db))
	sensor.Put("/", handlers.UpdateSensor(s.db))
	sensor.Delete("/:id", handlers.DeleteSensor(s.db))
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

}

// Run starts the server
func (s *Server) Run(address string) error {
	return s.app.Listen(address)
}

// Shutdown shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
