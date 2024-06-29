package service

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	"log/slog"
)

func CreateOrder(orderCreateDto *common.OrderCreateDto) (*model.Order, error) {
	order := model.Order{
		UserID:         0,
		AddressID:      0,
		TotalAmount:    0,
		DiscountAmount: 0,
		FinalAmount:    0,
		DiscountCode:   "",
		Status:         "",
	}

	db := database.NewDatabaseConnection()

	res := db.Create(order)

	if res.Error != nil {
		slog.Error(res.Error.Error())
		return nil, res.Error
	}

	return &order, nil
}
