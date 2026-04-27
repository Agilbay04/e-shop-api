package dto

import "e-shop-api/internal/constants"

type OrderResponse struct {
	ID          	string				`json:"id"`
	OrderNumber 	string 				`json:"order_number"`
	UserID      	string 				`json:"user_id"`
	Username    	string				`json:"username"`
	GrandTotal  	int       			`json:"grand_total"`
	Status      	constant.OrderStatus 	`json:"status"`
	OrderItems 	[]OrderItemResponse 	`json:"order_items"`
}

type OrderItemResponse struct {
	StoreID   	string 	`json:"store_id"`
	StoreName 	string  `json:"store_name"`
	ProductID 	string 	`json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity  	int     `json:"quantity"`
	Unit      	string  `json:"unit"`
	Price     	int     `json:"price"`
	SubTotal  	int     `json:"sub_total"`
}