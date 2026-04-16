package repository

import (
	"e-shop-api/internal/model"
	"gorm.io/gorm"
)

type ProductQueryRepository interface {
	FindAll() ([]model.Product, error)
	FindBySlug(slug string) (*model.Product, error)
	FindByID(id string) (*model.Product, error)
}

type productQueryRepository struct {
	db *gorm.DB
}

func NewProductQueryRepository(db *gorm.DB) ProductQueryRepository {
	return &productQueryRepository{db}
}

func (r *productQueryRepository) FindAll() ([]model.Product, error) {
	var products []model.Product
	// Preload "Store" agar data toko muncul di JSON produk
	err := r.db.Preload("Store").Find(&products).Error
	return products, err
}

func (r *productQueryRepository) FindBySlug(slug string) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Store").Where("slug = ?", slug).First(&product).Error
	return &product, err
}

func (r *productQueryRepository) FindByID(id string) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Store").Where("id = ?", id).First(&product).Error
	return &product, err
}