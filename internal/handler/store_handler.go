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

func (h *StoreHandler) CreateStore(ctx *gin.Context) {
	user := util.GetCurrentUser(ctx)
	var req dto.CreateStoreRequest
	req.UserID = user.ID

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(util.BadRequestException("Invalid request body", err.Error()))
		return
	}

	res, err := h.storeService.CreateStore(req, user)
	if err != nil {
		ctx.Error(util.UnprocessableEntityException(err.Error()))
		return
	}

	Created(ctx, res, "Success create data")
}

func (h *StoreHandler) GetStores(ctx *gin.Context) {
	var req dto.QueryStoreParam
	user := util.GetCurrentUser(ctx)

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(util.BadRequestException("Invalid query parameters", err))
		return
	}

	stores, total, err := h.storeService.GetStores(req, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	OkPagination(ctx, stores, total, req.PaginationParam, "Success retrieve data")
}

func (h *StoreHandler) UpdateStore(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(util.BadRequestException("Store ID is required", nil))
		return
	}

	var req dto.UpdateStoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(util.BadRequestException("Invalid request body", err))
		return
	}

	user := util.GetCurrentUser(ctx)
	res, err := h.storeService.UpdateStore(id, req, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Success update data")
}

func (h *StoreHandler) ActivateStore(ctx *gin.Context) {
	var req dto.ActivateStoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(util.BadRequestException("Invalid request body", err))
		return
	}

	user := util.GetCurrentUser(ctx)
	res, err := h.storeService.ActivateStore(req, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	var status string
	if req.IsActive {
		status = "activate"
	} else {
		status = "deactivate"
	}

	Ok(ctx, res, "Success "+status+" data")
}

func (h *StoreHandler) DeleteStore(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(util.BadRequestException("Store ID is required", nil))
		return
	}

	user := util.GetCurrentUser(ctx)
	res, err := h.storeService.DeleteStore(id, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Success delete data")
}
