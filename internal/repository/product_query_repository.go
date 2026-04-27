package repository

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/model"
	"e-shop-api/internal/pkg/util"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductQueryRepository interface {
	FindAllPagination(req dto.QueryProductRequest, user dto.CurrentUser) ([]dto.ProductResponse, int64, error)
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

func (r *productQueryRepository) FindAllPagination(req dto.QueryProductRequest, user dto.CurrentUser) ([]dto.ProductResponse, int64, error) {
	var products []model.Product
	var total int64

	// Apply filters for COUNT query
	countQuery := r.db.Model(&model.Product{})
	r.applyFilters(countQuery, req)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply filters for SELECT query (fresh query)
	query := r.db.Model(&model.Product{})
	r.applyFilters(query, req)

	// Pagination & Sorting
	err := query.Scopes(util.Paginate(req.Page, req.Limit)).
		Order(req.SortBy + " " + req.OrderBy).
		Preload("Store").
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	// Map model.Product to dto.ProductResponse
	productsResponse := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		productsResponse[i] = dto.ProductResponse{
			ID:          product.ID.String(),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Unit:        product.Unit,
			IsActive:    product.IsActive,
			CreatedAt:   product.CreatedAt.Format(time.RFC3339),
			CreatedBy:  product.CreatedBy.String(),
			UpdatedAt: product.UpdatedAt.Format(time.RFC3339),
			UpdatedBy:  product.UpdatedBy.String(),
			DeletedAt:  product.DeletedAt.Time.Format(time.RFC3339),
			StoreID:    product.StoreID.String(),
			StoreName: product.Store.Name,
		}
	}

	return productsResponse, total, nil
}

func (r *productQueryRepository) applyFilters(query *gorm.DB, req dto.QueryProductRequest) *gorm.DB {
	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}
	if req.StoreID != nil {
		query = query.Where("store_id = ?", *req.StoreID)
	}
	if req.ID != nil {
		query = query.Where("id = ?", *req.ID)
	}
	if req.MinPrice != nil {
		query = query.Where("price >= ?", *req.MinPrice)
	}
	if req.MaxPrice != nil {
		query = query.Where("price <= ?", *req.MaxPrice)
	}
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}
	return query
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