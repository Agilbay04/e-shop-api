package service

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/model"
	"e-shop-api/internal/pkg/util"
	"e-shop-api/internal/repository"
	"slices"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(req dto.OrderRequest, user dto.CurrentUser) (dto.OrderResponse, error)
    UpdateOrder(orderId string, req dto.OrderRequest, user dto.CurrentUser) (dto.OrderResponse, error)
    CancelOrder(orderId string, user dto.CurrentUser) (dto.OrderResponse, error)
    ConfirmOrder(orderId string, user dto.CurrentUser) (dto.OrderResponse, error)
}

type orderService struct {
	db *gorm.DB
	orderRepo repository.OrderRepository
    orderQueryRepo repository.OrderQueryRepository
	productRepo repository.ProductRepository
	productQueryRepo repository.ProductQueryRepository
}

func NewOrderService(
	db *gorm.DB, 
	orderRepo repository.OrderRepository, 
    orderQueryRepo repository.OrderQueryRepository,
	productRepo repository.ProductRepository,
	productQueryRepo repository.ProductQueryRepository,
) OrderService {
	return &orderService{
		db,
		orderRepo,
        orderQueryRepo,
		productRepo,
		productQueryRepo,
	}
}

func (o *orderService) CreateOrder(
	req dto.OrderRequest, 
	user dto.CurrentUser,
) (dto.OrderResponse, error) {
    // Set Status
    if req.IsCheckout {
        req.Status = model.Pending
    } else {
        req.Status = model.Draft
    }

    // Begin Transaction
    tx := o.db.Begin()
    
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
        return dto.OrderResponse{}, err
    }

    // Save Order
    newOrder, err := o.saveOrder(tx, user.ID, totalOrderPrice, req.Status)
    if err != nil {
        tx.Rollback()
        return dto.OrderResponse{}, err
    }

    // Save Bulk Order Items
    if err := o.saveOrderItems(tx, newOrder.ID, orderItems); err != nil {
        tx.Rollback()
        return dto.OrderResponse{}, err
    }

    // Commit Transaction
    if err := tx.Commit().Error; err != nil {
        return dto.OrderResponse{}, err
    }

    return dto.OrderResponse{
        ID:         newOrder.ID.String(),
		UserID: 	user.ID.String(),
		Username: 	user.Username,
		GrandTotal: totalOrderPrice,
        Status:     newOrder.Status,
        OrderItems: itemResponses,
    }, nil
}

