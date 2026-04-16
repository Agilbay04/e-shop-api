package handler

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/pkg/util"
	"e-shop-api/internal/service"

	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	storeService service.StoreService
}

func NewStoreHandler(ss service.StoreService) *StoreHandler {
	return &StoreHandler{ss}
}

func (h *StoreHandler) Create(ctx *gin.Context) {
	user := util.GetCurrentUser(ctx)
	var req dto.CreateStoreRequest
	req.UserID = user.ID

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(util.BadRequestException("Invalid request body", err.Error()))
		return
	}

	res, err := h.storeService.CreateStore(&req, &user)
	if err != nil {
		ctx.Error(util.UnprocessableEntityException(err.Error()))
		return
	}

	Success(ctx, res, "Success create data")
}