package handler

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/pkg/util"
	"e-shop-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(as service.AuthService) *AuthHandler {
	return &AuthHandler{as}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(util.BadRequestException("Invalid request body", err))
		return
	}

	res, err := h.authService.Register(req)
	if err != nil {
		ctx.Error(err)
		return
	}

	Created(ctx, res, "Success create data")
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(util.BadRequestException("Invalid request body", err))
		return
	}

	res, err := h.authService.Login(req)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Login success")
}

func (h *AuthHandler) Profile(ctx *gin.Context) {
	user := util.GetCurrentUser(ctx)

	res, err := h.authService.Profile(user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Success get profile")
}

func (h *AuthHandler) UploadPicture(ctx *gin.Context) {
	var req dto.UploadPictureRequest
	user := util.GetCurrentUser(ctx)

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(util.BadRequestException("Invalid request body", err))
		return
	}

	res, err := h.authService.UploadPicture(req, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Success update profile")
}
