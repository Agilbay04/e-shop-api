package repository

import (
	"e-shop-api/internal/model"
	"gorm.io/gorm"
)

type StoreRepository interface {
	Create(tx *gorm.DB, store *model.Store) error
	Update(tx *gorm.DB, store *model.Store) error
	Delete(tx *gorm.DB, id string) error
}

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db}
}

func (r *storeRepository) Create(tx *gorm.DB, store *model.Store) error {
	if tx != nil {
		return tx.Create(store).Error
	}
	return r.db.Create(store).Error
}

func (r *storeRepository) Update(tx *gorm.DB, store *model.Store) error {
	if tx != nil {
		return tx.Save(store).Error
	}
	return r.db.Save(store).Error
}

func (r *storeRepository) Delete(tx *gorm.DB, id string) error {
	if tx != nil {
		return tx.Delete(&model.Store{}, "id = ?", id).Error
	}
	return r.db.Delete(&model.Store{}, "id = ?", id).Error
}
