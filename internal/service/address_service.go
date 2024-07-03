package service

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
)

func CheckAddressExists(addressId uint) bool {
	db := database.GetDbConn()
	address := model.Address{}

	db.Find(&address, "id = ?", addressId)
	return address.ID != 0
}
