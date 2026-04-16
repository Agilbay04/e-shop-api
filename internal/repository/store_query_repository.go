package repository

import (
	"e-shop-api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoreQueryRepository interface {
	FindAll() ([]model.Store, error)
	FindByID(id string) (*model.Store, error)
	FindByUserID(userID uuid.UUID) (*model.Store, error)
}

type storeQueryRepository struct {
	db *gorm.DB
}

func NewStoreQueryRepository(db *gorm.DB) StoreQueryRepository {
	return &storeQueryRepository{db}
}

func (r *storeQueryRepository) FindAll() ([]model.Store, error) {
	var stores []model.Store
	err := r.db.Preload("User").Find(&stores).Error
	return stores, err
}

func (r *storeQueryRepository) FindByID(id string) (*model.Store, error) {
	var store model.Store
	err := r.db.Preload("User").Where("id = ?", id).First(&store).Error
	return &store, err
}

func (r *storeQueryRepository) FindByUserID(userID uuid.UUID) (*model.Store, error) {
	var store model.Store
	err := r.db.Where("user_id = ?", userID).First(&store).Error
	return &store, err
}