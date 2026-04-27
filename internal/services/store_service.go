package services

import (
	"e-shop-api/internal/constants"
	"e-shop-api/internal/dtos"
	"e-shop-api/internal/models"
	"e-shop-api/internal/pkg/utils"
	"e-shop-api/internal/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoreService interface {
	CreateStore(req dto.CreateStoreRequest, user dto.CurrentUser) (dto.CreateStoreResponse, error)
	GetStores(req dto.QueryStoreParam, user dto.CurrentUser) ([]dto.StoreResponse, int64, error)
	UpdateStore(id string, req dto.UpdateStoreRequest, user dto.CurrentUser) (dto.StoreResponse, error)
	ActivateStore(req dto.ActivateStoreRequest, user dto.CurrentUser) (dto.StoreResponse, error)
	DeleteStore(id string, user dto.CurrentUser) (dto.StoreResponse, error)
}

type storeService struct {
	db             *gorm.DB
	storeRepo      repositories.StoreRepository
	storeQueryRepo repositories.StoreQueryRepository
	orderQueryRepo repositories.OrderQueryRepository
	userQueryRepo  repositories.UserQueryRepository
}

func NewStoreService(
	db 				*gorm.DB,
	storeRepo 		repositories.StoreRepository,
	storeQueryRepo 	repositories.StoreQueryRepository,
	orderQueryRepo 	repositories.OrderQueryRepository,
	userQueryRepo 	repositories.UserQueryRepository,
) StoreService {
	return &storeService{
		db,
		storeRepo,
		storeQueryRepo,
		orderQueryRepo,
		userQueryRepo,
	}
}

func (s *storeService) CreateStore(
	req dto.CreateStoreRequest,
	user dto.CurrentUser,
) (dto.CreateStoreResponse, error) {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if user.Role != constant.Seller {
		tx.Rollback()
		return dto.CreateStoreResponse{}, utils.UnauthorizedException("User is not a seller")
	}

	existingStore, err := s.storeQueryRepo.FindByUserID(req.UserID)

	if err == nil && existingStore != nil {
		tx.Rollback()
		return dto.CreateStoreResponse{}, utils.BadRequestException("User " + user.Username + " already has a store", err)
	}

	store := &model.Store{
		Name:        req.Name,
		Description: req.Description,
		UserID:      uuid.MustParse(req.UserID),
		Base: model.Base{
			CreatedBy: uuid.MustParse(user.ID),
			UpdatedBy: uuid.MustParse(user.ID),
		},
	}

	if err := s.storeRepo.Create(tx, store); err != nil {
		tx.Rollback()
		return dto.CreateStoreResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return dto.CreateStoreResponse{}, err
	}

	return dto.CreateStoreResponse{
		ID:          store.ID,
		Name:        store.Name,
		Description: store.Description,
		UserID:      store.UserID,
		Username:    user.Username,
	}, nil
}

func (s *storeService) GetStores(
	req dto.QueryStoreParam,
	user dto.CurrentUser,
) ([]dto.StoreResponse, int64, error) {
	if user.Role == constant.Seller {
		userStore, _ := s.storeQueryRepo.FindByUserID(user.ID)

		if userStore == nil {
			return []dto.StoreResponse{}, 0, nil
		}

		userID := user.ID
		req.UserID = &userID
	}

	return s.storeQueryRepo.FindAllPagination(req)
}

func (s *storeService) UpdateStore(
	id string,
	req dto.UpdateStoreRequest,
	user dto.CurrentUser,
) (dto.StoreResponse, error) {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	store, err := s.storeQueryRepo.FindByIDWithLock(tx, id)
	if err != nil {
		tx.Rollback()
		return dto.StoreResponse{}, utils.NotFoundException("Store not found")
	}

	if store.UserID.String() != user.ID {
		tx.Rollback()
		return dto.StoreResponse{}, utils.ForbiddenException("You don't have permission to update this store")
	}

	if req.Name != nil {
		store.Name = *req.Name
	}
	if req.Description != nil {
		store.Description = *req.Description
	}

	store.UpdatedBy = uuid.MustParse(user.ID)

	if err := s.storeRepo.Update(tx, store); err != nil {
		tx.Rollback()
		return dto.StoreResponse{}, utils.InternalServerErrorException("Failed to update store")
	}

	if err := tx.Commit().Error; err != nil {
		return dto.StoreResponse{}, err
	}

	deletedAt := ""
	if store.DeletedAt.Valid {
		deletedAt = store.DeletedAt.Time.Format(time.RFC3339)
	}

	return dto.StoreResponse{
		ID:          store.ID,
		Name:        store.Name,
		Description: store.Description,
		IsActive:    store.IsActive,
		UserID:      store.UserID,
		Username:    store.User.Username,
		CreatedAt:   store.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   store.UpdatedAt.Format(time.RFC3339),
		DeletedAt:   deletedAt,
	}, nil
}

