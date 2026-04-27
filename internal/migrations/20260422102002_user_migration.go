package migrations

import (
	"e-shop-api/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func UserMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260422102002",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.User{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("users")
		},
	}
}
