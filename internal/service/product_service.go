package service

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/model"
	"e-shop-api/internal/pkg/util"
	"e-shop-api/internal/repository"
)

type ProductService interface {
	CreateProduct(product dto.CreateProductRequest, user dto.CurrentUser) (dto.CreateProductResponse, error)
	GetAllProducts() ([]model.Product, error)
	GetProductBySlug(slug string) (*model.Product, error)
}

type productService struct {
	productRepo      repository.ProductRepository
	productQueryRepo repository.ProductQueryRepository
	storeQueryRepo   repository.StoreQueryRepository
}

func NewProductService(
	productRepo repository.ProductRepository, 
	productQueryRepo repository.ProductQueryRepository,
	storeQueryRepo repository.StoreQueryRepository,
) ProductService {
	return &productService {
		productRepo,
		productQueryRepo,
		storeQueryRepo,
	}
}

func (s *productService) CreateProduct(
	req dto.CreateProductRequest, 
	user dto.CurrentUser,
) (dto.CreateProductResponse, error) {
	userStore, err := s.storeQueryRepo.FindByUserID(user.ID)
	if err != nil {
		return dto.CreateProductResponse{}, 
			util.ForbiddenException("User " + user.Username + " doesn't have a store, can't create product")
	}

	product := model.Product {
		StoreID: 		userStore.ID,
		Name:    		req.Name,
		Description: 	req.Description,
		Price:   		*req.Price,
		Stock:   		*req.Stock,
		Unit:    		req.Unit,
		Base: model.Base {
			CreatedBy: user.ID,
			UpdatedBy: user.ID,
		},
	}

	if err := s.productRepo.Create(&product); err != nil {
		return dto.CreateProductResponse{}, util.InternalServerErrorException("Failed to save product")
	}

	return dto.CreateProductResponse{
		ID:    product.ID.String(),
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
		Unit:  product.Unit,
		StoreID: product.StoreID.String(),
		StoreName: userStore.Name,
	}, nil
}

func (s *productService) GetAllProducts() ([]model.Product, error) {
	return s.productQueryRepo.FindAll()
}

func (s *productService) GetProductBySlug(slug string) (*model.Product, error) {
	return s.productQueryRepo.FindBySlug(slug)
}