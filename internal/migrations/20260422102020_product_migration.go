package migrations

import (
	"e-shop-api/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func ProductMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260422102020",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.Product{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("products")
		},
	}
}
