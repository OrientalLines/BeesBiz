package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/orientallines/beesbiz/internal/database"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

// ObservationLog handlers
func GetObservationLog(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid observation log ID: %v", err)})
		}
		log, err := db.GetObservationLog(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get observation log: %v", err)})
		}
		return c.JSON(log)
	}
}

func CreateObservationLog(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var log types.ObservationLog
		if err := c.BodyParser(&log); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid observation log data: %v", err)})
		}
		createdLog, err := db.CreateObservationLog(log)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create observation log: %v", err)})
		}
		return c.JSON(createdLog)
	}
}

func UpdateObservationLog(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var log types.ObservationLog
		if err := c.BodyParser(&log); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid observation log data: %v", err)})
		}
		updatedLog, err := db.UpdateObservationLog(log)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update observation log: %v", err)})
		}
		return c.JSON(updatedLog)
	}
}

func DeleteObservationLog(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid observation log ID: %v", err)})
		}
		if err := db.DeleteObservationLog(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete observation log: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllObservationLogs(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logs, err := db.GetAllObservationLogs()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all observation logs: %v", err)})
		}
		return c.JSON(logs)
	}
}

// MaintenancePlan handlers
func GetMaintenancePlan(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid maintenance plan ID: %v", err)})
		}
		plan, err := db.GetMaintenancePlan(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get maintenance plan: %v", err)})
		}
		return c.JSON(plan)
	}
}

func CreateMaintenancePlan(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var plan types.MaintenancePlan
		if err := c.BodyParser(&plan); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid maintenance plan data: %v", err)})
		}
		createdPlan, err := db.CreateMaintenancePlan(plan)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create maintenance plan: %v", err)})
		}
		return c.JSON(createdPlan)
	}
}

func UpdateMaintenancePlan(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var plan types.MaintenancePlan
		if err := c.BodyParser(&plan); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid maintenance plan data: %v", err)})
		}
		updatedPlan, err := db.UpdateMaintenancePlan(plan)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update maintenance plan: %v", err)})
		}
		return c.JSON(updatedPlan)
	}
}

func DeleteMaintenancePlan(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid maintenance plan ID: %v", err)})
		}
		if err := db.DeleteMaintenancePlan(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete maintenance plan: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllMaintenancePlans(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		plans, err := db.GetAllMaintenancePlans()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all maintenance plans: %v", err)})
		}
		return c.JSON(plans)
	}
}

// Incident handlers
func GetIncident(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid incident ID: %v", err)})
		}
		incident, err := db.GetIncident(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get incident: %v", err)})
		}
		return c.JSON(incident)
	}
}

func CreateIncident(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var incident types.Incident
		if err := c.BodyParser(&incident); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid incident data: %v", err)})
		}
		createdIncident, err := db.CreateIncident(incident)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create incident: %v", err)})
		}
		return c.JSON(createdIncident)
	}
}

func UpdateIncident(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var incident types.Incident
		if err := c.BodyParser(&incident); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid incident data: %v", err)})
		}
		updatedIncident, err := db.UpdateIncident(incident)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update incident: %v", err)})
		}
		return c.JSON(updatedIncident)
	}
}

func DeleteIncident(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid incident ID: %v", err)})
		}
		if err := db.DeleteIncident(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete incident: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllIncidents(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		incidents, err := db.GetAllIncidents()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all incidents: %v", err)})
		}
		return c.JSON(incidents)
	}
}
