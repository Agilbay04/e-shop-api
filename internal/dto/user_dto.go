package dto

import "e-shop-api/internal/model"

type CurrentUser struct {
    ID       string         `json:"id"`
    Username string         `json:"username"`
    Email    string         `json:"email"`
    Role     model.UserRole `json:"role"`
    Picture  string         `json:"picture"`
}