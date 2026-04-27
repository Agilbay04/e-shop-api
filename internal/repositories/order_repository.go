package repositories

import (
	"e-shop-api/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(tx *gorm.DB, order *models.Order) error
	CreateOrderItems(tx *gorm.DB, items []models.OrderItem) error
	DeleteOrderItems(tx *gorm.DB, orderID string) error
	UpdateOrder(tx *gorm.DB, order *models.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) CreateOrder(tx *gorm.DB, order *models.Order) error {
	if tx != nil {
		return tx.Create(order).Error
	}

	return r.db.Create(order).Error
}

func (r *orderRepository) CreateOrderItems(tx *gorm.DB, items []models.OrderItem) error {
	if tx != nil {
		return tx.Create(&items).Error
	}

	return r.db.Create(&items).Error
}

func (r *orderRepository) DeleteOrderItems(tx *gorm.DB, orderID string) error {
	if tx != nil {
		return tx.Where("order_id = ?", orderID).Delete(&models.OrderItem{}).Error
	}
	return r.db.Where("order_id = ?", orderID).Delete(&models.OrderItem{}).Error
}

func (r *orderRepository) UpdateOrder(tx *gorm.DB, order *models.Order) error {
	if tx != nil {
		return tx.Save(order).Error
	}
	return r.db.Save(order).Error
}
