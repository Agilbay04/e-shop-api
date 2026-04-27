package main

import (
	"e-shop-api/internal/configs"
	"e-shop-api/internal/pkg/logger"
	"e-shop-api/internal/seeders"
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