func (o *orderService) prepareOrderData(
	tx *gorm.DB, 
	req dto.OrderRequest, 
	user dto.CurrentUser,
) (int, []model.OrderItem, []dto.OrderItemResponse, error) {
    var total int
    var items []model.OrderItem
    var responses []dto.OrderItemResponse

    for _, reqItem := range req.OrderItems {
        // Locking for prevent race condition product stock
		product, err := o.productQueryRepo.FindByIDWithLock(tx, reqItem.ProductID.String())
        if err != nil {
            return 0, 
				nil, 
				nil, 
				util.NotFoundException("Product not found: " + reqItem.ProductID.String())
        }
		
		// Validation stock
        if product.Stock < reqItem.Quantity {
            return 0, 
				nil, 
				nil, 
				util.BadRequestException("Insufficient stock for "+product.Name, nil)
        }

        subTotal := product.Price * reqItem.Quantity
        total += subTotal

        items = append(items, model.OrderItem{
            Base:      model.Base{CreatedBy: user.ID},
            StoreID:   product.StoreID,
            ProductID: product.ID,
            Quantity:  reqItem.Quantity,
            Price:     product.Price,
			SubTotal:  subTotal,
        })

        responses = append(responses, dto.OrderItemResponse{
			StoreID:    product.StoreID.String(),
			StoreName:  product.Store.Name,
            ProductID: 	product.ID.String(),
            ProductName:product.Name,
            Quantity:  	reqItem.Quantity,
			Price:  	product.Price,
			Unit: 		product.Unit,	
			SubTotal: 	subTotal,
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
	userID uuid.UUID, 
    total int,
    status model.OrderStatus,
) (*model.Order, error) {
    newOrder := &model.Order{
        Base:       model.Base{CreatedBy: userID},
        UserID:     userID,
        GrandTotal: total,
        Status:     status,
    }
    if err := o.orderRepo.CreateOrder(tx, newOrder); err != nil {
        return nil, err
    }
    return newOrder, nil
}

func (o *orderService) saveOrderItems(
	tx *gorm.DB, 
	orderID uuid.UUID, 
	items []model.OrderItem,
) error {
    // Mapping OrderID before bulk insert OrderItems
    for i := range items {
        items[i].OrderID = orderID
    }
    return o.orderRepo.CreateOrderItems(tx, items)
}

func (o *orderService) UpdateOrder(orderID string, req dto.OrderRequest, user dto.CurrentUser) (dto.OrderResponse, error) {
    // Set Status
    if req.IsCheckout {
        req.Status = model.Pending
    } else {
        req.Status = model.Draft
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
        return dto.OrderResponse{}, 
            util.NotFoundException("Order not found")
    }

    // Validate only admin or order creator can update order
    if user.Role != model.Admin && order.UserID != user.ID {
        tx.Rollback()
        return dto.OrderResponse{}, 
            util.ForbiddenException("You are not authorized to update this order")
    }

    // Validate only draft or pending order can be updated
    validStatus := []model.OrderStatus{model.Draft}
    if !slices.Contains(validStatus, order.Status) {
        tx.Rollback()
        return dto.OrderResponse{}, 
            util.BadRequestException("Only "+string(model.Draft)+" order can be updated", nil)
    }

    // Update Order
    total, items, responses, err := o.prepareOrderData(tx, req, user)
    if err != nil {
        tx.Rollback()
        return dto.OrderResponse{}, err
    }
    order.GrandTotal = total
    order.Status = req.Status
    if err := o.orderRepo.UpdateOrder(tx, order); err != nil {
        tx.Rollback()
        return dto.OrderResponse{}, err
    }

    // Update OrderItems
    if err := o.orderRepo.DeleteOrderItems(tx, orderID); err != nil {
        tx.Rollback()
        return dto.OrderResponse{}, err
    }
    if err := o.saveOrderItems(tx, order.ID, items); err != nil {
        tx.Rollback()
        return dto.OrderResponse{}, err
    }

    // Commit Transaction
    if err := tx.Commit().Error; err != nil {
        return dto.OrderResponse{}, err
    }

    return dto.OrderResponse{
        ID:          order.ID.String(),
        UserID:      order.UserID.String(),
        Username:    order.User.Username,
        GrandTotal:  order.GrandTotal,
        Status:      order.Status,
        OrderItems:  responses,
    }, nil
}

func (s *orderService) CancelOrder(orderID string, user dto.CurrentUser) (dto.OrderResponse, error) {
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
        return dto.OrderResponse{}, 
            util.NotFoundException("Order not found")
    }

    // Validate only admin or order creator can cancel order
    if user.Role != model.Admin && order.UserID != user.ID {
        tx.Rollback()
        return dto.OrderResponse{}, 
            util.ForbiddenException("You are not authorized to cancel this order")
    }

    // Validate only draft or pending order can be cancelled
    validStatus := []model.OrderStatus{model.Draft, model.Pending}
    if !slices.Contains(validStatus, order.Status) {
        tx.Rollback()
        return dto.OrderResponse{}, 
            util.BadRequestException(
                "Only "+string(model.Draft)+" and "+string(model.Pending)+" orders can be cancelled. Current status: "+string(order.Status), 
            nil)
    }

    // Update Order Status
    order.Status = model.Cancelled
    order.UpdatedBy = user.ID
    if err := s.orderRepo.UpdateOrder(tx, order); err != nil {
        tx.Rollback()
        return dto.OrderResponse{}, err
    }

    // Rollback Stock for each OrderItem
    if order.Status != model.Draft {
        for _, item := range order.OrderItems {
            if err := s.productRepo.AddStock(tx, item.ProductID.String(), item.Quantity); err != nil {
                tx.Rollback()
                return dto.OrderResponse{}, err
            }
        }
    }

    // Commit Transaction
    if err := tx.Commit().Error; err != nil {
        return dto.OrderResponse{}, err
    }

    orderItems := make([]dto.OrderItemResponse, len(order.OrderItems))
    for i, item := range order.OrderItems {
        orderItems[i] = dto.OrderItemResponse{
            StoreID:    item.StoreID.String(),
            StoreName:  item.Store.Name,
            ProductID:  item.ProductID.String(),
            ProductName: item.Product.Name,
            Quantity:   item.Quantity,
            Price:      item.Price,
            Unit:       item.Product.Unit,
            SubTotal:   item.SubTotal,
        }
    }

    return dto.OrderResponse{
        ID:          order.ID.String(),
        UserID:      order.UserID.String(),
        Username:    order.User.Username,
        GrandTotal:  order.GrandTotal,
        Status:      order.Status,
        OrderItems:  orderItems,
    }, nil
}

func (s *orderService) ConfirmOrder(orderID string, user dto.CurrentUser) (dto.OrderResponse, error) {
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
        return dto.OrderResponse{}, 
            util.NotFoundException("Order not found")
    }

    // Validate only admin or order creator can confirm order
    if user.Role != model.Admin && order.UserID != user.ID {
        tx.Rollback()
        return dto.OrderResponse{}, 
            util.ForbiddenException("You are not authorized to cancel this order")
    }

     // Validate only pending order can be confirmed
    if order.Status != model.Pending {
        tx.Rollback()
        return dto.OrderResponse{}, 
            util.BadRequestException(
                "Only "+string(model.Pending)+" orders can be confirmed. Current status: "+string(order.Status), 
            nil)
    }

    // Update Order Status
    order.Status = model.Paid
    order.UpdatedBy = user.ID
    if err := s.orderRepo.UpdateOrder(tx, order); err != nil {
        tx.Rollback()
        return dto.OrderResponse{}, err
    }

    // Commit Transaction
    if err := tx.Commit().Error; err != nil {
        return dto.OrderResponse{}, err
    }

    orderItems := make([]dto.OrderItemResponse, len(order.OrderItems))
    for i, item := range order.OrderItems {
        orderItems[i] = dto.OrderItemResponse{
            StoreID:    item.StoreID.String(),
            StoreName:  item.Store.Name,
            ProductID:  item.ProductID.String(),
            ProductName: item.Product.Name,
            Quantity:   item.Quantity,
            Price:      item.Price,
            Unit:       item.Product.Unit,
            SubTotal:   item.SubTotal,
        }
    }

    return dto.OrderResponse{
        ID:          order.ID.String(),
        UserID:      order.UserID.String(),
        Username:    order.User.Username,
        GrandTotal:  order.GrandTotal,
        Status:      order.Status,
        OrderItems:  orderItems,
    }, nil
}

