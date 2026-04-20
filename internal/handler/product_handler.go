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
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(util.BadRequestException("Invalid query parameters", err))
		return
	}

	products, total, err := ph.ProductService.GetPagination(req)
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