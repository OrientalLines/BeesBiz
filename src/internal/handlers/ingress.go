package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/orientallines/beesbiz/internal/database"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

// ObservationLog handlers
func GetObservationLog(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid observation log ID"})
		}
		log, err := db.GetObservationLog(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get observation log"})
		}
		return c.JSON(log)
	}
}

func CreateObservationLog(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var log types.ObservationLog
		if err := c.BodyParser(&log); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid observation log data"})
		}
		createdLog, err := db.CreateObservationLog(log)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create observation log"})
		}
		return c.JSON(createdLog)
	}
}

func UpdateObservationLog(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var log types.ObservationLog
		if err := c.BodyParser(&log); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid observation log data"})
		}
		updatedLog, err := db.UpdateObservationLog(log)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update observation log"})
		}
		return c.JSON(updatedLog)
	}
}

func DeleteObservationLog(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid observation log ID"})
		}
		if err := db.DeleteObservationLog(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete observation log"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllObservationLogs(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logs, err := db.GetAllObservationLogs()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get all observation logs"})
		}
		return c.JSON(logs)
	}
}

// MaintenancePlan handlers
func GetMaintenancePlan(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid maintenance plan ID"})
		}
		plan, err := db.GetMaintenancePlan(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get maintenance plan"})
		}
		return c.JSON(plan)
	}
}

func CreateMaintenancePlan(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var plan types.MaintenancePlan
		if err := c.BodyParser(&plan); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid maintenance plan data"})
		}
		createdPlan, err := db.CreateMaintenancePlan(plan)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create maintenance plan"})
		}
		return c.JSON(createdPlan)
	}
}

func UpdateMaintenancePlan(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var plan types.MaintenancePlan
		if err := c.BodyParser(&plan); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid maintenance plan data"})
		}
		updatedPlan, err := db.UpdateMaintenancePlan(plan)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update maintenance plan"})
		}
		return c.JSON(updatedPlan)
	}
}

func DeleteMaintenancePlan(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid maintenance plan ID"})
		}
		if err := db.DeleteMaintenancePlan(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete maintenance plan"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllMaintenancePlans(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		plans, err := db.GetAllMaintenancePlans()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get all maintenance plans"})
		}
		return c.JSON(plans)
	}
}

// Incident handlers
func GetIncident(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid incident ID"})
		}
		incident, err := db.GetIncident(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get incident"})
		}
		return c.JSON(incident)
	}
}

func CreateIncident(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var incident types.Incident
		if err := c.BodyParser(&incident); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid incident data"})
		}
		createdIncident, err := db.CreateIncident(incident)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create incident"})
		}
		return c.JSON(createdIncident)
	}
}

func UpdateIncident(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var incident types.Incident
		if err := c.BodyParser(&incident); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid incident data"})
		}
		updatedIncident, err := db.UpdateIncident(incident)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update incident"})
		}
		return c.JSON(updatedIncident)
	}
}

func DeleteIncident(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid incident ID"})
		}
		if err := db.DeleteIncident(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete incident"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllIncidents(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		incidents, err := db.GetAllIncidents()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get all incidents"})
		}
		return c.JSON(incidents)
	}
}
