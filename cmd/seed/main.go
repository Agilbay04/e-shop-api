package main

import (
	"e-shop-api/internal/config"
	"e-shop-api/internal/config/seeder"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db := config.ConnectDatabase()

	if (os.Getenv("APP_ENV") == "development") {
		log.Println("Starting Seeder...")
		seeder.RunSeeder(db)	
	}
}