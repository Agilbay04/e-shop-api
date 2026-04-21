package repository

import (
	"e-shop-api/internal/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(tx *gorm.DB, product *model.Product) error
	Update(tx *gorm.DB, product *model.Product) error
	Delete(tx *gorm.DB, product *model.Product) error
	UpdateStock(tx *gorm.DB, id string, newStock int) error
	AddStock(tx *gorm.DB, id string, quantity int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(tx *gorm.DB, product *model.Product) error {
	if tx != nil {
		return tx.Create(product).Error
	}
	return r.db.Create(product).Error
}

func (r *productRepository) Update(tx *gorm.DB, product *model.Product) error {
	if tx != nil {
		return tx.Save(product).Error
	}
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(tx *gorm.DB, product *model.Product) error {
	if tx != nil {
		return tx.Save(product).Error
	}
	return r.db.Save(product).Error
}

func (r *productRepository) UpdateStock(tx *gorm.DB, id string, newStock int) error {
	return tx.Model(&model.Product{}).Where("id = ?", id).Update("stock", newStock).Error
}

func (r *productRepository) AddStock(tx *gorm.DB, id string, quantity int) error {
    return tx.Model(&model.Product{}).
        Where("id = ?", id).
        Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

