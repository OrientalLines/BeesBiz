package database

import (
	"fmt"

	types "github.com/orientallines/beesbiz/internal/types/db"
	"go.uber.org/zap"
)

func (db *DB) GetVeterinaryPassport(id int) (types.VeterinaryPassport, error) {
	var passport types.VeterinaryPassport
	err := db.Get(&passport, "SELECT * FROM veterinary_passport WHERE passport_id = $1", id)
	if err != nil {
		zap.S().Error("Error getting veterinary passport: ", err)
		return types.VeterinaryPassport{}, fmt.Errorf("error getting veterinary passport: %w", err)
	}
	return passport, nil
}

func (db *DB) CreateVeterinaryPassport(passport types.VeterinaryPassport) (types.VeterinaryPassport, error) {
	var createdPassport types.VeterinaryPassport
	err := db.Get(&createdPassport, "INSERT INTO veterinary_passport (bee_community_id, issue_date, health_status, last_inspection_date) VALUES ($1, $2, $3, $4) RETURNING *", passport.BeeCommunityID, passport.IssueDate, passport.HealthStatus, passport.LastInspectionDate)
	if err != nil {
		zap.S().Error("Error creating veterinary passport: ", err)
		return types.VeterinaryPassport{}, fmt.Errorf("error creating veterinary passport: %w", err)
	}
	return createdPassport, nil
}

func (db *DB) UpdateVeterinaryPassport(passport types.VeterinaryPassport) (types.VeterinaryPassport, error) {
	var updatedPassport types.VeterinaryPassport
	err := db.Get(&updatedPassport, "UPDATE veterinary_passport SET bee_community_id = $1, issue_date = $2, health_status = $3, last_inspection_date = $4 WHERE passport_id = $5 RETURNING *", passport.BeeCommunityID, passport.IssueDate, passport.HealthStatus, passport.LastInspectionDate, passport.PassportID)
	if err != nil {
		zap.S().Error("Error updating veterinary passport: ", err)
		return types.VeterinaryPassport{}, fmt.Errorf("error updating veterinary passport: %w", err)
	}
	return updatedPassport, nil
}

func (db *DB) DeleteVeterinaryPassport(id int) error {
	_, err := db.Exec("DELETE FROM veterinary_passport WHERE passport_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting veterinary passport: ", err)
		return fmt.Errorf("error deleting veterinary passport: %w", err)
	}
	return nil
}

func (db *DB) GetAllVeterinaryPassports() ([]types.VeterinaryPassport, error) {
	var passports []types.VeterinaryPassport
	err := db.Select(&passports, "SELECT * FROM veterinary_passport")
	if err != nil {
		zap.S().Error("Error getting all veterinary passports: ", err)
		return []types.VeterinaryPassport{}, fmt.Errorf("error getting all veterinary passports: %w", err)
	}
	return passports, nil
}

func (db *DB) GetVeterinaryRecord(id int) (types.VeterinaryRecord, error) {
	var record types.VeterinaryRecord
	err := db.Get(&record, "SELECT * FROM veterinary_record WHERE record_id = $1", id)
	if err != nil {
		zap.S().Error("Error getting veterinary record: ", err)
		return types.VeterinaryRecord{}, fmt.Errorf("error getting veterinary record: %w", err)
	}
	return record, nil
}

func (db *DB) CreateVeterinaryRecord(record types.VeterinaryRecord) (types.VeterinaryRecord, error) {
	var createdRecord types.VeterinaryRecord
	err := db.Get(&createdRecord, "INSERT INTO veterinary_record (passport_id, record_date, description, treatment) VALUES ($1, $2, $3, $4) RETURNING *", record.PassportID, record.RecordDate, record.Description, record.Treatment)
	if err != nil {
		zap.S().Error("Error creating veterinary record: ", err)
		return types.VeterinaryRecord{}, fmt.Errorf("error creating veterinary record: %w", err)
	}
	return createdRecord, nil
}

func (db *DB) UpdateVeterinaryRecord(record types.VeterinaryRecord) (types.VeterinaryRecord, error) {
	var updatedRecord types.VeterinaryRecord
	err := db.Get(&updatedRecord, "UPDATE veterinary_record SET passport_id = $1, record_date = $2, description = $3, treatment = $4 WHERE record_id = $5 RETURNING *", record.PassportID, record.RecordDate, record.Description, record.Treatment, record.RecordID)
	if err != nil {
		zap.S().Error("Error updating veterinary record: ", err)
		return types.VeterinaryRecord{}, fmt.Errorf("error updating veterinary record: %w", err)
	}
	return updatedRecord, nil
}

func (db *DB) DeleteVeterinaryRecord(id int) error {
	_, err := db.Exec("DELETE FROM veterinary_record WHERE record_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting veterinary record: ", err)
		return fmt.Errorf("error deleting veterinary record: %w", err)
	}
	return nil
}

func (db *DB) GetAllVeterinaryRecords() ([]types.VeterinaryRecord, error) {
	var records []types.VeterinaryRecord
	err := db.Select(&records, "SELECT * FROM veterinary_record")
	if err != nil {
		zap.S().Error("Error getting all veterinary records: ", err)
		return []types.VeterinaryRecord{}, fmt.Errorf("error getting all veterinary records: %w", err)
	}
	return records, nil
}
