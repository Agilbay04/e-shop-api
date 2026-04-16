package main

import (
	"log"
	"e-shop-api/internal/config"
	"e-shop-api/internal/model"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 2. Connect DB
	db := config.ConnectDatabase()

	log.Println("Starting Migration...")

	// 3. Run Migration
	err := db.AutoMigrate(
		&model.User{},
		&model.Store{},
		&model.Product{},
	)

	if err != nil {
		log.Fatal("Migration Failed: ", err)
	}

	log.Println("Migration Success!")
}