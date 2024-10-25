package database

import (
	"fmt"

	types "github.com/orientallines/beesbiz/internal/types/db"
	"go.uber.org/zap"
)

func (db *DB) GetProductionReport(id int) (types.ProductionReport, error) {
	var report types.ProductionReport
	err := db.Get(&report, "SELECT * FROM production_report WHERE report_id = $1", id)
	if err != nil {
		zap.S().Error("Error getting production report: ", err)
		return types.ProductionReport{}, fmt.Errorf("error getting production report: %w", err)
	}
	return report, nil
}

func (db *DB) CreateProductionReport(report types.ProductionReport) (types.ProductionReport, error) {
	var createdReport types.ProductionReport
	err := db.Get(&createdReport, "INSERT INTO production_report (apiary_id, start_date, end_date, total_honey_produced, total_expenses) VALUES ($1, $2, $3, $4, $5) RETURNING *", report.ApiaryID, report.StartDate, report.EndDate, report.TotalHoney, report.TotalExpenses)
	if err != nil {
		zap.S().Error("Error creating production report: ", err)
		return types.ProductionReport{}, fmt.Errorf("error creating production report: %w", err)
	}
	return createdReport, nil
}

func (db *DB) UpdateProductionReport(report types.ProductionReport) (types.ProductionReport, error) {
	var updatedReport types.ProductionReport
	err := db.Get(&updatedReport, "UPDATE production_report SET apiary_id = $1, start_date = $2, end_date = $3, total_honey_produced = $4, total_expenses = $5 WHERE report_id = $6 RETURNING *", report.ApiaryID, report.StartDate, report.EndDate, report.TotalHoney, report.TotalExpenses, report.ReportID)
	if err != nil {
		zap.S().Error("Error updating production report: ", err)
		return types.ProductionReport{}, fmt.Errorf("error updating production report: %w", err)
	}
	return updatedReport, nil
}

func (db *DB) DeleteProductionReport(id int) error {
	_, err := db.Exec("DELETE FROM production_report WHERE report_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting production report: ", err)
		return fmt.Errorf("error deleting production report: %w", err)
	}
	return nil
}

func (db *DB) GetAllProductionReports() ([]types.ProductionReport, error) {
	var reports []types.ProductionReport
	err := db.Select(&reports, "SELECT * FROM production_report")
	if err != nil {
		zap.S().Error("Error getting all production reports: ", err)
		return []types.ProductionReport{}, fmt.Errorf("error getting all production reports: %w", err)
	}
	return reports, nil
}

func (db *DB) GetRecentProductionReports(limit int) ([]types.ProductionReport, error) {
	var reports []types.ProductionReport
	query := `
        SELECT * FROM production_report
        ORDER BY end_date DESC
        LIMIT $1
    `
	err := db.Select(&reports, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent production reports: %w", err)
	}
	return reports, nil
}
