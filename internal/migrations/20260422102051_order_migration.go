package migrations

import (
	"e-shop-api/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func OrderMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260422102051",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.Order{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("orders")
		},
	}
}
