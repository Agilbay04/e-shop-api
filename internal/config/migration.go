package config

import (
	"log"
	"e-shop-api/internal/model"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	log.Println("Running database migration...")
	
	// Register models
	err := db.AutoMigrate(
		&model.User{},
		&model.Store{},
		&model.Product{},
		&model.Order{},
		&model.OrderItem{},
		
		// Add more models here
	)

	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	log.Println("Migration completed successfully!")
}