package database

import (
	"fmt"

	types "github.com/orientallines/beesbiz/internal/types/db"
)

func (db *DB) GetRegion(id int) (types.Region, error) {
	var region types.Region
	err := db.Get(&region, "SELECT * FROM region WHERE region_id = $1", id)
	if err != nil {
		return types.Region{}, fmt.Errorf("error getting region: %w", err)
	}
	return region, nil
}

func (db *DB) CreateRegion(region types.Region) (types.Region, error) {
	var createdRegion types.Region
	err := db.Get(&createdRegion, "INSERT INTO region (name, climate_zone) VALUES ($1, $2) RETURNING *", region.Name, region.ClimateZone)
	if err != nil {
		return types.Region{}, fmt.Errorf("error creating region: %w", err)
	}
	return createdRegion, nil
}

func (db *DB) UpdateRegion(region types.Region) (types.Region, error) {
	var updatedRegion types.Region
	err := db.Get(&updatedRegion, "UPDATE region SET name = $1, climate_zone = $2 WHERE region_id = $3 RETURNING *", region.Name, region.ClimateZone, region.RegionID)
	if err != nil {
		return types.Region{}, fmt.Errorf("error updating region: %w", err)
	}
	return updatedRegion, nil
}

func (db *DB) DeleteRegion(id int) error {
	_, err := db.Exec("DELETE FROM region WHERE region_id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting region: %w", err)
	}
	return nil
}

func (db *DB) GetAllRegions() ([]types.Region, error) {
	var regions []types.Region
	err := db.Select(&regions, "SELECT * FROM region")
	if err != nil {
		return []types.Region{}, fmt.Errorf("error getting all regions: %w", err)
	}
	return regions, nil
}

func (db *DB) CreateAllowedRegion(allowedRegion types.AllowedRegion) (types.AllowedRegion, error) {
	var createdAllowedRegion types.AllowedRegion
	err := db.Get(&createdAllowedRegion, "INSERT INTO allowed_region (user_id, region_id) VALUES ($1, $2) RETURNING *", allowedRegion.UserID, allowedRegion.RegionID)
	if err != nil {
		return types.AllowedRegion{}, fmt.Errorf("error creating allowed region: %w", err)
	}
	return createdAllowedRegion, nil
}

func (db *DB) UpdateAllowedRegion(allowedRegion types.AllowedRegion) (types.AllowedRegion, error) {
	var updatedAllowedRegion types.AllowedRegion
	err := db.Get(&updatedAllowedRegion, "UPDATE allowed_region SET user_id = $1, region_id = $2 WHERE id = $3 RETURNING *", allowedRegion.UserID, allowedRegion.RegionID, allowedRegion.ID)
	if err != nil {
		return types.AllowedRegion{}, fmt.Errorf("error updating allowed region: %w", err)
	}
	return updatedAllowedRegion, nil
}

func (db *DB) DeleteAllowedRegion(id int) error {
	_, err := db.Exec("DELETE FROM allowed_region WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting allowed region: %w", err)
	}
	return nil
}

func (db *DB) GetAllAllowedRegions() ([]types.AllowedRegion, error) {
	var allowedRegions []types.AllowedRegion
	err := db.Select(&allowedRegions, "SELECT * FROM allowed_region")
	if err != nil {
		return []types.AllowedRegion{}, fmt.Errorf("error getting all allowed regions: %w", err)
	}
	return allowedRegions, nil
}

func (db *DB) GetRegionApiary(id int) (types.RegionApiary, error) {
	var regionApiary types.RegionApiary
	err := db.Get(&regionApiary, "SELECT * FROM region_apiary WHERE id = $1", id)
	if err != nil {
		return types.RegionApiary{}, fmt.Errorf("error getting region apiary: %w", err)
	}
	return regionApiary, nil
}

func (db *DB) CreateRegionApiary(regionApiary types.RegionApiary) (types.RegionApiary, error) {
	var createdRegionApiary types.RegionApiary
	err := db.Get(&createdRegionApiary, "INSERT INTO region_apiary (apiary_id, region_id) VALUES ($1, $2) RETURNING *", regionApiary.ApiaryID, regionApiary.RegionID)
	if err != nil {
		return types.RegionApiary{}, fmt.Errorf("error creating region apiary: %w", err)
	}
	return createdRegionApiary, nil
}

func (db *DB) UpdateRegionApiary(regionApiary types.RegionApiary) (types.RegionApiary, error) {
	var updatedRegionApiary types.RegionApiary
	err := db.Get(&updatedRegionApiary, "UPDATE region_apiary SET apiary_id = $1, region_id = $2 WHERE id = $3 RETURNING *", regionApiary.ApiaryID, regionApiary.RegionID, regionApiary.ID)
	if err != nil {
		return types.RegionApiary{}, fmt.Errorf("error updating region apiary: %w", err)
	}
	return updatedRegionApiary, nil
}

func (db *DB) DeleteRegionApiary(id int) error {
	_, err := db.Exec("DELETE FROM region_apiary WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting region apiary: %w", err)
	}
	return nil
}

func (db *DB) GetAllRegionApiaries() ([]types.RegionApiary, error) {
	var regionApiaries []types.RegionApiary
	err := db.Select(&regionApiaries, "SELECT * FROM region_apiary")
	if err != nil {
		return []types.RegionApiary{}, fmt.Errorf("error getting all region apiaries: %w", err)
	}
	return regionApiaries, nil
}
