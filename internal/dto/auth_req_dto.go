package dto

import (
	"e-shop-api/internal/model"
	"mime/multipart"
)

type RegisterRequest struct {
	Username string 		`json:"username" binding:"required,min=3"`
	Email    string 		`json:"email" binding:"required,email"`
	Password string 		`json:"password" binding:"required,min=6"`
	Role     model.UserRole `json:"role" binding:"required,oneof=admin seller buyer"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UploadPictureRequest struct {
	Picture *multipart.FileHeader `form:"picture" binding:"required"`
}

type ForgotPasswordRequest struct {
    Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
    Token           string `json:"token" binding:"required"`
    NewPassword     string `json:"new_password" binding:"required,min=6"`
    ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

