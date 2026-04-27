package repositories

import (
	"e-shop-api/internal/dtos"
	"e-shop-api/internal/models"
	"e-shop-api/internal/pkg/utils"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StoreQueryRepository interface {
	FindAll() ([]models.Store, error)
	FindByID(id string) (*models.Store, error)
	FindByUserID(userID string) (*models.Store, error)
	FindAllPagination(req dtos.QueryStoreParam) ([]dtos.StoreResponse, int64, error)
	FindByIDWithLock(tx *gorm.DB, id string) (*models.Store, error)
}

type storeQueryRepository struct {
	db *gorm.DB
}

func NewStoreQueryRepository(db *gorm.DB) StoreQueryRepository {
	return &storeQueryRepository{db}
}

func (r *storeQueryRepository) FindAll() ([]models.Store, error) {
	var stores []models.Store
	err := r.db.Preload("User").Find(&stores).Error
	return stores, err
}

func (r *storeQueryRepository) FindByID(id string) (*models.Store, error) {
	var store models.Store
	err := r.db.Preload("User").Where("id = ?", id).First(&store).Error
	return &store, err
}

func (r *storeQueryRepository) FindByUserID(userID string) (*models.Store, error) {
	var store models.Store
	err := r.db.Where("user_id = ?", userID).First(&store).Error
	return &store, err
}

func (r *storeQueryRepository) FindAllPagination(req dtos.QueryStoreParam) ([]dtos.StoreResponse, int64, error) {
	var stores []models.Store
	var total int64

	countQuery := r.db.Model(&models.Store{})
	r.applyFilters(countQuery, req)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Model(&models.Store{})
	r.applyFilters(query, req)

	err := query.Scopes(utils.Paginate(req.Page, req.Limit)).
		Order(req.SortBy + " " + req.OrderBy).
		Preload("User").
		Find(&stores).Error

	if err != nil {
		return nil, 0, err
	}

	storesResponse := make([]dtos.StoreResponse, len(stores))
	for i, store := range stores {
		deletedAt := ""
		if store.DeletedAt.Valid {
			deletedAt = store.DeletedAt.Time.Format(time.RFC3339)
		}
		storesResponse[i] = dtos.StoreResponse{
			ID:          store.ID,
			Name:        store.Name,
			Description: store.Description,
			IsActive:    store.IsActive,
			UserID:      store.UserID,
			Username:    store.User.Username,
			CreatedAt:   store.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   store.UpdatedAt.Format(time.RFC3339),
			DeletedAt:   deletedAt,
		}
	}

	return storesResponse, total, nil
}

func (r *storeQueryRepository) applyFilters(query *gorm.DB, req dtos.QueryStoreParam) *gorm.DB {
	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	return query
}

func (r *storeQueryRepository) FindByIDWithLock(tx *gorm.DB, id string) (*models.Store, error) {
	var store models.Store

	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("User").
		First(&store, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &store, nil
}
