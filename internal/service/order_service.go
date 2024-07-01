package service

import (
	"errors"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	"log/slog"
)

var orderSlice = []string{
	"created",
	"paid",
	"cancelled",
	"completed",
}

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
		Status:         "created",
	}

	db := database.GetDbConn()

	res := db.Create(&order)

	if res.Error != nil {
		slog.Error(res.Error.Error())
		return nil, res.Error
	}

	// create order items
	err := createOrderItems(orderCreateDto.CartId, order.ID)

	if err != nil {
		return nil, err
	}

	return MakeOrderResp(order), nil
}

func createOrderItems(cartId uint, orderId uint) error {
	db := database.GetDbConn()
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

	orderDetails := common.OrderDetails{
		ID:             order.ID,
		UserID:         order.UserID,
		AddressID:      order.AddressID,
		TotalAmount:    order.TotalAmount,
		DiscountAmount: order.DiscountAmount,
		FinalAmount:    order.FinalAmount,
		DiscountCode:   order.DiscountCode,
		Status:         order.Status,
	}

	return &orderDetails, nil
}
