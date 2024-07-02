package service

import (
	"errors"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	"gorm.io/gorm"
	"log/slog"
	"slices"
)

func CreateOrder(orderCreateDto *common.OrderCreateDto) (*common.OrderResp, error) {
	totalAmt := FetchTotalAmtFromCart(orderCreateDto.CartId)
	discount := CalDiscount(orderCreateDto.DiscountCode, totalAmt)
	finalAmt := totalAmt - discount

	order := model.Order{
		UserID:         orderCreateDto.UserId,
		AddressID:      orderCreateDto.AddressId,
		TotalAmount:    totalAmt,
		DiscountAmount: discount,
		FinalAmount:    finalAmt,
		DiscountCode:   orderCreateDto.DiscountCode,
		Status:         model.DeliveryState,
	}

	db := database.GetDbConn()

	err := db.Transaction(func(tx *gorm.DB) error {
		// creating order
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// creating order items
		if err := createOrderItems(orderCreateDto.CartId, order.ID, tx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return MakeOrderResp(order), nil
}

func createOrderItems(cartId uint, orderId uint, db *gorm.DB) error {
	var results []common.OrderItemCreate

	// pvId, quantity, price
	err := db.Raw("select cart_items.quantity as quantity, product_variants.price as price, product_variants.id as pv_id from cart_items join product_variants on product_variants.id = cart_items.product_variant_id where cart_id = ?", cartId).Scan(&results).Error

	if err != nil {
		return err
	}

	var items = make([]model.OrderItem, len(results))

	for index, result := range results {
		orderItem := model.OrderItem{
			OrderID:          orderId,
			ProductVariantID: result.PvId,
			Quantity:         result.Quantity,
			Price:            result.Price,
		}

		items[index] = orderItem
	}

	return db.Create(items).Error
}

func FetchOrder(orderId uint) (*model.Order, error) {
	var order model.Order

	db := database.GetDbConn()

	db.Find(&order, orderId)

	if order.ID == 0 {
		return nil, errors.New("order not found")
	}

	return &order, nil
}

func UpdateOrderStatus(orderId uint, status model.OrderStatus) error {
	order, err := FetchOrder(orderId)

	if err != nil {
		return err
	}

	if slices.Contains(model.FinalOrderState, string(order.Status)) {
		return errors.New("order already in a final state")
	}

	order.Status = status
	db := database.GetDbConn()
	return db.Save(order).Error
}

func MakeOrderResp(order model.Order) *common.OrderResp {
	return &common.OrderResp{
		TotalAmount:    order.TotalAmount,
		DiscountAmount: order.DiscountAmount,
		FinalAmount:    order.FinalAmount,
		DiscountCode:   order.DiscountCode,
		Status:         order.Status,
		Id:             order.ID,
	}
}

func FetchTotalAmtFromCart(cartId uint) float64 {
	total := 0.0
	err := database.GetDbConn().Raw("select sum(product_variants.price * cart_items.quantity) from cart_items join product_variants on product_variants.id = cart_items.product_variant_id where cart_items.cart_id = ?", cartId).Scan(&total).Error

	if err != nil {
		slog.Error(err.Error())
	}
	return total
}

func FetchOrderDetails(orderId uint) (*common.OrderDetails, error) {
	order, err := FetchOrder(orderId)

	if err != nil {
		return nil, err
	}

	items, err := FetchOrderItem(orderId)

	if err != nil {
		slog.Error(err.Error())
		return nil, errors.New("failed to fetch order items")
	}

	orderDetails := common.OrderDetails{
		ID:             order.ID,
		UserID:         order.UserID,
		AddressID:      order.AddressID,
		TotalAmount:    order.TotalAmount,
		DiscountAmount: order.DiscountAmount,
		FinalAmount:    order.FinalAmount,
		DiscountCode:   order.DiscountCode,
		Status:         order.Status,
		Items:          items,
	}

	return &orderDetails, nil
}

func FetchOrderItem(orderId uint) ([]*common.OrderItemDetail, error) {
	var items []*common.OrderItemDetail
	return items, database.GetDbConn().Raw("select order_items.id as id, order_items.quantity, order_items.price, order_items.product_variant_id as pv_id from order_items where order_id = ?", orderId).Scan(&items).Error
}
