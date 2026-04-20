package main

import (
	"log"
	"e-shop-api/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 2. Connect DB
	db := config.ConnectDatabase()

	// 3. Run Migration
	config.RunMigration(db)
}