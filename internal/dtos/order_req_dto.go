package dtos

import (
	"e-shop-api/internal/constants"

	"github.com/google/uuid"
)

type OrderRequest struct {
	IsCheckout bool               	 `json:"is_checkout"`
	GrandTotal int                	 `json:"grand_total"`
	Status     constants.OrderStatus `json:"status"`
	OrderItems []OrderItemRequest 	 `json:"order_items" binding:"required"`
}

type OrderItemRequest struct {
	StoreID   uuid.UUID `json:"store_id" binding:"required"`
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
	Price     int       `json:"price" binding:"required,min=0"`
	SubTotal  int       `json:"sub_total" binding:"required,min=0"`
}

func QueryOrderRequest() QueryOrderParam {
	return QueryOrderParam{
		PaginationParam: PaginationParam{
			Page:    1,
			Limit:   10,
			SortBy:  "created_at",
			OrderBy: "desc",
		},
	}
}

type QueryOrderParam struct {
	PaginationParam
	Status *constants.OrderStatus `form:"status"`
}
