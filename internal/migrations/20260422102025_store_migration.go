package migrations

import (
	"e-shop-api/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func StoreMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260422102025",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Store{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("stores")
		},
	}
}
