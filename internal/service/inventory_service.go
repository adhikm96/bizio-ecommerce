package service

import (
	"errors"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
)

func FetchInventory(variantId uint) (*common.InventoryDetail, error) {
	var inv common.InventoryDetail
	db := database.GetDbConn()
	db.Raw("select id, quantity, reorder_level, variant_id from inventories where variant_id = ?", variantId).Scan(&inv)

	if inv.Id == 0 {
		return nil, errors.New("inventory not found")
	}
	return &inv, nil
}

func UpdateInventory(variantId uint, dto common.InventoryUpdateDto) error {
	inventory := model.Inventory{}

	db := database.GetDbConn()

	db.Find(&inventory, "variant_id = ?", variantId)

	if inventory.ID == 0 {
		inventory.VariantID = variantId
	}

	inventory.Quantity = dto.Quantity
	inventory.ReorderLevel = dto.ReorderLevel

	db.Save(&inventory)
	return nil
}
