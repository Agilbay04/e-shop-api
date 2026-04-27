package services

import (
	"e-shop-api/internal/constants"
	"e-shop-api/internal/dtos"
	"e-shop-api/internal/models"
	"e-shop-api/internal/pkg/utils"
	"e-shop-api/internal/repositories"
	"slices"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(req dtos.OrderRequest, user dtos.CurrentUser) (dtos.OrderResponse, error)
	UpdateOrder(orderId string, req dtos.OrderRequest, user dtos.CurrentUser) (dtos.OrderResponse, error)
	CancelOrder(orderId string, user dtos.CurrentUser) (dtos.OrderResponse, error)
	ConfirmOrder(orderId string, user dtos.CurrentUser) (dtos.OrderResponse, error)
	GetOrders(req dtos.QueryOrderParam, user dtos.CurrentUser) ([]dtos.OrderResponse, int64, error)
}

type orderService struct {
	db                *gorm.DB
	orderRepo         repositories.OrderRepository
	orderQueryRepo    repositories.OrderQueryRepository
	productRepo       repositories.ProductRepository
	productQueryRepo  repositories.ProductQueryRepository
	storeQueryRepo    repositories.StoreQueryRepository
	orderSequenceRepo repositories.OrderSequenceRepository
	notifService      NotificationService
}

func NewOrderService(
	db *gorm.DB,
	orderRepo repositories.OrderRepository,
	orderQueryRepo repositories.OrderQueryRepository,
	productRepo repositories.ProductRepository,
	productQueryRepo repositories.ProductQueryRepository,
	storeQueryRepo repositories.StoreQueryRepository,
	orderSequenceRepo repositories.OrderSequenceRepository,
	notifService NotificationService,
) OrderService {
	return &orderService{
		db,
		orderRepo,
		orderQueryRepo,
		productRepo,
		productQueryRepo,
		storeQueryRepo,
		orderSequenceRepo,
		notifService,
	}
}

func (o *orderService) CreateOrder(
	req dtos.OrderRequest,
	user dtos.CurrentUser,
) (dtos.OrderResponse, error) {
	// Set Status
	if req.IsCheckout {
		req.Status = constants.Pending
	} else {
		req.Status = constants.Draft
	}

	// Begin Transaction
	tx := o.db.Begin()

	// Generate Order Number
	orderNumber, err := o.generateOrderNumber(tx)
	if err != nil {
		tx.Rollback()
		return dtos.OrderResponse{}, err
	}

	// Safety net: Rollback if panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Preparing Data and Validation
	totalOrderPrice, orderItems, itemResponses, err := o.prepareOrderData(tx, req, user)
	if err != nil {
		tx.Rollback()
		return dtos.OrderResponse{}, err
	}

	// Save Order
	newOrder, err := o.saveOrder(tx, user.ID, totalOrderPrice, req.Status, orderNumber)
	if err != nil {
		tx.Rollback()
		return dtos.OrderResponse{}, err
	}

	// Save Bulk Order Items
	if err := o.saveOrderItems(tx, newOrder.ID, orderItems); err != nil {
		tx.Rollback()
		return dtos.OrderResponse{}, err
	}

	// Commit Transaction
	if err := tx.Commit().Error; err != nil {
		return dtos.OrderResponse{}, err
	}

	// send notification
	if req.IsCheckout {
		o.notifService.QueueSendEmail(
			user.Email,
			"Order Confirmation",
			"Your order "+newOrder.ID.String()+" has been successfully placed.",
		)
	}

	return dtos.OrderResponse{
		ID:          newOrder.ID.String(),
		OrderNumber: newOrder.OrderNumber,
		UserID:      user.ID,
		Username:    user.Username,
		GrandTotal:  totalOrderPrice,
		Status:      newOrder.Status,
		OrderItems:  itemResponses,
	}, nil
}

func (o *orderService) prepareOrderData(
	tx *gorm.DB,
	req dtos.OrderRequest,
	user dtos.CurrentUser,
) (int, []models.OrderItem, []dtos.OrderItemResponse, error) {
	var total int
	var items []models.OrderItem
	var responses []dtos.OrderItemResponse

	for _, reqItem := range req.OrderItems {
		// Locking for prevent race condition product stock
		product, err := o.productQueryRepo.FindByIDWithLock(tx, reqItem.ProductID.String())
		if err != nil {
			return 0,
				nil,
				nil,
				utils.NotFoundException("Product not found: " + reqItem.ProductID.String())
		}

		// Validation stock
		if product.Stock < reqItem.Quantity {
			return 0,
				nil,
				nil,
				utils.BadRequestException("Insufficient stock for "+product.Name, nil)
		}

		subTotal := product.Price * reqItem.Quantity
		total += subTotal

		items = append(items, models.OrderItem{
			Base:      models.Base{CreatedBy: uuid.MustParse(user.ID)},
			StoreID:   product.StoreID,
			ProductID: product.ID,
			Quantity:  reqItem.Quantity,
			Price:     product.Price,
			SubTotal:  subTotal,
		})

		responses = append(responses, dtos.OrderItemResponse{
			StoreID:     product.StoreID.String(),
			StoreName:   product.Store.Name,
			ProductID:   product.ID.String(),
			ProductName: product.Name,
			Quantity:    reqItem.Quantity,
			Price:       product.Price,
			Unit:        product.Unit,
			SubTotal:    subTotal,
		})

		// Update stock product
		if req.IsCheckout {
			newStock := product.Stock - reqItem.Quantity
			if err := o.productRepo.UpdateStock(tx, product.ID.String(), newStock); err != nil {
				return 0, nil, nil, err
			}
		}
	}
	return total, items, responses, nil
}

