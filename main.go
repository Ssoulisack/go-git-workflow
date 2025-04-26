package main

import (
	"fmt"
	"go-fiber/api/rest/routes"
	"go-fiber/bootstrap"
	"log"
)

func main() {
	app := bootstrap.App()
	globalEnv := app.Env
	fiber := app.Fiber
	db := app.DB
	routes.Setup(fiber, db)

	log.Fatal(fiber.Listen(fmt.Sprintf(":%d", globalEnv.App.Port)))
}
