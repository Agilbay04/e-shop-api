package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AlterOrderAddOrderNumberColumnMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260427115756",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
				ALTER TABLE orders 
				ADD COLUMN IF NOT EXISTS order_number VARCHAR(50) UNIQUE
			`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`
				ALTER TABLE orders DROP COLUMN IF EXISTS order_number
			`).Error
		},
	}
}
