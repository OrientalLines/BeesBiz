package database

import (
	"fmt"

	types "github.com/orientallines/beesbiz/internal/types/db"
)

func (db *DB) GetProductionReport(id int) (types.ProductionReport, error) {
	var report types.ProductionReport
	err := db.Get(&report, "SELECT * FROM production_report WHERE report_id = $1", id)
	if err != nil {
		return types.ProductionReport{}, fmt.Errorf("error getting production report: %w", err)
	}
	return report, nil
}

func (db *DB) CreateProductionReport(report types.ProductionReport) (types.ProductionReport, error) {
	var createdReport types.ProductionReport
	err := db.Get(&createdReport, "INSERT INTO production_report (apiary_id, start_date, end_date, total_honey_produced, total_expenses) VALUES ($1, $2, $3, $4, $5) RETURNING *", report.ApiaryID, report.StartDate, report.EndDate, report.TotalHoney, report.TotalExpenses)
	if err != nil {
		return types.ProductionReport{}, fmt.Errorf("error creating production report: %w", err)
	}
	return createdReport, nil
}

func (db *DB) UpdateProductionReport(report types.ProductionReport) (types.ProductionReport, error) {
	var updatedReport types.ProductionReport
	err := db.Get(&updatedReport, "UPDATE production_report SET apiary_id = $1, start_date = $2, end_date = $3, total_honey_produced = $4, total_expenses = $5 WHERE report_id = $6 RETURNING *", report.ApiaryID, report.StartDate, report.EndDate, report.TotalHoney, report.TotalExpenses, report.ReportID)
	if err != nil {
		return types.ProductionReport{}, fmt.Errorf("error updating production report: %w", err)
	}
	return updatedReport, nil
}

func (db *DB) DeleteProductionReport(id int) error {
	_, err := db.Exec("DELETE FROM production_report WHERE report_id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting production report: %w", err)
	}
	return nil
}

func (db *DB) GetAllProductionReports() ([]types.ProductionReport, error) {
	var reports []types.ProductionReport
	err := db.Select(&reports, "SELECT * FROM production_report")
	if err != nil {
		return []types.ProductionReport{}, fmt.Errorf("error getting all production reports: %w", err)
	}
	return reports, nil
}
