package service

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/model"
	"e-shop-api/internal/repository"
	"errors"
)

type StoreService interface {
	CreateStore(req *dto.CreateStoreRequest, user *dto.CurrentUser) (dto.CreateStoreResponse, error)
}

type storeService struct {
	storeRepo repository.StoreRepository
	storeQueryRepo repository.StoreQueryRepository
	userQueryRepo repository.UserQueryRepository
}

func NewStoreService(
	storeRepo repository.StoreRepository, 
	storeQueryRepo repository.StoreQueryRepository, 
	userQueryRepo repository.UserQueryRepository,
) StoreService {
	return &storeService {
		storeRepo, 
		storeQueryRepo,
		userQueryRepo,
	}
}

func (s *storeService) CreateStore(
	req *dto.CreateStoreRequest, 
	user *dto.CurrentUser,
) (dto.CreateStoreResponse, error) {
    if user.Role != "seller" {
        return dto.CreateStoreResponse{}, errors.New("User is not a seller")
    }

    existingStore, err := s.storeQueryRepo.FindByUserID(req.UserID)
    
    if err == nil && existingStore != nil {
        return dto.CreateStoreResponse{}, errors.New("User " + user.Username + " already has a store")
    }

	store := &model.Store {
		Name: req.Name,
		Description: req.Description,
		UserID: req.UserID,
	}

	if err := s.storeRepo.Create(store); err != nil {
		return dto.CreateStoreResponse{}, err
	}

    return dto.CreateStoreResponse{
		ID: store.ID,
		Name: store.Name,
		Description: store.Description,
		UserID: store.UserID,
	}, nil
}