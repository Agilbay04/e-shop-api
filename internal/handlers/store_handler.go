package handlers

import (
	"e-shop-api/internal/dtos"
	"e-shop-api/internal/pkg/utils"
	"e-shop-api/internal/services"

	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	storeService services.StoreService
}

func NewStoreHandler(ss services.StoreService) *StoreHandler {
	return &StoreHandler{ss}
}

var storeAllowedSortBy = map[string]bool{
	"created_at": true,
	"updated_at": true,
	"name":       true,
}

func (h *StoreHandler) Index(ctx *gin.Context) {
	var req dto.QueryStoreParam
	user := utils.GetCurrentUser(ctx)

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(utils.BadRequestException("Invalid query parameters", err))
		return
	}

	if !dto.ValidateSortByPattern(req.SortBy) || !dto.IsAllowedSortBy(req.SortBy, storeAllowedSortBy) {
		ctx.Error(utils.BadRequestException("Invalid sort_by value. Allowed: created_at, updated_at, name", nil))
		return
	}

	stores, total, err := h.storeService.GetStores(req, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	OkPagination(ctx, stores, total, req.PaginationParam, "Success retrieve data")
}

func (h *StoreHandler) CreateStore(ctx *gin.Context) {
	user := utils.GetCurrentUser(ctx)
	var req dto.CreateStoreRequest
	req.UserID = user.ID

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.BadRequestException("Invalid request body", err.Error()))
		return
	}

	res, err := h.storeService.CreateStore(req, user)
	if err != nil {
		ctx.Error(utils.UnprocessableEntityException(err.Error()))
		return
	}

	Created(ctx, res, "Success create data")
}

func (h *StoreHandler) UpdateStore(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(utils.BadRequestException("Store ID is required", nil))
		return
	}

	var req dto.UpdateStoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.BadRequestException("Invalid request body", err))
		return
	}

	user := utils.GetCurrentUser(ctx)
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
		ctx.Error(utils.BadRequestException("Invalid request body", err))
		return
	}

	user := utils.GetCurrentUser(ctx)
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
		ctx.Error(utils.BadRequestException("Store ID is required", nil))
		return
	}

	user := utils.GetCurrentUser(ctx)
	res, err := h.storeService.DeleteStore(id, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Success delete data")
}
