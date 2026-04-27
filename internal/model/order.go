package model

import (
	"e-shop-api/internal/constant"

	"github.com/google/uuid"
)

type Order struct {
	Base
	UserID      uuid.UUID   			`gorm:"type:uuid;column:user_id" json:"user_id"`
	GrandTotal  int     				`gorm:"type:int;not null;column:grand_total" json:"grand_total"`
	Status      constant.OrderStatus	`gorm:"type:varchar(20);column:status;default:draft" json:"status"`
	OrderItems  []OrderItem 			`gorm:"foreignKey:OrderID"`
	User        User        			`gorm:"foreignKey:UserID" json:"user"`
}

func (Order) TableName() string {
	return "orders"
}