package seeders

import (
	"e-shop-api/internal/models"
	"errors"

	"gorm.io/gorm"
)

func SeedStores(db *gorm.DB) error {
	var seller models.User
	// Find user with role "seller"
	err := db.Where("role = ?", "seller").First(&seller).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("failed seed store: user seller not found")
		}
		return err
	}

	stores := []models.Store{
		{
			Name:        "Gadget Store",
			Description: "Toko elektronik terlengkap",
			UserID:      seller.ID,
		},
	}

	for _, s := range stores {
		err := db.FirstOrCreate(&s, models.Store{Name: s.Name}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
