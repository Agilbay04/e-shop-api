package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		UserMigration(),
		ProductMigration(),
		StoreMigration(),
		OrderMigration(),
		AlterUserAddPictureColumnMigration(),
		OrderSequenceMigration(),
		AlterOrderAddOrderNumberColumnMigration(),
	})

	return m.Migrate()
}
