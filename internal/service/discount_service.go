package service

import (
	"errors"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	"log/slog"
	"time"
)

func CalDiscount(discountCode string, amt float64) float64 {
	discount, err := FetchDiscount(discountCode)

	if err != nil {
		slog.Debug(err.Error())
		return 0
	}

	// if not valid return discount as 0
	if !IsValidDiscount(discount) {
		return 0
	}

	return common.TwoDigitPrecision(amt * discount.DiscountPercentage / 100)
}

func IsValidDiscount(discount *model.Discount) bool {
	if discount.ID == 0 {
		return false
	}

	if discount.CurrentUses >= discount.MaxUses {
		slog.Debug(discount.Code + "discount has reached max usage ")
		return false
	}

	if discount.EndDate.Before(time.Now()) || discount.StartDate.After(time.Now()) {
		slog.Debug(discount.Code + " has expired or not yet started")
		return false
	}

	return true
}

func FetchDiscount(discountCode string) (*model.Discount, error) {
	discount := model.Discount{}

	db := database.GetDbConn()

	db.Find(&discount, "code = ?", discountCode)

	if discount.ID == 0 {
		return nil, errors.New("discount not found")
	}

	return &discount, nil
}
