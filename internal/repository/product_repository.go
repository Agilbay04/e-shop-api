package repository

import (
	"e-shop-api/internal/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *model.Product) error
	Update(product *model.Product) error
	Delete(product *model.Product) error
	UpdateStock(tx *gorm.DB, id string, newStock int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) UpdateStock(tx *gorm.DB, id string, newStock int) error {
    return tx.Model(&model.Product{}).Where("id = ?", id).Update("stock", newStock).Error
}