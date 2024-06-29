package service

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
)

func FindOutOfStockItems(cartId uint) []uint {
	db := database.NewDatabaseConnection()
	var outOfStockCartItems []uint

	// fetch cartItemId where variant's quantity < cart item quantity
	db.Raw("select cart_items.id from cart_items join inventories on inventories.variant_id = cart_items.product_variant_id where cart_items.cart_id = ? and inventories.quantity < cart_items.quantity", cartId).Scan(&outOfStockCartItems)
	return outOfStockCartItems
}
