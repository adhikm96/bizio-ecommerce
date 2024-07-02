package service

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
)

func FindOutOfStockItems(cartId uint) ([]uint, error) {
	db := database.GetDbConn()
	var outOfStockCartItems []uint

	// fetch cartItemId where variant's quantity < cart item quantity
	return outOfStockCartItems, db.Raw("select cart_items.id from cart_items left join inventories on inventories.variant_id = cart_items.product_variant_id where cart_items.cart_id = ? and (inventories.id is null or inventories.quantity < cart_items.quantity)", cartId).Scan(&outOfStockCartItems).Error
}
