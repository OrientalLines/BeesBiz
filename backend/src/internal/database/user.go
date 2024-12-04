package database

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	types "github.com/orientallines/beesbiz/internal/types/db"
)

func (db *DB) GetUser(id int) (types.User, error) {
	var user types.User
	err := db.Get(&user, "SELECT * FROM \"user\" WHERE user_id = $1", id)
	if err != nil {
		zap.S().Error("Error getting user: ", err)
		return types.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return user, nil
}

func (db *DB) CreateUser(user types.User) (types.User, error) {
	var createdUser types.User
	err := db.Get(&createdUser, "INSERT INTO \"user\" (username, full_name, role, email, password, last_login) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *", user.Username, user.FullName, strings.ToUpper(string(user.Role)), user.Email, user.Password, user.LastLogin)
	if err != nil {
		zap.S().Error("Error creating user: ", err)
		return types.User{}, fmt.Errorf("error creating user: %w", err)
	}
	return createdUser, nil
}

func (db *DB) UpdateUser(user types.User) (types.User, error) {
	var updatedUser types.User
	err := db.Get(&updatedUser, "UPDATE \"user\" SET username = $1, full_name = $2, role = $3, email = $4, password = $5, last_login = $6 WHERE user_id = $7 RETURNING *", user.Username, user.FullName, user.Role, user.Email, user.Password, user.LastLogin, user.UserID)
	if err != nil {
		zap.S().Error("Error updating user: ", err)
		return types.User{}, fmt.Errorf("error updating user: %w", err)
	}
	return updatedUser, nil
}

func (db *DB) DeleteUser(id int) error {
	_, err := db.Exec("DELETE FROM \"user\" WHERE user_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting user: ", err)
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}

func (db *DB) GetAllUsers() ([]types.User, error) {
	var users []types.User
	err := db.Select(&users, "SELECT * FROM \"user\"")
	if err != nil {
		zap.S().Error("Error getting all users: ", err)
		return []types.User{}, fmt.Errorf("error getting all users: %w", err)
	}
	return users, nil
}
func (db *DB) GetAllowedRegions(userID int) ([]types.AllowedRegion, error) {
	var allowedRegions []types.AllowedRegion
	err := db.Select(&allowedRegions, "SELECT * FROM allowed_region WHERE user_id = $1", userID)
	if err != nil {
		zap.S().Error("Error getting allowed regions: ", err)
		return []types.AllowedRegion{}, fmt.Errorf("error getting allowed regions: %w", err)
	}
	return allowedRegions, nil
}

func (db *DB) GetUserByEmail(email string) (types.User, error) {
	var user types.User
	err := db.Get(&user, "SELECT * FROM \"user\" WHERE email = $1", email)
	if err != nil {
		zap.S().Error("Error getting user by email: ", err)
		return types.User{}, fmt.Errorf("error getting user by email: %w", err)
	}
	return user, nil
}

func (db *DB) GetUserByUsername(username string) (types.User, error) {
	var user types.User
	err := db.Get(&user, "SELECT * FROM \"user\" WHERE username = $1", username)
	if err != nil {
		zap.S().Error("Error getting user by username: ", err)
		return types.User{}, fmt.Errorf("error getting user by username: %w", err)
	}
	return user, nil
}

// GetWorkerGroup retrieves a worker group by its ID.
func (db *DB) GetWorkerGroup(id int) (types.WorkerGroup, error) {
	var group types.WorkerGroup
	err := db.Get(&group, "SELECT * FROM worker_group WHERE group_id = $1", id)
	if err != nil {
		zap.S().Error("Error getting worker group: ", err)
		return types.WorkerGroup{}, fmt.Errorf("error getting worker group: %w", err)
	}
	return group, nil
}

// CreateWorkerGroup creates a new worker group and returns the created group.
func (db *DB) CreateWorkerGroup(group types.WorkerGroup) (types.WorkerGroup, error) {
	var createdGroup types.WorkerGroup
	err := db.Get(&createdGroup, "INSERT INTO worker_group (manager_id, worker_id, group_name) VALUES ($1, $2, $3) RETURNING *", group.ManagerID, group.WorkerID, group.GroupName)
	if err != nil {
		zap.S().Error("Error creating worker group: ", err)
		return types.WorkerGroup{}, fmt.Errorf("error creating worker group: %w", err)
	}
	return createdGroup, nil
}

// UpdateWorkerGroup updates an existing worker group and returns the updated group.
func (db *DB) UpdateWorkerGroup(group types.WorkerGroup) (types.WorkerGroup, error) {
	var updatedGroup types.WorkerGroup
	err := db.Get(&updatedGroup, "UPDATE worker_group SET manager_id = $1, worker_id = $2, group_name = $3 WHERE group_id = $4 RETURNING *", group.ManagerID, group.WorkerID, group.GroupName, group.GroupID)
	if err != nil {
		zap.S().Error("Error updating worker group: ", err)
		return types.WorkerGroup{}, fmt.Errorf("error updating worker group: %w", err)
	}
	return updatedGroup, nil
}

// DeleteWorkerGroup deletes a worker group by its ID.
func (db *DB) DeleteWorkerGroup(id int) error {
	_, err := db.Exec("DELETE FROM worker_group WHERE group_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting worker group: ", err)
		return fmt.Errorf("error deleting worker group: %w", err)
	}
	return nil
}

// GetAllWorkerGroups retrieves all worker groups.
func (db *DB) GetAllWorkerGroups() ([]types.WorkerGroup, error) {
	var groups []types.WorkerGroup
	err := db.Select(&groups, "SELECT * FROM worker_group")
	if err != nil {
		zap.S().Error("Error getting all worker groups: ", err)
		return []types.WorkerGroup{}, fmt.Errorf("error getting all worker groups: %w", err)
	}
	return groups, nil
}

// GetWorkerGroupsByManager retrieves all worker groups managed by a specific manager.
func (db *DB) GetWorkerGroupsByManager(managerID int) ([]types.WorkerGroup, error) {
	var groups []types.WorkerGroup
	query := `
        SELECT * FROM worker_group
        WHERE manager_id = $1
    `
	err := db.Select(&groups, query, managerID)
	if err != nil {
		zap.S().Error("Error getting worker groups by manager: ", err)
		return nil, fmt.Errorf("error getting worker groups by manager: %w", err)
	}
	return groups, nil
}
