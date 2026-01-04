package routes

import (
	"fmt"
	"go-fiber/bootstrap"
	"go-fiber/core/logs"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {

	api := app.Group("/api/v1", func(ctx *fiber.Ctx) error {
		return ctx.Next()
	})

	api.Post("/health", func(c *fiber.Ctx) error {
		hostname, err := os.Hostname()
		logs.Info(fmt.Sprintf("hostname: %s", hostname))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Error getting hostname: %s", err),
			})
		}
		currentTime := time.Now().Format(time.RFC3339)
		if err := c.Status(200).JSON(fiber.Map{
			"hostname":  hostname,
			"timestamp": currentTime,
			"version":   bootstrap.GlobalEnv.App.Version,
			"msg":       "Connect OK...!",
		}); err != nil {
			return err
		}
		return nil
	})

	NewUserRouter(api, db)

}

//update config
//try
