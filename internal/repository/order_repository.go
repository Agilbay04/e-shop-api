package repository

import (
	"e-shop-api/internal/model"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(tx *gorm.DB, order *model.Order) error
	CreateOrderItems(tx *gorm.DB, items []model.OrderItem) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository {db}
}

func (r *orderRepository) CreateOrder(tx *gorm.DB, order *model.Order) error {
	if tx != nil {
		return tx.Create(order).Error
	}

	return r.db.Create(order).Error
}

func (r *orderRepository) CreateOrderItems(tx *gorm.DB, items []model.OrderItem) error {
	if tx != nil {
		return tx.Create(&items).Error
	}
	
	return r.db.Create(&items).Error
}