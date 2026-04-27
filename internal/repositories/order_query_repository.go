package repositories

import (
	"e-shop-api/internal/constants"
	"e-shop-api/internal/dtos"
	"e-shop-api/internal/models"
	"e-shop-api/internal/pkg/utils"

	"gorm.io/gorm"
)

type OrderQueryRepository interface {
	CountOrderItemsByStoreAndOrderStatus(tx *gorm.DB, storeID string, statuses []constants.OrderStatus) (int64, error)
	FindByIDWithLock(tx *gorm.DB, orderID string) (*models.Order, error)
	FindAllPagination(userID string, storeID string, statuses []constants.OrderStatus, req dtos.QueryOrderParam) ([]dtos.OrderResponse, int64, error)
}

type orderQueryRepository struct {
	db *gorm.DB
}

func NewOrderQueryRepository(db *gorm.DB) OrderQueryRepository {
	return &orderQueryRepository{db}
}

func (r *orderQueryRepository) CountOrderItemsByStoreAndOrderStatus(tx *gorm.DB, storeID string, statuses []constants.OrderStatus) (int64, error) {
	var count int64

	query := r.db.Model(&models.OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("order_items.store_id = ?", storeID).
		Where("orders.status IN ?", statuses)

	if tx != nil {
		query = tx.Model(&models.OrderItem{}).
			Joins("JOIN orders ON orders.id = order_items.order_id").
			Where("order_items.store_id = ?", storeID).
			Where("orders.status IN ?", statuses)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *orderQueryRepository) FindByIDWithLock(tx *gorm.DB, orderID string) (*models.Order, error) {
	var order models.Order
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

func (r *orderQueryRepository) FindAllPagination(userID string, storeID string, statuses []constants.OrderStatus, req dtos.QueryOrderParam) ([]dtos.OrderResponse, int64, error) {
	var orders []models.Order
	var total int64

	countQuery := r.db.Model(&models.Order{})
	r.applyFilters(countQuery, userID, storeID, statuses)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Model(&models.Order{})
	r.applyFilters(query, userID, storeID, statuses)

	err := query.Scopes(utils.Paginate(req.Page, req.Limit)).
		Order(req.SortBy + " " + req.OrderBy).
		Preload("User").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Preload("OrderItems.Product.Store").
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	responses := make([]dtos.OrderResponse, len(orders))
	for i, order := range orders {
		orderItems := make([]dtos.OrderItemResponse, len(order.OrderItems))
		for j, item := range order.OrderItems {
			orderItems[j] = dtos.OrderItemResponse{
				StoreID:     item.StoreID.String(),
				StoreName:   item.Product.Store.Name,
				ProductID:   item.ProductID.String(),
				ProductName: item.Product.Name,
				Quantity:    item.Quantity,
				Unit:        item.Product.Unit,
				Price:       item.Price,
				SubTotal:    item.SubTotal,
			}
		}

		responses[i] = dtos.OrderResponse{
			ID:          order.ID.String(),
			UserID:      order.UserID.String(),
			Username:    order.User.Username,
			GrandTotal:  order.GrandTotal,
			Status:      order.Status,
			OrderItems:  orderItems,
			OrderNumber: order.OrderNumber,
		}
	}

	return responses, total, nil
}

func (r *orderQueryRepository) applyFilters(query *gorm.DB, userID string, storeID string, statuses []constants.OrderStatus) *gorm.DB {
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if storeID != "" {
		query = query.Joins("JOIN order_items ON order_items.order_id = orders.id").
			Where("order_items.store_id = ?", storeID)
	}
	if len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
	}
	return query
}
