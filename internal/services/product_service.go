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

type ProductService interface {
	CreateProduct(product dtos.CreateProductRequest, user dtos.CurrentUser) (dtos.ProductResponse, error)
	GetPagination(req dtos.QueryProductRequest, user dtos.CurrentUser) ([]dtos.ProductResponse, int64, error)
	GetProductBySlug(slug string) (*models.Product, error)
	UpdateProduct(id string, req dtos.UpdateProductRequest, user dtos.CurrentUser) (dtos.ProductResponse, error)
	DeleteProduct(id string, user dtos.CurrentUser) (dtos.ProductResponse, error)
	ActivateProduct(req dtos.ActivateProductRequest, user dtos.CurrentUser) (dtos.ProductResponse, error)
}

type productService struct {
	db               *gorm.DB
	productRepo      repositories.ProductRepository
	productQueryRepo repositories.ProductQueryRepository
	storeQueryRepo   repositories.StoreQueryRepository
}

func NewProductService(
	db *gorm.DB,
	productRepo repositories.ProductRepository,
	productQueryRepo repositories.ProductQueryRepository,
	storeQueryRepo repositories.StoreQueryRepository,
) ProductService {
	return &productService{
		db,
		productRepo,
		productQueryRepo,
		storeQueryRepo,
	}
}

func (s *productService) CreateProduct(
	req dtos.CreateProductRequest,
	user dtos.CurrentUser,
) (dtos.ProductResponse, error) {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	userStore, err := s.storeQueryRepo.FindByUserID(user.ID)
	if err != nil {
		tx.Rollback()
		return dtos.ProductResponse{},
			utils.ForbiddenException("User " + user.Username + " doesn't have a store, can't create product")
	}

	product := models.Product{
		StoreID:     userStore.ID,
		Name:        req.Name,
		Description: req.Description,
		Price:       *req.Price,
		Stock:       *req.Stock,
		Unit:        req.Unit,
		Base: models.Base{
			CreatedBy: uuid.MustParse(user.ID),
			UpdatedBy: uuid.MustParse(user.ID),
		},
	}

	if err := s.productRepo.Create(tx, &product); err != nil {
		tx.Rollback()
		return dtos.ProductResponse{}, utils.InternalServerErrorException("Failed to save product")
	}

	if err := tx.Commit().Error; err != nil {
		return dtos.ProductResponse{}, err
	}

	return dtos.ProductResponse{
		ID:        product.ID.String(),
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		Unit:      product.Unit,
		IsActive:  product.IsActive,
		CreatedAt: product.CreatedAt.Format(time.RFC3339),
		CreatedBy: product.CreatedBy.String(),
		UpdatedAt: product.UpdatedAt.Format(time.RFC3339),
		UpdatedBy: product.UpdatedBy.String(),
		DeletedAt: product.DeletedAt.Time.Format(time.RFC3339),
		StoreID:   product.StoreID.String(),
		StoreName: userStore.Name,
	}, nil
}

func (s *productService) GetPagination(req dtos.QueryProductRequest, user dtos.CurrentUser) ([]dtos.ProductResponse, int64, error) {
	if user.Role == constants.Seller {
		userStore, err := s.storeQueryRepo.FindByUserID(user.ID)
		if err != nil {
			return nil, 0, err
		}

		if userStore == nil {
			return []dtos.ProductResponse{}, 0, nil
		}

		strID := userStore.ID.String()
		req.StoreID = &strID
	}

	if req.IsActive == nil {
		req.IsActive = new(bool)
		*req.IsActive = true
	}

	return s.productQueryRepo.FindAllPagination(req, user)
}

func (s *productService) GetProductBySlug(slug string) (*models.Product, error) {
	return s.productQueryRepo.FindBySlug(slug)
}

