package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AlterOrderAddOrderNumberColumnMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260427115756",
		Migrate: func(tx *gorm.DB) error {
			return tx.Migrator().AddColumn(&struct {
			}{}, "order_number")
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&struct {
			}{}, "order_number")
		},
	}
}