package main

import (
	"e-shop-api/internal/config"
	"e-shop-api/internal/migrations"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Connect DB
	db := config.ConnectDatabase()

	// Run Migration
	log.Println("Starting migrations...")

    if err := migrations.RunMigrations(db); err != nil {
        log.Fatalf("Migration failed: %v", err)
    }

    log.Println("Migrations completed successfully!")
}