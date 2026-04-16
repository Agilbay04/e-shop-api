package seeder

import (
	"log"
	"gorm.io/gorm"
)

func RunSeeder(db *gorm.DB) {
	log.Println("Running database seeding...")

	// 1. Seed Users
	if err := SeedUsers(db); err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	// 2. Seed Stores
	if err := SeedStores(db); err != nil {
		log.Fatalf("Failed to seed stores: %v", err)
	}

	// 3. Seed Products
	if err := SeedProducts(db); err != nil {
		log.Fatalf("Failed to seed products: %v", err)
	}

	log.Println("Seeding completed successfully!")
}