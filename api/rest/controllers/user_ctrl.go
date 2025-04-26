package controllers

import (
	"go-fiber/data/services"

	"github.com/gofiber/fiber/v2"
)

type UserCtrl interface {
	//Methods
	GetAllUsers(ctx *fiber.Ctx) error
}

type userCtrl struct {
	userSvc services.UserService
}

func (u *userCtrl) GetAllUsers(ctx *fiber.Ctx) error {
	u.userSvc.GetAllUsers()
	return nil
}

func NewUserCtrl(userSvc services.UserService) UserCtrl {
	return &userCtrl{
		userSvc: userSvc,
	}
}