func (s *storeService) ActivateStore(
	req dto.ActivateStoreRequest,
	user dto.CurrentUser,
) (dto.StoreResponse, error) {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	store, err := s.storeQueryRepo.FindByIDWithLock(tx, req.ID)
	if err != nil {
		tx.Rollback()
		return dto.StoreResponse{}, utils.NotFoundException("Store not found")
	}

	if store.UserID.String() != user.ID {
		tx.Rollback()
		return dto.StoreResponse{}, utils.ForbiddenException("You don't have permission to activate/deactivate this store")
	}

	if store.IsActive == req.IsActive {
		tx.Rollback()
		status := "active"
		if !req.IsActive {
			status = "inactive"
		}
		return dto.StoreResponse{}, utils.BadRequestException("Store is already "+status, nil)
	}

	if !req.IsActive {
		count, err := s.orderQueryRepo.CountOrderItemsByStoreAndOrderStatus(tx, req.ID, []constant.OrderStatus{constant.Pending})
		if err != nil {
			tx.Rollback()
			return dto.StoreResponse{}, utils.InternalServerErrorException("Failed to check order items")
		}
		if count > 0 {
			tx.Rollback()
			return dto.StoreResponse{}, utils.BadRequestException("Cannot deactivate store with pending order items", nil)
		}
	}

	store.UpdatedBy = uuid.MustParse(user.ID)
	store.IsActive = req.IsActive

	if err := s.storeRepo.Update(tx, store); err != nil {
		tx.Rollback()
		return dto.StoreResponse{}, utils.InternalServerErrorException("Failed to activate/deactivate store")
	}

	if err := tx.Commit().Error; err != nil {
		return dto.StoreResponse{}, err
	}

	deletedAt := ""
	if store.DeletedAt.Valid {
		deletedAt = store.DeletedAt.Time.Format(time.RFC3339)
	}

	return dto.StoreResponse{
		ID:          store.ID,
		Name:        store.Name,
		Description: store.Description,
		IsActive:    store.IsActive,
		UserID:      store.UserID,
		Username:    store.User.Username,
		CreatedAt:   store.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   store.UpdatedAt.Format(time.RFC3339),
		DeletedAt:   deletedAt,
	}, nil
}

func (s *storeService) DeleteStore(
	id string,
	user dto.CurrentUser,
) (dto.StoreResponse, error) {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	store, err := s.storeQueryRepo.FindByIDWithLock(tx, id)
	if err != nil {
		tx.Rollback()
		return dto.StoreResponse{}, utils.NotFoundException("Store not found")
	}

	if store.UserID.String() != user.ID {
		tx.Rollback()
		return dto.StoreResponse{}, utils.ForbiddenException("You don't have permission to delete this store")
	}

	count, err := s.orderQueryRepo.CountOrderItemsByStoreAndOrderStatus(tx, id, []constant.OrderStatus{constant.Draft, constant.Pending})
	if err != nil {
		tx.Rollback()
		return dto.StoreResponse{}, utils.InternalServerErrorException("Failed to check order items")
	}
	if count > 0 {
		tx.Rollback()
		return dto.StoreResponse{}, utils.BadRequestException("Cannot delete store with draft or pending order items", nil)
	}

	store.UpdatedBy = uuid.MustParse(user.ID)
	store.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}

	if err := s.storeRepo.Delete(tx, store.ID.String()); err != nil {
		tx.Rollback()
		return dto.StoreResponse{}, utils.InternalServerErrorException("Failed to delete store")
	}

	if err := tx.Commit().Error; err != nil {
		return dto.StoreResponse{}, err
	}

	deletedAt := store.DeletedAt.Time.Format(time.RFC3339)

	return dto.StoreResponse{
		ID:          store.ID,
		Name:        store.Name,
		Description: store.Description,
		IsActive:    store.IsActive,
		UserID:      store.UserID,
		Username:    store.User.Username,
		CreatedAt:   store.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   store.UpdatedAt.Format(time.RFC3339),
		DeletedAt:   deletedAt,
	}, nil
}
