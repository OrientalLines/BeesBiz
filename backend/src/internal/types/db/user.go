package types

import "github.com/guregu/null"

type Role string

const (
	Admin   Role = "ADMIN"
	Worker  Role = "WORKER"
	Manager Role = "MANAGER"
)

type User struct {
	UserID    int       `json:"user_id,omitempty" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	FullName  string    `json:"full_name" db:"full_name"`
	Role      Role      `json:"role" db:"role"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	LastLogin null.Time `json:"last_login" db:"last_login"`
}

type WorkerGroup struct {
	GroupID   int       `db:"group_id" json:"group_id"`
	ManagerID int       `db:"manager_id" json:"manager_id"`
	WorkerID  int       `db:"worker_id" json:"worker_id"`
	GroupName string    `db:"group_name" json:"group_name"`
	CreatedAt null.Time `db:"created_at" json:"created_at"`
	UpdatedAt null.Time `db:"updated_at" json:"updated_at"`
}
