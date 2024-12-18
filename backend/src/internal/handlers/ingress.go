package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/orientallines/beesbiz/internal/database"
	"github.com/orientallines/beesbiz/internal/rabbitmq"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

const IncidentQueue = "incident_queue"

// ObservationLog handlers

// GetObservationLog gets an observation log by ID
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

// CreateObservationLog creates a new observation log
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

// UpdateObservationLog updates an observation log
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

// DeleteObservationLog deletes an observation log
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

// GetAllObservationLogs gets all observation logs
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

// GetMaintenancePlan gets a maintenance plan by ID
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

// CreateMaintenancePlan creates a new maintenance plan
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

// UpdateMaintenancePlan updates a maintenance plan
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

// DeleteMaintenancePlan deletes a maintenance plan
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

// GetAllMaintenancePlans gets all maintenance plans
func GetAllMaintenancePlans(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		plans, err := db.GetAllMaintenancePlans()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all maintenance plans: %v", err)})
		}
		return c.JSON(plans)
	}
}

// Add this new handler after the other maintenance plan handlers
// UpdateMaintenancePlanStatus updates a maintenance plan's status
func UpdateMaintenancePlanStatus(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get plan ID from params
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid maintenance plan ID: %v", err),
			})
		}

		// Parse request body
		var updateData struct {
			Status string `json:"status"`
		}
		if err := c.BodyParser(&updateData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid request data: %v", err),
			})
		}

		// Get existing plan
		plan, err := db.GetMaintenancePlan(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to get maintenance plan: %v", err),
			})
		}

		// Update status
		plan.Status = updateData.Status

		// Save updated plan
		updatedPlan, err := db.UpdateMaintenancePlan(plan)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to update maintenance plan status: %v", err),
			})
		}

		return c.JSON(updatedPlan)
	}
}

// Incident handlers

// GetIncident gets an incident by ID
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

// CreateIncident creates a new incident
func CreateIncident(db *database.DB, rmq *rabbitmq.RabbitMQ) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var incident types.Incident
		if err := c.BodyParser(&incident); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid incident data: %v", err)})
		}
		createdIncident, err := db.CreateIncident(incident)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create incident: %v", err)})
		}

		// Publish RabbitMQ notification
		if err := rmq.PublishMessage("incident_queue", createdIncident); err != nil {
			// Log the error but don't fail the request
			zap.L().Error("Failed to publish incident notification", zap.Error(err))
		}

		return c.JSON(createdIncident)
	}
}

// UpdateIncident updates an incident
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

// DeleteIncident deletes an incident
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

// GetAllIncidents gets all incidents
func GetAllIncidents(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		incidents, err := db.GetAllIncidents()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all incidents: %v", err)})
		}
		return c.JSON(incidents)
	}
}

// Add this new handler
// UpdateIncidentStatus updates an incident's status
func UpdateIncidentStatus(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get incident ID from params
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid incident ID: %v", err),
			})
		}

		// Parse request body
		var updateData struct {
			Severity string `json:"severity"`
		}
		if err := c.BodyParser(&updateData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid request data: %v", err),
			})
		}

		// Get existing incident
		incident, err := db.GetIncident(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to get incident: %v", err),
			})
		}

		// Update severity
		incident.Severity = updateData.Severity

		// Save updated incident
		updatedIncident, err := db.UpdateIncident(incident)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to update incident status: %v", err),
			})
		}

		return c.JSON(updatedIncident)
	}
}
