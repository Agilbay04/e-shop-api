package repository

import (
	"e-shop-api/internal/model"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(tx *gorm.DB, order *model.Order) error
	CreateOrderItems(tx *gorm.DB, items []model.OrderItem) error
	DeleteOrderItems(tx *gorm.DB, orderID string) error
	UpdateOrder(tx *gorm.DB, order *model.Order) error
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

func (r *orderRepository) DeleteOrderItems(tx *gorm.DB, orderID string) error {
	if tx != nil {
		return tx.Where("order_id = ?", orderID).Delete(&model.OrderItem{}).Error
	}
	return r.db.Where("order_id = ?", orderID).Delete(&model.OrderItem{}).Error
}

func (r *orderRepository) UpdateOrder(tx *gorm.DB, order *model.Order) error {
	if tx != nil {
		return tx.Save(order).Error
	}
	return r.db.Save(order).Error
}
