package migrations

import (
	"e-shop-api/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AlterUserAddPictureColumnMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260422110417",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.User{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&models.User{}, "picture")
		},
	}
}
