package seeders

import (
	"e-shop-api/internal/models"
	"errors"

	"gorm.io/gorm"
)

func SeedProducts(db *gorm.DB) error {
	var store models.Store
	// Find the store that was created previously
	err := db.First(&store).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("failed seed product: store not found")
		}
		return err
	}

	products := []models.Product{
		{
			Name:    "Macbook Pro M2",
			Price:   20000000,
			Stock:   10,
			StoreID: store.ID,
		},
		{
			Name:    "iPhone 15 Pro",
			Price:   18000000,
			Stock:   5,
			StoreID: store.ID,
		},
	}

	for _, p := range products {
		err := db.FirstOrCreate(&p, models.Product{Name: p.Name}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
