package handler

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/pkg/util"
	"e-shop-api/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{us}
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	user := util.GetCurrentUser(ctx)

	res, err := h.userService.Profile(user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Success get profile")
}

func (h *UserHandler) UploadPicture(ctx *gin.Context) {
	var req dto.UploadPictureRequest
	user := util.GetCurrentUser(ctx)

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(util.BadRequestException("Invalid request body", err))
		return
	}

	res, err := h.userService.UploadPicture(req, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Success update profile")
}