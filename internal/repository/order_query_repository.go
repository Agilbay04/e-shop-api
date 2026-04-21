package repository

import (
	"e-shop-api/internal/model"

	"gorm.io/gorm"
)

type OrderQueryRepository interface {
	CountOrderItemsByStoreAndOrderStatus(tx *gorm.DB, storeID string, statuses []model.OrderStatus) (int64, error)
	FindByIDWithLock(tx *gorm.DB, orderID string) (*model.Order, error)
}

type orderQueryRepository struct {
	db *gorm.DB
}

func NewOrderQueryRepository(db *gorm.DB) OrderQueryRepository {
	return &orderQueryRepository{db}
}

func (r *orderQueryRepository) CountOrderItemsByStoreAndOrderStatus(tx *gorm.DB, storeID string, statuses []model.OrderStatus) (int64, error) {
	var count int64

	query := r.db.Model(&model.OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("order_items.store_id = ?", storeID).
		Where("orders.status IN ?", statuses)

	if tx != nil {
		query = tx.Model(&model.OrderItem{}).
			Joins("JOIN orders ON orders.id = order_items.order_id").
			Where("order_items.store_id = ?", storeID).
			Where("orders.status IN ?", statuses)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *orderQueryRepository) FindByIDWithLock(tx *gorm.DB, orderID string) (*model.Order, error) {
    var order model.Order
    err := tx.Set("gorm:query_option", "FOR UPDATE").
		Preload("User").
        Preload("OrderItems").
		Preload("OrderItems.Product").
		Preload("OrderItems.Product.Store").
        First(&order, "id = ?", orderID).Error
    
    if err != nil {
        return nil, err
    }
    return &order, nil
}

