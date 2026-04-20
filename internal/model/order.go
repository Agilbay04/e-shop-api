package model

import "github.com/google/uuid"

type OrderStatus string

const (
	Draft 		OrderStatus = "draft"
	Pending 	OrderStatus = "pending"
	Paid    	OrderStatus = "paid"
	Cancelled 	OrderStatus = "cancelled"
)

type Order struct {
	Base
	UserID      uuid.UUID   `gorm:"type:uuid;column:user_id" json:"user_id"`
	GrandTotal  int     	`gorm:"type:int;not null;column:grand_total" json:"grand_total"`
	Status      OrderStatus	`gorm:"type:varchar(20);column:status;default:draft" json:"status"`
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID"`
	User        User        `gorm:"foreignKey:UserID" json:"user"`
}

func (s OrderStatus) IsValid() bool {
	return s == Draft || s == Pending || s == Paid || s == Cancelled
}

func (Order) TableName() string {
	return "orders"
}