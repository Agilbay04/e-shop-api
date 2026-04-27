package migrations

import (
	"e-shop-api/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AlterOrderAddOrderNumberColumnMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260427115756",
		Migrate: func(tx *gorm.DB) error {
			return tx.Migrator().AddColumn(&models.Order{}, "order_number")
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&models.Order{}, "order_number")
		},
	}
}
