package service

import (
	"errors"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
)

func FetchInventory(variantId uint) (*common.InventoryDetail, error) {
	var inv common.InventoryDetail
	db := database.NewDatabaseConnection()
	db.Raw("select id, quantity, reorder_level, variant_id from inventories where variant_id = ?", variantId).Scan(&inv)

	if inv.Id == 0 {
		return nil, errors.New("inventory not found")
	}
	return &inv, nil
}

func UpdateInventory(variantId uint, dto common.InventoryUpdateDto) error {
	inventory := model.Inventory{}

	db := database.NewDatabaseConnection()

	db.Find(&inventory, "variant_id = ?", variantId)

	if inventory.ID == 0 {
		return errors.New("inventory does not exists with given variant id")
	}

	inventory.Quantity = dto.Quantity
	inventory.ReorderLevel = dto.ReorderLevel

	db.Save(&inventory)
	return nil
}
