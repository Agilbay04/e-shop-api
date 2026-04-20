package service

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/model"
	"e-shop-api/internal/repository"
	"errors"

	"gorm.io/gorm"
)

type StoreService interface {
	CreateStore(req dto.CreateStoreRequest, user dto.CurrentUser) (dto.CreateStoreResponse, error)
}

type storeService struct {
	db             *gorm.DB
	storeRepo      repository.StoreRepository
	storeQueryRepo repository.StoreQueryRepository
	userQueryRepo  repository.UserQueryRepository
}

func NewStoreService(
	db *gorm.DB,
	storeRepo repository.StoreRepository,
	storeQueryRepo repository.StoreQueryRepository,
	userQueryRepo repository.UserQueryRepository,
) StoreService {
	return &storeService{
		db,
		storeRepo,
		storeQueryRepo,
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

	if user.Role != "seller" {
		tx.Rollback()
		return dto.CreateStoreResponse{}, errors.New("User is not a seller")
	}

	existingStore, err := s.storeQueryRepo.FindByUserID(req.UserID.String())

	if err == nil && existingStore != nil {
		tx.Rollback()
		return dto.CreateStoreResponse{}, errors.New("User " + user.Username + " already has a store")
	}

	store := &model.Store{
		Name:        req.Name,
		Description: req.Description,
		UserID:      req.UserID,
		Base: model.Base{
			CreatedBy: user.ID,
			UpdatedBy: user.ID,
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
