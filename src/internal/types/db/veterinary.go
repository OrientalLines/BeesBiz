package types

import "time"

type VeterinaryPassport struct {
	PassportID         int       `json:"passport_id" db:"passport_id"`
	CommunityID        int       `json:"community_id" db:"community_id"`
	IssueDate          time.Time `json:"issue_date" db:"issue_date"`
	HealthStatus       string    `json:"health_status" db:"health_status"`
	LastInspectionDate time.Time `json:"last_inspection_date" db:"last_inspection_date"`
}

type VeterinaryRecord struct {
	RecordID    int       `json:"record_id" db:"record_id"`
	PassportID  int       `json:"passport_id" db:"passport_id"`
	RecordDate  time.Time `json:"record_date" db:"record_date"`
	Description string    `json:"description" db:"description"`
	Treatment   string    `json:"treatment" db:"treatment"`
}
