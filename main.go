package main

import (
	"log"

	"cookaholic/internal/app"
	config "cookaholic/internal/infrastructure"
)

func main() {
	// Load environment variables from .env file
	config.LoadConfig()

	// Initialize application
	application, err := app.NewApplication()
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}

	// Start the application
	if err := application.Start(); err != nil {
		log.Fatal("Failed to start application:", err)
	}
}
