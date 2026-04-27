package models

import "github.com/google/uuid"

type OrderItem struct{
	Base
	OrderID   uuid.UUID `gorm:"type:uuid;column:order_id" json:"order_id"`
	StoreID   uuid.UUID `gorm:"type:uuid;column:store_id" json:"store_id"`
	ProductID uuid.UUID `gorm:"type:uuid;column:product_id" json:"product_id"`
	Quantity  int       `gorm:"type:int;not null;column:quantity" json:"quantity"`
	Price     int       `gorm:"type:int;not null;column:price" json:"price"`
	SubTotal  int       `gorm:"type:int;not null;column:sub_total" json:"sub_total"`
	Order     Order     `gorm:"foreignKey:OrderID"`
	Store     Store     `gorm:"foreignKey:StoreID"`
	Product   Product   `gorm:"foreignKey:ProductID"`
}

func (OrderItem) TableName() string {
	return "order_items"
}