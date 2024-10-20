package types

type Role string

const (
	Admin  Role = "ADMIN"
	Worker Role = "WORKER"
)

type User struct {
	UserID    int    `json:"user_id" db:"user_id"`
	Username  string `json:"username" db:"username"`
	Role      Role   `json:"role" db:"role"`
	Email     string `json:"email" db:"email"`
	LastLogin string `json:"last_login" db:"last_login"`
}
