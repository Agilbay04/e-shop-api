package main

import (
	"e-shop-api/internal/app"
	"e-shop-api/internal/config"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect database
	db := config.ConnectDatabase()

	// Setup router
	r := app.SetupRouter(db)

	// Run server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8001"
	}

	log.Printf("Server e-shop-api is starting on port %s", port)
	
	// Menjalankan server Gin
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}