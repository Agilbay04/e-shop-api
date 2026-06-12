package migrations

import (
	"bagogo-boiler/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AlterOrderAddOrderNumberColumnMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260427115756",
		Migrate: func(tx *gorm.DB) error {
			if !tx.Migrator().HasColumn(&models.Order{}, "order_number") {
				return tx.Migrator().AddColumn(&models.Order{}, "order_number")
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&models.Order{}, "order_number")
		},
	}
}
