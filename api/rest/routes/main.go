package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api/v1", func(ctx *fiber.Ctx) error {
		return ctx.Next()
	})

	NewUserRouter(api, db)

}
//update config
//try 