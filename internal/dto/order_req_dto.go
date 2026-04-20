package dto

import (
	"e-shop-api/internal/model"

	"github.com/google/uuid"
)

type OrderRequest struct {
    IsCheckout bool                 `json:"is_checkout"`
    GrandTotal int                  `json:"grand_total"`
    Status     model.OrderStatus    `json:"status"`
    OrderItems []OrderItemRequest   `json:"order_items" binding:"required"`
}

type OrderItemRequest struct {
    StoreID   uuid.UUID `json:"store_id" binding:"required"`
    ProductID uuid.UUID `json:"product_id" binding:"required"`
    Quantity  int       `json:"quantity" binding:"required,min=1"`
    Price     int       `json:"price" binding:"required,min=0"`
    SubTotal  int       `json:"sub_total" binding:"required,min=0"`
}