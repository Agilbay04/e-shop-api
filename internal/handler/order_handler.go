package handler

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/pkg/util"
	"e-shop-api/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler {
		orderService: orderService,
	}
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	req := dto.OrderRequest {
		IsCheckout: true,
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(util.BadRequestException("Invalid request body", err))
		return
	}

	user := util.GetCurrentUser(ctx)
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
		ctx.Error(util.BadRequestException("Order ID is required", nil))
		return
	}

	var req dto.OrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(util.BadRequestException("Invalid request body", err))
		return
	}

	user := util.GetCurrentUser(ctx)
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
		ctx.Error(util.BadRequestException("Order ID is required", nil))
		return
	}

	user := util.GetCurrentUser(ctx)
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
		ctx.Error(util.BadRequestException("Order ID is required", nil))
		return
	}

	user := util.GetCurrentUser(ctx)
	res, err := h.orderService.ConfirmOrder(id, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Order has been confirmed")
}