func (s *productService) UpdateProduct(
	id string,
	req dtos.UpdateProductRequest,
	user dtos.CurrentUser,
) (dtos.ProductResponse, error) {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	product, err := s.productQueryRepo.FindByIDPreloadStore(id)
	if err != nil {
		tx.Rollback()
		return dtos.ProductResponse{}, utils.NotFoundException("Product not found")
	}

	userStore, err := s.storeQueryRepo.FindByUserID(user.ID)
	if err != nil {
		tx.Rollback()
		return dtos.ProductResponse{}, utils.ForbiddenException("User doesn't have a store")
	}

	if product.StoreID != userStore.ID {
		tx.Rollback()
		return dtos.ProductResponse{}, utils.ForbiddenException("You don't have permission to update this product")
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.Unit != nil {
		product.Unit = *req.Unit
	}

	product.UpdatedBy = uuid.MustParse(user.ID)

	if err := s.productRepo.Update(tx, product); err != nil {
		tx.Rollback()
		return dtos.ProductResponse{}, utils.InternalServerErrorException("Failed to update product")
	}

	if err := tx.Commit().Error; err != nil {
		return dtos.ProductResponse{}, err
	}

	return dtos.ProductResponse{
		ID:        product.ID.String(),
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		Unit:      product.Unit,
		IsActive:  product.IsActive,
		CreatedAt: product.CreatedAt.Format(time.RFC3339),
		CreatedBy: product.CreatedBy.String(),
		UpdatedAt: product.UpdatedAt.Format(time.RFC3339),
		UpdatedBy: product.UpdatedBy.String(),
		DeletedAt: product.DeletedAt.Time.Format(time.RFC3339),
		StoreID:   product.StoreID.String(),
		StoreName: product.Store.Name,
	}, nil
}

func (s *productService) DeleteProduct(
	id string,
	user dtos.CurrentUser,
) (dtos.ProductResponse, error) {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	product, err := s.productQueryRepo.FindByIDPreloadStore(id)
	if err != nil {
		tx.Rollback()
		return dtos.ProductResponse{},
			utils.NotFoundException("Product not found")
	}

	userStore, err := s.storeQueryRepo.FindByUserID(user.ID)
	if err != nil {
		tx.Rollback()
		return dtos.ProductResponse{},
			utils.ForbiddenException("User doesn't have a store")
	}

	if product.StoreID != userStore.ID {
		tx.Rollback()
		return dtos.ProductResponse{},
			utils.ForbiddenException("You don't have permission to delete this product")
	}

	product.UpdatedBy = uuid.MustParse(user.ID)
	product.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}

	if err := s.productRepo.Delete(tx, product); err != nil {
		tx.Rollback()
		return dtos.ProductResponse{},
			utils.InternalServerErrorException("Failed to delete product")
	}

	if err := tx.Commit().Error; err != nil {
		return dtos.ProductResponse{}, err
	}

	return dtos.ProductResponse{
		ID:        product.ID.String(),
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		Unit:      product.Unit,
		IsActive:  product.IsActive,
		CreatedAt: product.CreatedAt.Format(time.RFC3339),
		CreatedBy: product.CreatedBy.String(),
		UpdatedAt: product.UpdatedAt.Format(time.RFC3339),
		UpdatedBy: product.UpdatedBy.String(),
		DeletedAt: product.DeletedAt.Time.Format(time.RFC3339),
		StoreID:   product.StoreID.String(),
		StoreName: product.Store.Name,
	}, nil
}

func (s *productService) ActivateProduct(
	req dtos.ActivateProductRequest,
	user dtos.CurrentUser,
) (dtos.ProductResponse, error) {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	product, err := s.productQueryRepo.FindByIDPreloadStore(req.ID)
	if err != nil {
		tx.Rollback()
		return dtos.ProductResponse{},
			utils.NotFoundException("Product not found")
	}

	if product.IsActive == *req.IsActive {
		tx.Rollback()
		status := "active"
		if !*req.IsActive {
			status = "inactive"
		}
		return dtos.ProductResponse{},
			utils.BadRequestException("Product is already "+status, nil)
	}

	userStore, err := s.storeQueryRepo.FindByUserID(user.ID)
	if err != nil {
		tx.Rollback()
		return dtos.ProductResponse{},
			utils.ForbiddenException("User doesn't have a store")
	}

	if product.StoreID != userStore.ID {
		tx.Rollback()
		return dtos.ProductResponse{},
			utils.ForbiddenException("You don't have permission to activate/deactivate this product")
	}

	product.UpdatedBy = uuid.MustParse(user.ID)
	product.IsActive = *req.IsActive

	if err := s.productRepo.Update(tx, product); err != nil {
		tx.Rollback()
		return dtos.ProductResponse{},
			utils.InternalServerErrorException("Failed to activate/deactivate product")
	}

	if err := tx.Commit().Error; err != nil {
		return dtos.ProductResponse{}, err
	}

	return dtos.ProductResponse{
		ID:        product.ID.String(),
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		Unit:      product.Unit,
		IsActive:  product.IsActive,
		CreatedAt: product.CreatedAt.Format(time.RFC3339),
		CreatedBy: product.CreatedBy.String(),
		UpdatedAt: product.UpdatedAt.Format(time.RFC3339),
		UpdatedBy: product.UpdatedBy.String(),
		DeletedAt: product.DeletedAt.Time.Format(time.RFC3339),
		StoreID:   product.StoreID.String(),
		StoreName: product.Store.Name,
	}, nil
}