func (o *orderService) saveOrder(
	tx *gorm.DB,
	userID string,
	total int,
	status constants.OrderStatus,
	orderNumber string,
) (*models.Order, error) {
	newOrder := &models.Order{
		Base:        models.Base{CreatedBy: uuid.MustParse(userID)},
		UserID:      uuid.MustParse(userID),
		GrandTotal:  total,
		Status:      status,
		OrderNumber: orderNumber,
	}
	if err := o.orderRepo.CreateOrder(tx, newOrder); err != nil {
		return nil, err
	}
	return newOrder, nil
}

func (o *orderService) generateOrderNumber(tx *gorm.DB) (string, error) {
	dateStr := time.Now().Format("20060102")

	_, err := o.orderSequenceRepo.GetNextSequence(tx, dateStr)
	if err != nil {
		return "", err
	}

	randomStr, err := utils.GenerateRandomString(6)
	if err != nil {
		return "", err
	}

	return "ORD-" + dateStr + "-" + randomStr, nil
}

func (o *orderService) saveOrderItems(
	tx *gorm.DB,
	orderID uuid.UUID,
	items []models.OrderItem,
) error {
	// Mapping OrderID before bulk insert OrderItems
	for i := range items {
		items[i].OrderID = orderID
	}
	return o.orderRepo.CreateOrderItems(tx, items)
}

func (o *orderService) UpdateOrder(orderID string, req dtos.OrderRequest, user dtos.CurrentUser) (dtos.OrderResponse, error) {
	// Set Status
	if req.IsCheckout {
		req.Status = constants.Pending
	} else {
		req.Status = constants.Draft
	}

	// Begin Transaction
	tx := o.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get Order preload OrderItems with UPDATE locking
	order, err := o.orderQueryRepo.FindByIDWithLock(tx, orderID)
	if err != nil {
		tx.Rollback()
		return dtos.OrderResponse{},
			utils.NotFoundException("Order not found")
	}

	// Validate only admin or order creator can update order
	if user.Role != constants.Admin && order.UserID.String() != user.ID {
		tx.Rollback()
		return dtos.OrderResponse{},
			utils.ForbiddenException("You are not authorized to update this order")
	}

	// Validate only draft or pending order can be updated
	validStatus := []constants.OrderStatus{constants.Draft}
	if !slices.Contains(validStatus, order.Status) {
		tx.Rollback()
		return dtos.OrderResponse{},
			utils.BadRequestException("Only "+string(constants.Draft)+" order can be updated", nil)
	}

	// Update Order
	total, items, responses, err := o.prepareOrderData(tx, req, user)
	if err != nil {
		tx.Rollback()
		return dtos.OrderResponse{}, err
	}
	order.GrandTotal = total
	order.Status = req.Status
	if err := o.orderRepo.UpdateOrder(tx, order); err != nil {
		tx.Rollback()
		return dtos.OrderResponse{}, err
	}

	// Update OrderItems
	if err := o.orderRepo.DeleteOrderItems(tx, orderID); err != nil {
		tx.Rollback()
		return dtos.OrderResponse{}, err
	}
	if err := o.saveOrderItems(tx, order.ID, items); err != nil {
		tx.Rollback()
		return dtos.OrderResponse{}, err
	}

	// Commit Transaction
	if err := tx.Commit().Error; err != nil {
		return dtos.OrderResponse{}, err
	}

	// send notification
	if req.IsCheckout {
		o.notifService.QueueSendEmail(
			user.Email,
			"Order Confirmation",
			"Your order "+order.ID.String()+" has been successfully placed.",
		)
	}

	return dtos.OrderResponse{
		ID:         order.ID.String(),
		UserID:     order.UserID.String(),
		Username:   order.User.Username,
		GrandTotal: order.GrandTotal,
		Status:     order.Status,
		OrderItems: responses,
	}, nil
}

func (s *orderService) GetOrders(req dtos.QueryOrderParam, user dtos.CurrentUser) ([]dtos.OrderResponse, int64, error) {
	var userID, storeID string
	var statuses []constants.OrderStatus

	switch user.Role {
	case constants.Buyer:
		userID = user.ID
	case constants.Seller:
		userStore, err := s.storeQueryRepo.FindByUserID(user.ID)
		if err != nil || userStore == nil {
			return []dtos.OrderResponse{}, 0, nil
		}
		storeID = userStore.ID.String()
	}

	if req.Status != nil {
		statuses = []constants.OrderStatus{*req.Status}
	}

	return s.orderQueryRepo.FindAllPagination(userID, storeID, statuses, req)
}

