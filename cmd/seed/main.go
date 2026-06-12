package main

import (
	"bagogo-boiler/internal/configs"
	"bagogo-boiler/internal/pkg/logger"
	"bagogo-boiler/internal/seeders"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	logger.InitLogger()
	defer logger.L.Sync()

	db := configs.ConnectDatabase()

	if os.Getenv("APP_ENV") == "development" {
		logger.L.Info("Starting Seeder...")
		seeders.RunSeeder(db)
	}
}