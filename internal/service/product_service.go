package service

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
)

func CheckProductExists(productId uint) bool {
	db := database.GetDbConn()
	var product model.Product
	if err := db.First(&product, productId).Error; err != nil {
		return false
	}
	return true
}
