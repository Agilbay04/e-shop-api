package dtos

import "e-shop-api/internal/constants"

type CurrentUser struct {
    ID       string              `json:"id"`
    Username string              `json:"username"`
    Email    string              `json:"email"`
    Role     constants.UserRole  `json:"role"`
    Picture  string              `json:"picture"`
}