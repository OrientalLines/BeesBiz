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
	err := db.Get(&createdGroup, `
		INSERT INTO worker_group (manager_id, group_name) 
		VALUES ($1, $2) 
		RETURNING *`, 
		group.ManagerID, group.GroupName)
	if err != nil {
		zap.S().Error("Error creating worker group: ", err)
		return types.WorkerGroup{}, fmt.Errorf("error creating worker group: %w", err)
	}
	return createdGroup, nil
}

// AddWorkerToGroup adds a worker to a group
func (db *DB) AddWorkerToGroup(groupID, workerID int) error {
	_, err := db.Exec(`
		INSERT INTO worker_group_member (group_id, worker_id)
		VALUES ($1, $2)`,
		groupID, workerID)
	if err != nil {
		zap.S().Error("Error adding worker to group: ", err)
		return fmt.Errorf("error adding worker to group: %w", err)
	}
	return nil
}

// RemoveWorkerFromGroup removes a worker from a group
func (db *DB) RemoveWorkerFromGroup(groupID, workerID int) error {
	_, err := db.Exec(`
		DELETE FROM worker_group_member 
		WHERE group_id = $1 AND worker_id = $2`,
		groupID, workerID)
	if err != nil {
		zap.S().Error("Error removing worker from group: ", err)
		return fmt.Errorf("error removing worker from group: %w", err)
	}
	return nil
}

// GetGroupMembers retrieves all workers in a specific group
func (db *DB) GetGroupMembers(groupID int) ([]types.User, error) {
	var users []types.User
	err := db.Select(&users, `
		SELECT u.* FROM "user" u
		JOIN worker_group_member wgm ON u.user_id = wgm.worker_id
		WHERE wgm.group_id = $1`,
		groupID)
	if err != nil {
		zap.S().Error("Error getting group members: ", err)
		return nil, fmt.Errorf("error getting group members: %w", err)
	}
	return users, nil
}

// GetWorkerGroups retrieves all groups a worker belongs to
func (db *DB) GetWorkerGroups(workerID int) ([]types.WorkerGroup, error) {
	var groups []types.WorkerGroup
	err := db.Select(&groups, `
		SELECT wg.* FROM worker_group wg
		JOIN worker_group_member wgm ON wg.group_id = wgm.group_id
		WHERE wgm.worker_id = $1`,
		workerID)
	if err != nil {
		zap.S().Error("Error getting worker's groups: ", err)
		return nil, fmt.Errorf("error getting worker's groups: %w", err)
	}
	return groups, nil
}

func (db *DB) GetWorkerGroupsByManager(managerID int) ([]types.WorkerGroup, error) {
	var groups []types.WorkerGroup
	err := db.Select(&groups, "SELECT * FROM worker_group WHERE manager_id = $1", managerID)
	if err != nil {
		zap.S().Error("Error getting worker groups by manager: ", err)
		return nil, fmt.Errorf("error getting worker groups by manager: %w", err)
	}
	return groups, nil
}
func (db *DB) GetFreeUsers() ([]types.User, error) {
	var users []types.User
	err := db.Select(&users, `
		SELECT u.* FROM "user" u
		WHERE u.role = 'WORKER'
		AND NOT EXISTS (
			SELECT 1 FROM worker_group_member wgm 
			WHERE wgm.worker_id = u.user_id
		)`)
	if err != nil {
		zap.S().Error("Error getting free users: ", err)
		return nil, fmt.Errorf("error getting free users: %w", err)
	}
	return users, nil
}

func (db *DB) DeleteWorkerGroup(id int) error {
	// First delete all members from this group
	_, err := db.Exec("DELETE FROM worker_group_member WHERE group_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting worker group members: ", err)
		return fmt.Errorf("error deleting worker group members: %w", err)
	}

	// Then delete the group itself
	_, err = db.Exec("DELETE FROM worker_group WHERE group_id = $1", id)
	if err != nil {
		zap.S().Error("Error deleting worker group: ", err)
		return fmt.Errorf("error deleting worker group: %w", err)
	}
	return nil
}

func (db *DB) UpdateWorkerGroup(id int, group types.WorkerGroup) (types.WorkerGroup, error) {
	var updatedGroup types.WorkerGroup
	err := db.Get(&updatedGroup, "UPDATE worker_group SET group_name = $1 WHERE group_id = $2 RETURNING *", group.GroupName, id)
	if err != nil {
		zap.S().Error("Error updating worker group: ", err)
		return types.WorkerGroup{}, fmt.Errorf("error updating worker group: %w", err)
	}
	return updatedGroup, nil
}

func (db *DB) GetAllWorkerGroups() ([]types.WorkerGroup, error) {
	var groups []types.WorkerGroup
	err := db.Select(&groups, "SELECT * FROM worker_group")
	if err != nil {
		zap.S().Error("Error getting all worker groups: ", err)
		return nil, fmt.Errorf("error getting all worker groups: %w", err)
	}
	return groups, nil
}
