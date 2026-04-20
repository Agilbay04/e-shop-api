package handler

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/pkg/util"
	"e-shop-api/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	ProductService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		ProductService: productService,
	}
}

func (ph *ProductHandler) Index(ctx *gin.Context) {
	var req dto.QueryProductRequest
	user := util.GetCurrentUser(ctx)

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(util.BadRequestException("Invalid query parameters", err))
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
		ctx.Error(util.BadRequestException("Invalid request body", err))
		return
	}

	user := util.GetCurrentUser(ctx)
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
		ctx.Error(util.BadRequestException("Product ID is required", nil))
		return
	}

	var req dto.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(util.BadRequestException("Invalid request body", err))
		return
	}

	user := util.GetCurrentUser(ctx)
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
		ctx.Error(util.BadRequestException("Product ID is required", nil))
		return
	}

	user := util.GetCurrentUser(ctx)
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
		ctx.Error(util.BadRequestException("Invalid request body", err))
		return
	}

	user := util.GetCurrentUser(ctx)
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
