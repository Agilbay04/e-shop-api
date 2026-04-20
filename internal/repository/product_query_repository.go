package repository

import (
	"e-shop-api/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductQueryRepository interface {
	FindAll() ([]model.Product, error)
	FindBySlug(slug string) (*model.Product, error)
	FindByID(id string) (*model.Product, error)
	FindByIDPreloadStore(id string) (*model.Product, error)
	FindByIDWithLock(tx *gorm.DB, id string) (*model.Product, error)
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
	err := r.db.Find(&products).Error
	return products, err
}

func (r *productQueryRepository) FindBySlug(slug string) (*model.Product, error) {
	var product model.Product
	err := r.db.Where("slug = ?", slug).First(&product).Error
	return &product, err
}

func (r *productQueryRepository) FindByID(id string) (*model.Product, error) {
	var product model.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	return &product, err
}

func (r *productQueryRepository) FindByIDPreloadStore(id string) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Store").Where("id = ?", id).First(&product).Error
	return &product, err
}

func (r *productQueryRepository) FindByIDWithLock(tx *gorm.DB, id string) (*model.Product, error) {
    var product model.Product

    err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Store").
		First(&product, "id = ?", id).
		Error
    
	if err != nil {
		return nil, err
	}

    return &product, nil
}