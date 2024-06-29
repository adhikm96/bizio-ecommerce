package service

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	"log/slog"
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
		Status:         "created",
	}

	db := database.NewDatabaseConnection()

	res := db.Create(order)

	if res.Error != nil {
		slog.Error(res.Error.Error())
		return nil, res.Error
	}

	return MakeOrderResp(order), nil
}

func MakeOrderResp(order model.Order) *common.OrderResp {
	return &common.OrderResp{
		TotalAmount:    order.TotalAmount,
		DiscountAmount: order.DiscountAmount,
		FinalAmount:    order.FinalAmount,
		DiscountCode:   order.DiscountCode,
		Status:         order.Status,
	}
}

func FetchTotalAmtFromCart(cartId uint) float64 {
	total := 0.0
	database.NewDatabaseConnection().Raw("select sum(product_variants.price) from cart_items join product_variants on product_variants.id = cart_items.product_variant_id = product_variants.id where product_variants.id = ?", cartId).Scan(&total)
	return total
}
