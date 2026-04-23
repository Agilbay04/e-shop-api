package seeder

import (
	"e-shop-api/internal/pkg/logger"
	"gorm.io/gorm"

	"go.uber.org/zap"
)

func RunSeeder(db *gorm.DB) {
	logger.InitLogger()
	defer logger.Log.Sync()

	logger.Log.Info("Running database seeding...")

	// 1. Seed Users
	if err := SeedUsers(db); err != nil {
		logger.Log.Fatal("Failed to seed users", zap.Error(err))
	}

	// 2. Seed Stores
	if err := SeedStores(db); err != nil {
		logger.Log.Fatal("Failed to seed stores", zap.Error(err))
	}

	// 3. Seed Products
	if err := SeedProducts(db); err != nil {
		logger.Log.Fatal("Failed to seed products", zap.Error(err))
	}

	logger.Log.Info("Seeding completed successfully!")
}