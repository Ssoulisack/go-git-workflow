package controllers

import (
	"go-fiber/api/rest/middleware"
	"go-fiber/data/services"
	"go-fiber/domain/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserCtrl interface {
	//Methods
	CreateUser(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
}

type userCtrl struct {
	userSvc services.UserService
}

// CreateUser implements UserCtrl.
func (u *userCtrl) CreateUser(c *fiber.Ctx) error {
	// Parse request body
	var userReq models.UserReq
	if err := c.BodyParser(&userReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Call service to create user
	if err := u.userSvc.CreateUser(userReq); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return middleware.NewSuccessMessageResponse(c, "User created successfully")
}

// GetUserByID implements UserCtrl.
func (u *userCtrl) GetUserByID(c *fiber.Ctx) error {
	// Get user ID from URL parameters
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}
	// Convert ID to integer
	userID, err := strconv.Atoi(id)
	if err != nil {
		return middleware.ErrorExpectationFailed("Invalid user ID")
	}

	// Call service to get user by ID
	user, err := u.userSvc.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return middleware.ErrorNotFound("User not found")
	}
	return middleware.NewSuccessResponse(c, user)

}

func NewUserCtrl(userSvc services.UserService) UserCtrl {
	return &userCtrl{
		userSvc: userSvc,
	}
}
