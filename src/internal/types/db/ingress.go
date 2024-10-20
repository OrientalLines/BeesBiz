package types

import "time"

type ObservationLog struct {
	LogID           int       `json:"log_id" db:"log_id"`
	HiveID          int       `json:"hive_id" db:"hive_id"`
	ObservationDate time.Time `json:"observation_date" db:"observation_date"`
	Description     string    `json:"description" db:"description"`
	Recommendations string    `json:"recommendations" db:"recommendations"`
}

type MaintenancePlan struct {
	PlanID      int       `json:"plan_id" db:"plan_id"`
	ApiaryID    int       `json:"apiary_id" db:"apiary_id"`
	PlannedDate time.Time `json:"planned_date" db:"planned_date"`
	WorkType    string    `json:"work_type" db:"work_type"`
	AssignedTo  int       `json:"assigned_to" db:"assigned_to"`
}

type Incident struct {
	IncidentID   int       `json:"incident_id" db:"incident_id"`
	HiveID       int       `json:"hive_id" db:"hive_id"`
	IncidentDate time.Time `json:"incident_date" db:"incident_date"`
	Description  string    `json:"description" db:"description"`
	Severity     string    `json:"severity" db:"severity"`
	ActionsTaken string    `json:"actions_taken" db:"actions_taken"`
}
