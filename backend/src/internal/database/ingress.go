package database

import (
	"fmt"

	types "github.com/orientallines/beesbiz/internal/types/db"
	"go.uber.org/zap"
)

func (db *DB) GetObservationLog(id int) (types.ObservationLog, error) {
	var log types.ObservationLog
	err := db.Get(&log, "SELECT * FROM observation_log WHERE log_id = $1", id)
	if err != nil {
		zap.S().Error("Error getting observation log: ", err)
		return types.ObservationLog{}, fmt.Errorf("error getting observation log: %w", err)
	}
	return log, nil
}

func (db *DB) CreateObservationLog(log types.ObservationLog) (types.ObservationLog, error) {
	var createdLog types.ObservationLog
	err := db.Get(&createdLog, "INSERT INTO observation_log (hive_id, observation_date, description, recommendations) VALUES ($1, $2, $3, $4) RETURNING *", log.HiveID, log.ObservationDate, log.Description, log.Recommendations)
	if err != nil {
		zap.S().Error("Error creating observation log: ", err)
		return types.ObservationLog{}, fmt.Errorf("error creating observation log: %w", err)
	}
	return createdLog, nil
}

func (db *DB) UpdateObservationLog(log types.ObservationLog) (types.ObservationLog, error) {
	var updatedLog types.ObservationLog
	err := db.Get(&updatedLog, "UPDATE observation_log SET hive_id = $1, observation_date = $2, description = $3, recommendations = $4 WHERE log_id = $5 RETURNING *", log.HiveID, log.ObservationDate, log.Description, log.Recommendations, log.LogID)
	if err != nil {
		zap.S().Error("Error updating observation log: ", err)
		return types.ObservationLog{}, fmt.Errorf("error updating observation log: %w", err)
	}
	return updatedLog, nil
}

func (db *DB) DeleteObservationLog(id int) error {
	_, err := db.Exec("DELETE FROM observation_log WHERE log_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting observation log: ", err)
		return fmt.Errorf("error deleting observation log: %w", err)
	}
	return nil
}

func (db *DB) GetAllObservationLogs() ([]types.ObservationLog, error) {
	var logs []types.ObservationLog
	err := db.Select(&logs, "SELECT * FROM observation_log")
	if err != nil {
		zap.S().Error("Error getting all observation logs: ", err)
		return []types.ObservationLog{}, fmt.Errorf("error getting all observation logs: %w", err)
	}
	return logs, nil
}

func (db *DB) GetMaintenancePlan(id int) (types.MaintenancePlan, error) {
	var plan types.MaintenancePlan
	err := db.Get(&plan, "SELECT * FROM maintenance_plan WHERE plan_id = $1", id)
	if err != nil {
		zap.S().Error("Error getting maintenance plan: ", err)
		return types.MaintenancePlan{}, fmt.Errorf("error getting maintenance plan: %w", err)
	}
	return plan, nil
}

func (db *DB) CreateMaintenancePlan(plan types.MaintenancePlan) (types.MaintenancePlan, error) {
	var createdPlan types.MaintenancePlan
	err := db.Get(&createdPlan, "INSERT INTO maintenance_plan (apiary_id, planned_date, work_type, assigned_to) VALUES ($1, $2, $3, $4) RETURNING *", plan.ApiaryID, plan.PlannedDate, plan.WorkType, plan.AssignedTo)
	if err != nil {
		zap.S().Error("Error creating maintenance plan: ", err)
		return types.MaintenancePlan{}, fmt.Errorf("error creating maintenance plan: %w", err)
	}
	return createdPlan, nil
}

func (db *DB) UpdateMaintenancePlan(plan types.MaintenancePlan) (types.MaintenancePlan, error) {
	var updatedPlan types.MaintenancePlan
	err := db.Get(&updatedPlan, "UPDATE maintenance_plan SET apiary_id = $1, planned_date = $2, work_type = $3, assigned_to = $4 WHERE plan_id = $5 RETURNING *", plan.ApiaryID, plan.PlannedDate, plan.WorkType, plan.AssignedTo, plan.PlanID)
	if err != nil {
		zap.S().Error("Error updating maintenance plan: ", err)
		return types.MaintenancePlan{}, fmt.Errorf("error updating maintenance plan: %w", err)
	}
	return updatedPlan, nil
}

func (db *DB) DeleteMaintenancePlan(id int) error {
	_, err := db.Exec("DELETE FROM maintenance_plan WHERE plan_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting maintenance plan: ", err)
		return fmt.Errorf("error deleting maintenance plan: %w", err)
	}
	return nil
}

func (db *DB) GetAllMaintenancePlans() ([]types.MaintenancePlan, error) {
	var plans []types.MaintenancePlan
	err := db.Select(&plans, "SELECT * FROM maintenance_plan")
	if err != nil {
		zap.S().Error("Error getting all maintenance plans: ", err)
		return []types.MaintenancePlan{}, fmt.Errorf("error getting all maintenance plans: %w", err)
	}
	return plans, nil
}

func (db *DB) GetIncident(id int) (types.Incident, error) {
	var incident types.Incident
	err := db.Get(&incident, "SELECT * FROM incident WHERE incident_id = $1", id)
	if err != nil {
		zap.S().Error("Error getting incident: ", err)
		return types.Incident{}, fmt.Errorf("error getting incident: %w", err)
	}
	return incident, nil
}

func (db *DB) CreateIncident(incident types.Incident) (types.Incident, error) {
	var createdIncident types.Incident
	err := db.Get(&createdIncident, "INSERT INTO incident (hive_id, incident_date, description, severity, actions_taken) VALUES ($1, $2, $3, $4, $5) RETURNING *", incident.HiveID, incident.IncidentDate, incident.Description, incident.Severity, incident.ActionsTaken)
	if err != nil {
		zap.S().Error("Error creating incident: ", err)
		return types.Incident{}, fmt.Errorf("error creating incident: %w", err)
	}
	return createdIncident, nil
}

func (db *DB) UpdateIncident(incident types.Incident) (types.Incident, error) {
	var updatedIncident types.Incident
	err := db.Get(&updatedIncident, "UPDATE incident SET hive_id = $1, incident_date = $2, description = $3, severity = $4, actions_taken = $5 WHERE incident_id = $6 RETURNING *", incident.HiveID, incident.IncidentDate, incident.Description, incident.Severity, incident.ActionsTaken, incident.IncidentID)
	if err != nil {
		zap.S().Error("Error updating incident: ", err)
		return types.Incident{}, fmt.Errorf("error updating incident: %w", err)
	}
	return updatedIncident, nil
}

func (db *DB) DeleteIncident(id int) error {
	_, err := db.Exec("DELETE FROM incident WHERE incident_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting incident: ", err)
		return fmt.Errorf("error deleting incident: %w", err)
	}
	return nil
}

func (db *DB) GetAllIncidents() ([]types.Incident, error) {
	var incidents []types.Incident
	err := db.Select(&incidents, "SELECT * FROM incident")
	if err != nil {
		zap.S().Error("Error getting all incidents: ", err)
		return []types.Incident{}, fmt.Errorf("error getting all incidents: %w", err)
	}
	return incidents, nil
}
