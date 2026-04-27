package repositories

import (
	"e-shop-api/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderSequenceRepository interface {
	GetNextSequence(tx *gorm.DB, date string) (int, error)
}

type orderSequenceRepository struct {
	db *gorm.DB
}

func NewOrderSequenceRepository(db *gorm.DB) OrderSequenceRepository {
	return &orderSequenceRepository{db}
}

func (r *orderSequenceRepository) GetNextSequence(tx *gorm.DB, date string) (int, error) {
	db := r.db
	if tx != nil {
		db = tx
	}

	var seq model.OrderSequence
	err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("date = ?", date).
		First(&seq).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := db.Create(&model.OrderSequence{
			Date:          date,
			LastSequence: 0,
		}).Error; err != nil {
			return 0, err
		}

		err = db.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("date = ?", date).
			First(&seq).Error
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}

	seq.LastSequence++
	if err := db.Save(&seq).Error; err != nil {
		return 0, err
	}

	return seq.LastSequence, nil
}

func InitOrderSequence(db *gorm.DB) error {
	today := time.Now().Format("2006-01-02")
	return db.FirstOrCreate(&model.OrderSequence{}, &model.OrderSequence{
		Date:          today,
		LastSequence:  0,
	}).Error
}