package types

import "github.com/guregu/null"

type Apiary struct {
	ApiaryID          int       `json:"apiary_id,omitempty" db:"apiary_id"`
	Location          string    `json:"location" db:"location"`
	ManagerID         int       `json:"manager_id" db:"manager_id"`
	EstablishmentDate null.Time `json:"establishment_date" db:"establishment_date"`
}

type Hive struct {
	HiveID           int       `json:"hive_id,omitempty" db:"hive_id"`
	ApiaryID         int       `json:"apiary_id" db:"apiary_id"`
	HiveType         string    `json:"hive_type" db:"hive_type"`
	InstallationDate null.Time `json:"installation_date" db:"installation_date"`
	CurrentStatus    string    `json:"current_status" db:"current_status"`
}

type BeeCommunity struct {
	CommunityID        int    `json:"community_id,omitempty" db:"community_id"`
	HiveID             int    `json:"hive_id" db:"hive_id"`
	QueenAge           int    `json:"queen_age" db:"queen_age"`
	PopulationEstimate int    `json:"population_estimate" db:"population_estimate"`
	HealthStatus       string `json:"health_status" db:"health_status"`
}

type HoneyHarvest struct {
	HarvestID    int       `json:"harvest_id,omitempty" db:"harvest_id"`
	HiveID       int       `json:"hive_id" db:"hive_id"`
	HarvestDate  null.Time `json:"harvest_date" db:"harvest_date"`
	Quantity     float32   `json:"quantity" db:"quantity"`
	QualityGrade string    `json:"quality_grade" db:"quality_grade"`
}
