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

	Created(ctx, res, "Success create data")

}