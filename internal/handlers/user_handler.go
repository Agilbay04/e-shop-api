package handlers

import (
	"e-shop-api/internal/dtos"
	"e-shop-api/internal/pkg/utils"
	"e-shop-api/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(us services.UserService) *UserHandler {
	return &UserHandler{us}
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	user := utils.GetCurrentUser(ctx)

	res, err := h.userService.Profile(user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Success get profile")
}

func (h *UserHandler) UploadPicture(ctx *gin.Context) {
	var req dto.UploadPictureRequest
	user := utils.GetCurrentUser(ctx)

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(utils.BadRequestException("Invalid request body", err))
		return
	}

	res, err := h.userService.UploadPicture(req, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Success update profile")
}