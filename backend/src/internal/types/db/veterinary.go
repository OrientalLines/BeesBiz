package types

import "github.com/guregu/null"

type VeterinaryPassport struct {
	PassportID         int       `json:"passport_id,omitempty" db:"passport_id"`
	BeeCommunityID     int       `json:"bee_community_id" db:"bee_community_id"`
	IssueDate          null.Time `json:"issue_date" db:"issue_date"`
	HealthStatus       string    `json:"health_status" db:"health_status"`
	LastInspectionDate null.Time `json:"last_inspection_date" db:"last_inspection_date"`
}

type VeterinaryRecord struct {
	RecordID    int       `json:"record_id,omitempty" db:"record_id"`
	PassportID  int       `json:"passport_id" db:"passport_id"`
	RecordDate  null.Time `json:"record_date" db:"record_date"`
	Description string    `json:"description" db:"description"`
	Treatment   string    `json:"treatment" db:"treatment"`
}
