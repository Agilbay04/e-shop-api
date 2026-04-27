package seeders

import (
	"e-shop-api/internal/models"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	users := []models.User{
		{Username: "admin", Email: "admin@mail.com", Password: "password123", Role: "admin"},
		{Username: "seller", Email: "seller@mail.com", Password: "password123", Role: "seller"},
		{Username: "buyer", Email: "buyer@mail.com", Password: "password123", Role: "buyer"},
	}

	for _, u := range users {
		// FirstOrCreate for preventing duplicate data
		err := db.FirstOrCreate(&u, models.User{Email: u.Email}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
