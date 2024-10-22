package types

type Region struct {
	RegionID    int    `json:"region_id,omitempty" db:"region_id"`
	Name        string `json:"name" db:"name"`
	ClimateZone string `json:"climate_zone" db:"climate_zone"`
}

type RegionApiary struct {
	ID       int    `json:"id,omitempty" db:"id"`
	ApiaryID int    `json:"apiary_id" db:"apiary_id"`
	RegionID int    `json:"region_id" db:"region_id"`
	Region   Region `json:"region" db:"region"`
	Apiary   Apiary `json:"apiary" db:"apiary"`
}

type AllowedRegion struct {
	ID       int    `json:"id,omitempty" db:"id"`
	UserID   int    `json:"user_id" db:"user_id"`
	RegionID int    `json:"region_id" db:"region_id"`
	Region   Region `json:"region" db:"region"`
	User     User   `json:"user" db:"user"`
}
