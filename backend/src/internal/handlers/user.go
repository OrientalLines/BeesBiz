package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/orientallines/beesbiz/internal/database"
	types "github.com/orientallines/beesbiz/internal/types/db"
)

// User handlers
func GetUser(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid user ID: %v", err)})
		}
		user, err := db.GetUser(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get user: %v", err)})
		}
		return c.JSON(user)
	}
}

func CreateUser(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user types.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid user data: %v", err)})
		}
		createdUser, err := db.CreateUser(user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create user: %v", err)})
		}
		return c.JSON(createdUser)
	}
}

func UpdateUser(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user types.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid user data: %v", err)})
		}
		updatedUser, err := db.UpdateUser(user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update user: %v", err)})
		}
		return c.JSON(updatedUser)
	}
}

func DeleteUser(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid user ID: %v", err)})
		}
		if err := db.DeleteUser(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete user: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetAllUsers(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := db.GetAllUsers()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all users: %v", err)})
		}
		return c.JSON(users)
	}
}

func GetUserAllowedRegions(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid user ID: %v", err)})
		}
		regions, err := db.GetAllowedRegions(userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get allowed regions: %v", err)})
		}
		return c.JSON(regions)
	}
}

func ModifyUserRole(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type RoleUpdate struct {
			UserID int        `json:"user_id"`
			Role   types.Role `json:"role"`
		}

		var update RoleUpdate
		if err := c.BodyParser(&update); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid request data: %v", err)})
		}

		// Get existing user
		user, err := db.GetUser(update.UserID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		// Update role
		user.Role = update.Role
		updatedUser, err := db.UpdateUser(user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update user role: %v", err)})
		}

		return c.JSON(updatedUser)
	}
}

func ModifyUserAllowedRegions(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type RegionsUpdate struct {
			UserID    int   `json:"user_id"`
			RegionIDs []int `json:"region_ids"`
		}

		var update RegionsUpdate
		if err := c.BodyParser(&update); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid request data: %v", err)})
		}

		// First, delete existing allowed regions
		if err := db.DeleteAllowedRegionsForUser(update.UserID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update allowed regions: %v", err)})
		}

		// Then add new allowed regions
		for _, regionID := range update.RegionIDs {
			allowedRegion := types.AllowedRegion{
				UserID:   update.UserID,
				RegionID: regionID,
			}
			if _, err := db.CreateAllowedRegion(allowedRegion); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to add allowed region: %v", err)})
			}
		}

		// Return updated list of allowed regions
		regions, err := db.GetAllowedRegions(update.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get updated allowed regions: %v", err)})
		}

		return c.JSON(regions)
	}
}

// GetWorkerGroup retrieves a worker group by its ID.
func GetWorkerGroup(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid worker group ID: %v", err)})
		}
		group, err := db.GetWorkerGroup(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get worker group: %v", err)})
		}
		return c.JSON(group)
	}
}

// CreateWorkerGroup creates a new worker group.
func CreateWorkerGroup(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var group types.WorkerGroup
		if err := c.BodyParser(&group); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid worker group data: %v", err)})
		}
		createdGroup, err := db.CreateWorkerGroup(group)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create worker group: %v", err)})
		}
		return c.JSON(createdGroup)
	}
}

// UpdateWorkerGroup updates an existing worker group.
func UpdateWorkerGroup(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var group types.WorkerGroup
		if err := c.BodyParser(&group); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid worker group data: %v", err)})
		}
		updatedGroup, err := db.UpdateWorkerGroup(group)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update worker group: %v", err)})
		}
		return c.JSON(updatedGroup)
	}
}

// DeleteWorkerGroup deletes a worker group by its ID.
func DeleteWorkerGroup(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid worker group ID: %v", err)})
		}
		if err := db.DeleteWorkerGroup(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to delete worker group: %v", err)})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// GetAllWorkerGroups retrieves all worker groups.
func GetAllWorkerGroups(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		groups, err := db.GetAllWorkerGroups()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get all worker groups: %v", err)})
		}
		return c.JSON(groups)
	}
}

// GetWorkerGroupsByManager retrieves all worker groups managed by a specific manager.
func GetWorkerGroupsByManager(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		managerID, err := c.ParamsInt("manager_id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid manager ID: %v", err)})
		}
		groups, err := db.GetWorkerGroupsByManager(managerID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to get worker groups by manager: %v", err)})
		}
		return c.JSON(groups)
	}
}
