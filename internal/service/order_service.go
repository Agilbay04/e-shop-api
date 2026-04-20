package service

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/model"
	"e-shop-api/internal/pkg/util"
	"e-shop-api/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(req dto.OrderRequest, user dto.CurrentUser) (dto.OrderResponse, error)
}

type orderService struct {
	db *gorm.DB
	orderRepo repository.OrderRepository
	productRepo repository.ProductRepository
	productQueryRepo repository.ProductQueryRepository
}

func NewOrderService(
	db *gorm.DB, 
	orderRepo repository.OrderRepository, 
	productRepo repository.ProductRepository,
	productQueryRepo repository.ProductQueryRepository,
) OrderService {
	return &orderService{
		db,
		orderRepo,
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
		newStock := product.Stock - reqItem.Quantity
        if err := o.productRepo.UpdateStock(tx, product.ID.String(), newStock); err != nil {
            return 0, nil, nil, err
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

