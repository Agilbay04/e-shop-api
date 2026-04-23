package main

import (
	"e-shop-api/internal/config"
	"e-shop-api/internal/config/seeder"
	"e-shop-api/internal/pkg/logger"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	logger.InitLogger()
	defer logger.L.Sync()

	db := config.ConnectDatabase()

	if os.Getenv("APP_ENV") == "development" {
		logger.L.Info("Starting Seeder...")
		seeder.RunSeeder(db)
	}
}