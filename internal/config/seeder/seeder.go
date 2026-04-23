package seeder

import (
	"e-shop-api/internal/pkg/logger"
	"gorm.io/gorm"

	"go.uber.org/zap"
)

func RunSeeder(db *gorm.DB) {
	logger.InitLogger()
	defer logger.L.Sync()

	logger.L.Info("Running database seeding...")

	// 1. Seed Users
	if err := SeedUsers(db); err != nil {
		logger.L.Fatal("Failed to seed users", zap.Error(err))
	}

	// 2. Seed Stores
	if err := SeedStores(db); err != nil {
		logger.L.Fatal("Failed to seed stores", zap.Error(err))
	}

	// 3. Seed Products
	if err := SeedProducts(db); err != nil {
		logger.L.Fatal("Failed to seed products", zap.Error(err))
	}

	logger.L.Info("Seeding completed successfully!")
}