package types

import "time"

type ProductionReport struct {
	ReportID      int       `json:"report_id" db:"report_id"`
	ApiaryID      int       `json:"apiary_id" db:"apiary_id"`
	StartDate     time.Time `json:"start_date" db:"start_date"`
	EndDate       time.Time `json:"end_date" db:"end_date"`
	TotalHoney    int       `json:"total_honey_produced" db:"total_honey_produced"`
	TotalExpenses int       `json:"total_expenses" db:"total_expenses"`
}
