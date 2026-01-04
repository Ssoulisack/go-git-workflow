package main

import (
	"fmt"
	"go-fiber/api/rest/routes"
	"go-fiber/bootstrap"
	"go-fiber/core/logs"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	app := bootstrap.App()
	globalEnv := app.Env
	fiber := app.Fiber
	db := app.DB
	routes.Setup(fiber, db)

	go func() {
		for range time.Tick(10 * time.Second) {
			pgDB, _ := db.DB()
			stats := pgDB.Stats()
			logs.Info(fmt.Sprintf(
				"[DB Pool] Open=%d InUse=%d Idle=%d WaitCount=%d MaxOpen=%v\n",
				stats.OpenConnections,
				stats.InUse,
				stats.Idle,
				stats.WaitCount,
				globalEnv.Database.MaxOpenConns,
			))
		}
	}()

	// --- Graceful shutdown handler ---
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Fatal(fiber.Listen(fmt.Sprintf(":%v", globalEnv.App.Port)))
	}()

	// Wait for interrupt signal to gracefully shutdown
	<-stop // wait for Ctrl+C or container stop

	fmt.Println("Shutting down gracefully...")

	pgDB, _ := db.DB()
	pgDB.Close() // close DB pool cleanly
	fmt.Println("Database connection closed")

	// Shutdown Fiber app
	if err := fiber.Shutdown(); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	fmt.Println("Server stopped")
	// ------------------------------------------------
}
