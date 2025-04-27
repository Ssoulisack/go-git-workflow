package routes

import (
	"go-fiber/api/rest/controllers"
	"go-fiber/data/repositories"
	"go-fiber/data/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewUserRouter(router fiber.Router, db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo)
	userCtrl := controllers.NewUserCtrl(userSvc)

	userRoute := router.Group("/users", func(ctx *fiber.Ctx) error {
		return ctx.Next()
	})

	userRoute.Post("/", userCtrl.CreateUser)
	userRoute.Get("/:id", userCtrl.GetUserByID)
}