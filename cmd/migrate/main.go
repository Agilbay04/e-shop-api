package main

import (
	"e-shop-api/internal/config"
	"e-shop-api/internal/migrations"
	"e-shop-api/internal/pkg/logger"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load()
	logger.InitLogger()
	defer logger.Log.Sync()

	logger.Log.Info("Starting migrations...")

	db := config.ConnectDatabase()

	if err := migrations.RunMigrations(db); err != nil {
		logger.Log.Fatal("Migration failed", zap.Error(err))
	}

	logger.Log.Info("Migrations completed successfully!")
}