package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Application struct {
	Env   *Env
	Fiber *fiber.App
	DB    *gorm.DB
}

var GlobalEnv Env

func App() *Application {
	app := &Application{}
	app.Env = NewEnv()
	GlobalEnv = *NewEnv()
	app.Fiber = NewFiber()
	app.DB = NewDatabaseConnectionPostgres(app.Env)
	return app
}
