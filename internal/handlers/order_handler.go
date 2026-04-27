package handlers

import (
	"e-shop-api/internal/dtos"
	"e-shop-api/internal/pkg/utils"
	"e-shop-api/internal/services"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

var orderAllowedSortBy = map[string]bool{
	"created_at":  true,
	"updated_at":  true,
	"grand_total": true,
	"status":     true,
}

func (h *OrderHandler) Index(ctx *gin.Context) {
	var req dto.QueryOrderParam
	user := utils.GetCurrentUser(ctx)

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(utils.BadRequestException("Invalid query parameters", err))
		return
	}

	if !dto.ValidateSortByPattern(req.SortBy) || !dto.IsAllowedSortBy(req.SortBy, orderAllowedSortBy) {
		ctx.Error(utils.BadRequestException("Invalid sort_by value. Allowed: created_at, updated_at, grand_total, status", nil))
		return
	}

	orders, total, err := h.orderService.GetOrders(req, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	OkPagination(ctx, orders, total, req.PaginationParam, "Success retrieve data")
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	req := dto.OrderRequest{
		IsCheckout: true,
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.BadRequestException("Invalid request body", err))
		return
	}

	user := utils.GetCurrentUser(ctx)
	res, err := h.orderService.CreateOrder(req, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Created(ctx, res, "Success create order")
}

func (h *OrderHandler) UpdateOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(utils.BadRequestException("Order ID is required", nil))
		return
	}

	var req dto.OrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.BadRequestException("Invalid request body", err))
		return
	}

	user := utils.GetCurrentUser(ctx)
	res, err := h.orderService.UpdateOrder(id, req, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Success update order")
}

func (h *OrderHandler) CancelOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(utils.BadRequestException("Order ID is required", nil))
		return
	}

	user := utils.GetCurrentUser(ctx)
	res, err := h.orderService.CancelOrder(id, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Order has been canceled")
}

func (h *OrderHandler) ConfirmOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(utils.BadRequestException("Order ID is required", nil))
		return
	}

	user := utils.GetCurrentUser(ctx)
	res, err := h.orderService.ConfirmOrder(id, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Order has been confirmed")
}
