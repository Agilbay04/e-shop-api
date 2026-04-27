package handlers

import (
	"e-shop-api/internal/dtos"
	"e-shop-api/internal/pkg/utils"
	"e-shop-api/internal/services"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	ProductService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{
		ProductService: productService,
	}
}

var productAllowedSortBy = map[string]bool{
	"created_at": 	true,
	"updated_at": 	true,
	"price":     	true,
	"name":     	true,
	"stock":     	true,
}

func (ph *ProductHandler) Index(ctx *gin.Context) {
	var req dto.QueryProductRequest
	user := utils.GetCurrentUser(ctx)

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(utils.BadRequestException("Invalid query parameters", err))
		return
	}

	if !dto.ValidateSortByPattern(req.SortBy) || !dto.IsAllowedSortBy(req.SortBy, productAllowedSortBy) {
		ctx.Error(utils.BadRequestException("Invalid sort_by value. Allowed: created_at, updated_at, price, name, stock", nil))
		return
	}

	products, total, err := ph.ProductService.GetPagination(req, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	OkPagination(ctx, products, total, req.PaginationParam, "Success retrieve data")
}

func (ph *ProductHandler) CreateProduct(ctx *gin.Context) {
	var req dto.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.BadRequestException("Invalid request body", err))
		return
	}

	user := utils.GetCurrentUser(ctx)
	res, err := ph.ProductService.CreateProduct(req, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Created(ctx, res, "Success create data")
}

func (ph *ProductHandler) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(utils.BadRequestException("Product ID is required", nil))
		return
	}

	var req dto.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.BadRequestException("Invalid request body", err))
		return
	}

	user := utils.GetCurrentUser(ctx)
	res, err := ph.ProductService.UpdateProduct(id, req, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Success update data")
}

func (ph *ProductHandler) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(utils.BadRequestException("Product ID is required", nil))
		return
	}

	user := utils.GetCurrentUser(ctx)
	res, err := ph.ProductService.DeleteProduct(id, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Success delete data")
}

func (ph *ProductHandler) ActivateProduct(ctx *gin.Context) {
	var req dto.ActivateProductRequest

	if req.IsActive == nil {
		req.IsActive = new(bool)
		*req.IsActive = true
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.BadRequestException("Invalid request body", err))
		return
	}

	user := utils.GetCurrentUser(ctx)
	res, err := ph.ProductService.ActivateProduct(req, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	var status string
	if *req.IsActive {
		status = "activate"
	} else {
		status = "deactivate"
	}

	Ok(ctx, res, "Success "+status+" data")
}
