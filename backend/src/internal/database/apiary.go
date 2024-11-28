package database

import (
	"fmt"

	types "github.com/orientallines/beesbiz/internal/types/db"
	"go.uber.org/zap"
)

func (db *DB) GetApiary(id int) (types.Apiary, error) {
	var apiary types.Apiary
	err := db.Get(&apiary, "SELECT * FROM apiary WHERE apiary_id = $1", id)
	if err != nil {
		zap.S().Error("Error getting apiary: ", err)
		return types.Apiary{}, fmt.Errorf("error getting apiary: %w", err)
	}
	return apiary, nil
}

func (db *DB) CreateApiary(apiary types.Apiary) (types.Apiary, error) {
	var createdApiary types.Apiary
	err := db.Get(&createdApiary, "INSERT INTO apiary (location, manager_id, establishment_date) VALUES ($1, $2, $3) RETURNING *", apiary.Location, apiary.ManagerID, apiary.EstablishmentDate)
	if err != nil {
		zap.S().Error("Error creating apiary: ", err)
		return types.Apiary{}, fmt.Errorf("error creating apiary: %w", err)
	}
	return createdApiary, nil
}

func (db *DB) UpdateApiary(apiary types.Apiary) (types.Apiary, error) {
	var updatedApiary types.Apiary
	err := db.Get(&updatedApiary, "UPDATE apiary SET location = $1, manager_id = $2, establishment_date = $3 WHERE id = $4 RETURNING *", apiary.Location, apiary.ManagerID, apiary.EstablishmentDate, apiary.ApiaryID)
	if err != nil {
		zap.S().Error("Error updating apiary: ", err)
		return types.Apiary{}, fmt.Errorf("error updating apiary: %w", err)
	}
	return updatedApiary, nil
}

func (db *DB) DeleteApiary(id int) error {
	_, err := db.Exec("DELETE FROM apiary WHERE apiary_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting apiary: ", err)
		return fmt.Errorf("error deleting apiary: %w", err)
	}
	return nil
}

func (db *DB) GetAllApiaries() ([]types.Apiary, error) {
	var apiaries []types.Apiary
	err := db.Select(&apiaries, "SELECT * FROM apiary")
	if err != nil {
		zap.S().Error("Error getting all apiaries: ", err)
		return []types.Apiary{}, fmt.Errorf("error getting all apiaries: %w", err)
	}
	return apiaries, nil
}

func (db *DB) GetAllHives() ([]types.Hive, error) {
	var hives []types.Hive
	err := db.Select(&hives, "SELECT * FROM hive")
	if err != nil {
		zap.S().Error("Error getting all hives: ", err)
		return []types.Hive{}, fmt.Errorf("error getting all hives: %w", err)
	}
	return hives, nil
}

func (db *DB) CreateHive(hive types.Hive) (types.Hive, error) {
	var createdHive types.Hive
	err := db.Get(&createdHive, "INSERT INTO hive (apiary_id, hive_type, installation_date, current_status) VALUES ($1, $2, $3, $4) RETURNING *", hive.ApiaryID, hive.HiveType, hive.InstallationDate, hive.CurrentStatus)
	if err != nil {
		zap.S().Error("Error creating hive: ", err)
		return types.Hive{}, fmt.Errorf("error creating hive: %w", err)
	}
	return createdHive, nil
}

func (db *DB) UpdateHive(hive types.Hive) (types.Hive, error) {
	var updatedHive types.Hive
	err := db.Get(&updatedHive, "UPDATE hive SET apiary_id = $1, hive_type = $2, installation_date = $3, current_status = $4 WHERE hive_id = $5 RETURNING *", hive.ApiaryID, hive.HiveType, hive.InstallationDate, hive.CurrentStatus, hive.HiveID)
	if err != nil {
		zap.S().Error("Error updating hive: ", err)
		return types.Hive{}, fmt.Errorf("error updating hive: %w", err)
	}
	return updatedHive, nil
}

func (db *DB) DeleteHive(id int) error {
	_, err := db.Exec("DELETE FROM hive WHERE hive_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting hive: ", err)
		return fmt.Errorf("error deleting hive: %w", err)
	}
	return nil
}

func (db *DB) GetAllHivesByApiaryID(apiaryID int) ([]types.Hive, error) {
	var hives []types.Hive
	err := db.Select(&hives, "SELECT * FROM hive WHERE apiary_id = $1", apiaryID)
	if err != nil {
		zap.S().Error("Error getting all hives by apiary id: ", err)
		return []types.Hive{}, fmt.Errorf("error getting all hives by apiary id: %w", err)
	}
	return hives, nil
}

func (db *DB) GetAllBeeCommunities() ([]types.BeeCommunity, error) {
	var beeCommunities []types.BeeCommunity
	err := db.Select(&beeCommunities, "SELECT * FROM bee_community")
	if err != nil {
		zap.S().Error("Error getting all bee communities: ", err)
		return []types.BeeCommunity{}, fmt.Errorf("error getting all bee communities: %w", err)
	}
	return beeCommunities, nil
}

