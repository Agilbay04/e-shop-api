package repository

import (
	"e-shop-api/internal/model"
	"gorm.io/gorm"
)

type StoreRepository interface {
	Create(store *model.Store) error
	Update(store *model.Store) error
	Delete(id string) error
}

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db}
}

func (r *storeRepository) Create(store *model.Store) error {
	return r.db.Create(store).Error
}

func (r *storeRepository) Update(store *model.Store) error {
	return r.db.Save(store).Error
}

func (r *storeRepository) Delete(id string) error {
	return r.db.Delete(&model.Store{}, "id = ?", id).Error
}