package database

import (
	"fmt"

	types "github.com/orientallines/beesbiz/internal/types/db"
)

func (db *DB) GetUser(id int) (types.User, error) {
	var user types.User
	err := db.Get(&user, "SELECT * FROM user WHERE user_id = $1", id)
	if err != nil {
		return types.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return user, nil
}

func (db *DB) CreateUser(user types.User) (types.User, error) {
	var createdUser types.User
	err := db.Get(&createdUser, "INSERT INTO user (username, role, email, last_login) VALUES ($1, $2, $3, $4) RETURNING *", user.Username, user.Role, user.Email, user.LastLogin)
	if err != nil {
		return types.User{}, fmt.Errorf("error creating user: %w", err)
	}
	return createdUser, nil
}

func (db *DB) UpdateUser(user types.User) (types.User, error) {
	var updatedUser types.User
	err := db.Get(&updatedUser, "UPDATE user SET username = $1, role = $2, email = $3, last_login = $4 WHERE user_id = $5 RETURNING *", user.Username, user.Role, user.Email, user.LastLogin, user.UserID)
	if err != nil {
		return types.User{}, fmt.Errorf("error updating user: %w", err)
	}
	return updatedUser, nil
}

func (db *DB) DeleteUser(id int) error {
	_, err := db.Exec("DELETE FROM user WHERE user_id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}

func (db *DB) GetAllUsers() ([]types.User, error) {
	var users []types.User
	err := db.Select(&users, "SELECT * FROM user")
	if err != nil {
		return []types.User{}, fmt.Errorf("error getting all users: %w", err)
	}
	return users, nil
}
func (db *DB) GetAllowedRegion(id int) (types.AllowedRegion, error) {
	var allowedRegion types.AllowedRegion
	err := db.Get(&allowedRegion, "SELECT * FROM allowed_region WHERE id = $1", id)
	if err != nil {
		return types.AllowedRegion{}, fmt.Errorf("error getting allowed region: %w", err)
	}
	return allowedRegion, nil
}