func (s *orderService) CancelOrder(orderID string, user dtos.CurrentUser) (dtos.OrderResponse, error) {
	// Begin Transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get Order preload OrderItems with UPDATE locking
	order, err := s.orderQueryRepo.FindByIDWithLock(tx, orderID)
	if err != nil {
		tx.Rollback()
		return dtos.OrderResponse{},
			utils.NotFoundException("Order not found")
	}

	// Validate only admin or order creator can cancel order
	if user.Role != constants.Admin && order.UserID.String() != user.ID {
		tx.Rollback()
		return dtos.OrderResponse{},
			utils.ForbiddenException("You are not authorized to cancel this order")
	}

	// Validate only draft or pending order can be cancelled
	validStatus := []constants.OrderStatus{constants.Draft, constants.Pending}
	if !slices.Contains(validStatus, order.Status) {
		tx.Rollback()
		return dtos.OrderResponse{},
			utils.BadRequestException(
				"Only "+string(constants.Draft)+" and "+string(constants.Pending)+" orders can be cancelled. Current status: "+string(order.Status),
				nil)
	}

	// Update Order Status
	order.Status = constants.Cancelled
	order.UpdatedBy = uuid.MustParse(user.ID)
	if err := s.orderRepo.UpdateOrder(tx, order); err != nil {
		tx.Rollback()
		return dtos.OrderResponse{}, err
	}

	// Rollback Stock for each OrderItem
	if order.Status != constants.Draft {
		for _, item := range order.OrderItems {
			if err := s.productRepo.AddStock(tx, item.ProductID.String(), item.Quantity); err != nil {
				tx.Rollback()
				return dtos.OrderResponse{}, err
			}
		}
	}

	// Commit Transaction
	if err := tx.Commit().Error; err != nil {
		return dtos.OrderResponse{}, err
	}

	orderItems := make([]dtos.OrderItemResponse, len(order.OrderItems))
	for i, item := range order.OrderItems {
		orderItems[i] = dtos.OrderItemResponse{
			StoreID:     item.StoreID.String(),
			StoreName:   item.Store.Name,
			ProductID:   item.ProductID.String(),
			ProductName: item.Product.Name,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Unit:        item.Product.Unit,
			SubTotal:    item.SubTotal,
		}
	}

	// Send Notification
	s.notifService.QueueSendEmail(
		user.Email,
		"Order has been cancelled",
		"Order with ID "+orderID+" has been cancelled",
	)

	return dtos.OrderResponse{
		ID:         order.ID.String(),
		UserID:     order.UserID.String(),
		Username:   order.User.Username,
		GrandTotal: order.GrandTotal,
		Status:     order.Status,
		OrderItems: orderItems,
	}, nil
}

func (s *orderService) ConfirmOrder(orderID string, user dtos.CurrentUser) (dtos.OrderResponse, error) {
	// Begin Transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get Order preload OrderItems with UPDATE locking
	order, err := s.orderQueryRepo.FindByIDWithLock(tx, orderID)
	if err != nil {
		tx.Rollback()
		return dtos.OrderResponse{},
			utils.NotFoundException("Order not found")
	}

	// Validate only admin or order creator can confirm order
	if user.Role != constants.Admin && order.UserID.String() != user.ID {
		tx.Rollback()
		return dtos.OrderResponse{},
			utils.ForbiddenException("You are not authorized to cancel this order")
	}

	// Validate only pending order can be confirmed
	if order.Status != constants.Pending {
		tx.Rollback()
		return dtos.OrderResponse{},
			utils.BadRequestException(
				"Only "+string(constants.Pending)+" orders can be confirmed. Current status: "+string(order.Status),
				nil)
	}

	// Update Order Status
	order.Status = constants.Paid
	order.UpdatedBy = uuid.MustParse(user.ID)
	if err := s.orderRepo.UpdateOrder(tx, order); err != nil {
		tx.Rollback()
		return dtos.OrderResponse{}, err
	}

	// Commit Transaction
	if err := tx.Commit().Error; err != nil {
		return dtos.OrderResponse{}, err
	}

	orderItems := make([]dtos.OrderItemResponse, len(order.OrderItems))
	for i, item := range order.OrderItems {
		orderItems[i] = dtos.OrderItemResponse{
			StoreID:     item.StoreID.String(),
			StoreName:   item.Store.Name,
			ProductID:   item.ProductID.String(),
			ProductName: item.Product.Name,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Unit:        item.Product.Unit,
			SubTotal:    item.SubTotal,
		}
	}

	// Send Notification
	s.notifService.QueueSendEmail(
		user.Email,
		"Order has been confirmed",
		"Order with ID "+orderID+" has been confirmed",
	)

	return dtos.OrderResponse{
		ID:         order.ID.String(),
		UserID:     order.UserID.String(),
		Username:   order.User.Username,
		GrandTotal: order.GrandTotal,
		Status:     order.Status,
		OrderItems: orderItems,
	}, nil
}
