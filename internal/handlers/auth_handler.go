package handlers

import (
	"e-shop-api/internal/dtos"
	"e-shop-api/internal/pkg/utils"
	"e-shop-api/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(as services.AuthService) *AuthHandler {
	return &AuthHandler{as}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req dtos.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.BadRequestException("Invalid request body", err))
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
	var req dtos.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.BadRequestException("Invalid request body", err))
		return
	}

	res, err := h.authService.Login(req)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Login success")
}

func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	var req dtos.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.BadRequestException("Invalid request body", err))
		return
	}

	res, err := h.authService.RefreshToken(req)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, res, "Token refreshed successfully")
}

func (h *AuthHandler) ForgotPassword(ctx *gin.Context) {
	var req dtos.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(utils.BadRequestException("Invalid request body", err))
		return
	}

	err := h.authService.ForgotPassword(req)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, nil, "Request reset password has been sent to your email: "+req.Email)
}

func (h *AuthHandler) ResetPassword(ctx *gin.Context) {
	var req dtos.ResetPasswordRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(utils.BadRequestException("Invalid request body", err))
		return
	}

	err := h.authService.ResetPassword(req)
	if err != nil {
		ctx.Error(err)
		return
	}

	Ok(ctx, nil, "Success reset password, try login again with new password")
}