func (db *DB) CreateBeeCommunity(beeCommunity types.BeeCommunity) (types.BeeCommunity, error) {
	var createdBeeCommunity types.BeeCommunity
	err := db.Get(&createdBeeCommunity, "INSERT INTO bee_community (hive_id, queen_age, population_estimate, health_status) VALUES ($1, $2, $3, $4) RETURNING *", beeCommunity.HiveID, beeCommunity.QueenAge, beeCommunity.PopulationEstimate, beeCommunity.HealthStatus)
	if err != nil {
		zap.S().Error("Error creating bee community: ", err)
		return types.BeeCommunity{}, fmt.Errorf("error creating bee community: %w", err)
	}
	return createdBeeCommunity, nil
}

func (db *DB) UpdateBeeCommunity(beeCommunity types.BeeCommunity) (types.BeeCommunity, error) {
	var updatedBeeCommunity types.BeeCommunity
	err := db.Get(&updatedBeeCommunity, "UPDATE bee_community SET hive_id = $1, queen_age = $2, population_estimate = $3, health_status = $4 WHERE community_id = $5 RETURNING *", beeCommunity.HiveID, beeCommunity.QueenAge, beeCommunity.PopulationEstimate, beeCommunity.HealthStatus, beeCommunity.CommunityID)
	if err != nil {
		zap.S().Error("Error updating bee community: ", err)
		return types.BeeCommunity{}, fmt.Errorf("error updating bee community: %w", err)
	}
	return updatedBeeCommunity, nil
}

func (db *DB) DeleteBeeCommunity(id int) error {
	_, err := db.Exec("DELETE FROM bee_community WHERE community_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting bee community: ", err)
		return fmt.Errorf("error deleting bee community: %w", err)
	}
	return nil
}

func (db *DB) GetAllBeeCommunitiesByHiveID(hiveID int) ([]types.BeeCommunity, error) {
	var beeCommunities []types.BeeCommunity
	err := db.Select(&beeCommunities, "SELECT * FROM bee_community WHERE hive_id = $1", hiveID)
	if err != nil {
		zap.S().Error("Error getting all bee communities by hive id: ", err)
		return []types.BeeCommunity{}, fmt.Errorf("error getting all bee communities by hive id: %w", err)
	}
	return beeCommunities, nil
}

func (db *DB) GetHoneyHarvest(id int) (types.HoneyHarvest, error) {
	var harvest types.HoneyHarvest
	err := db.Get(&harvest, "SELECT * FROM honey_harvest WHERE harvest_id = $1", id)
	if err != nil {
		zap.S().Error("Error getting honey harvest: ", err)
		return types.HoneyHarvest{}, fmt.Errorf("error getting honey harvest: %w", err)
	}
	return harvest, nil
}

func (db *DB) CreateHoneyHarvest(harvest types.HoneyHarvest) (types.HoneyHarvest, error) {
	var createdHarvest types.HoneyHarvest
	err := db.Get(&createdHarvest, "INSERT INTO honey_harvest (hive_id, harvest_date, quantity, quality_grade) VALUES ($1, $2, $3, $4) RETURNING *", harvest.HiveID, harvest.HarvestDate, harvest.Quantity, harvest.QualityGrade)
	if err != nil {
		zap.S().Error("Error creating honey harvest: ", err)
		return types.HoneyHarvest{}, fmt.Errorf("error creating honey harvest: %w", err)
	}
	return createdHarvest, nil
}

func (db *DB) UpdateHoneyHarvest(harvest types.HoneyHarvest) (types.HoneyHarvest, error) {
	var updatedHarvest types.HoneyHarvest
	err := db.Get(&updatedHarvest, "UPDATE honey_harvest SET hive_id = $1, harvest_date = $2, quantity = $3, quality_grade = $4 WHERE harvest_id = $5 RETURNING *", harvest.HiveID, harvest.HarvestDate, harvest.Quantity, harvest.QualityGrade, harvest.HarvestID)
	if err != nil {
		zap.S().Error("Error updating honey harvest: ", err)
		return types.HoneyHarvest{}, fmt.Errorf("error updating honey harvest: %w", err)
	}
	return updatedHarvest, nil
}

func (db *DB) DeleteHoneyHarvest(id int) error {
	_, err := db.Exec("DELETE FROM honey_harvest WHERE harvest_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting honey harvest: ", err)
		return fmt.Errorf("error deleting honey harvest: %w", err)
	}
	return nil
}

func (db *DB) GetAllHoneyHarvests() ([]types.HoneyHarvest, error) {
	var harvests []types.HoneyHarvest
	err := db.Select(&harvests, "SELECT * FROM honey_harvest")
	if err != nil {
		zap.S().Error("Error getting all honey harvests: ", err)
		return []types.HoneyHarvest{}, fmt.Errorf("error getting all honey harvests: %w", err)
	}
	return harvests, nil
}

func (db *DB) GetAllSensorsByHiveID(hiveID int) ([]types.Sensor, error) {
	var sensors []types.Sensor
	err := db.Select(&sensors, "SELECT * FROM sensor WHERE hive_id = $1", hiveID)
	if err != nil {
		zap.S().Error("Error getting sensors by hive id: ", err)
		return nil, fmt.Errorf("error getting sensors by hive id: %w", err)
	}
	return sensors, nil
}
