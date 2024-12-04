package types

import "github.com/guregu/null"

type ProductionReport struct {
	ReportID      int       `json:"report_id,omitempty" db:"report_id"`
	ApiaryID      int       `json:"apiary_id" db:"apiary_id"`
	StartDate     null.Time `json:"start_date" db:"start_date"`
	EndDate       null.Time `json:"end_date" db:"end_date"`
	TotalHoney    int       `json:"total_honey_produced" db:"total_honey_produced"`
	TotalExpenses int       `json:"total_expenses" db:"total_expenses"`
	CuratedBy     int       `json:"curated_by" db:"curated_by"`
}
