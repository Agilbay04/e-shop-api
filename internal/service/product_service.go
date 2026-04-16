package service

import (
	"errors"
	"e-shop-api/internal/model"
	"e-shop-api/internal/repository"
)

type ProductService interface {
	CreateProduct(product *model.Product) error
	GetAllProducts() ([]model.Product, error)
	GetProductBySlug(slug string) (*model.Product, error)
}

type productService struct {
	productRepo      repository.ProductRepository
	productQueryRepo repository.ProductQueryRepository
}

func NewProductService(
	productRepo repository.ProductRepository, 
	productQueryRepo repository.ProductQueryRepository,
) ProductService {
	return &productService {
		productRepo,
		productQueryRepo,
	}
}

func (s *productService) CreateProduct(product *model.Product) error {
	if product.Price <= 0 {
		return errors.New("harga produk harus lebih dari 0")
	}
	return s.productRepo.Create(product)
}

func (s *productService) GetAllProducts() ([]model.Product, error) {
	return s.productQueryRepo.FindAll()
}

func (s *productService) GetProductBySlug(slug string) (*model.Product, error) {
	return s.productQueryRepo.FindBySlug(slug)
}