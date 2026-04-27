package migrations

import (
	"e-shop-api/internal/models"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func OrderSequenceMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260427115731",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&models.OrderSequence{}); err != nil {
				return err
			}

			today := time.Now().Format("2006-01-02")
			return tx.Create(&models.OrderSequence{
				Date:         today,
				LastSequence: 0,
			}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("order_sequences")
		},
	}
}
