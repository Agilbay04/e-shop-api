package repository

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/model"
	"e-shop-api/internal/pkg/util"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StoreQueryRepository interface {
	FindAll() ([]model.Store, error)
	FindByID(id string) (*model.Store, error)
	FindByUserID(userID string) (*model.Store, error)
	FindAllPagination(req dto.QueryStoreParam) ([]dto.StoreResponse, int64, error)
	FindByIDWithLock(tx *gorm.DB, id string) (*model.Store, error)
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

func (r *storeQueryRepository) FindByUserID(userID string) (*model.Store, error) {
	var store model.Store
	err := r.db.Where("user_id = ?", userID).First(&store).Error
	return &store, err
}

func (r *storeQueryRepository) FindAllPagination(req dto.QueryStoreParam) ([]dto.StoreResponse, int64, error) {
	var stores []model.Store
	var total int64

	countQuery := r.db.Model(&model.Store{})
	r.applyFilters(countQuery, req)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Model(&model.Store{})
	r.applyFilters(query, req)

	err := query.Scopes(util.Paginate(req.Page, req.Limit)).
		Order(req.SortBy + " " + req.OrderBy).
		Preload("User").
		Find(&stores).Error

	if err != nil {
		return nil, 0, err
	}

	storesResponse := make([]dto.StoreResponse, len(stores))
	for i, store := range stores {
		deletedAt := ""
		if store.DeletedAt.Valid {
			deletedAt = store.DeletedAt.Time.Format(time.RFC3339)
		}
		storesResponse[i] = dto.StoreResponse{
			ID:          store.ID,
			Name:        store.Name,
			Description: store.Description,
			IsActive:    store.IsActive,
			UserID:      store.UserID,
			Username:    store.User.Username,
			CreatedAt:   store.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  store.UpdatedAt.Format(time.RFC3339),
			DeletedAt:  deletedAt,
		}
	}

	return storesResponse, total, nil
}

func (r *storeQueryRepository) applyFilters(query *gorm.DB, req dto.QueryStoreParam) *gorm.DB {
	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	return query
}

func (r *storeQueryRepository) FindByIDWithLock(tx *gorm.DB, id string) (*model.Store, error) {
	var store model.Store

	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("User").
		First(&store, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &store, nil
}
