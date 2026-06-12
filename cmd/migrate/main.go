package main

import (
	"bagogo-boiler/internal/configs"
	"bagogo-boiler/internal/migrations"
	"bagogo-boiler/internal/pkg/logger"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load()
	logger.InitLogger()
	defer logger.L.Sync()

	logger.L.Info("Starting migrations...")

	db := configs.ConnectDatabase()

	if err := migrations.RunMigrations(db); err != nil {
		logger.L.Fatal("Migration failed", zap.Error(err))
	}

	logger.L.Info("Migrations completed successfully!")
}