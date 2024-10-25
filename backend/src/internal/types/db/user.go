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