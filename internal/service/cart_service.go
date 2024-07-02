package service

import (
	"errors"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
)

func CheckStockAvailable(quantity int, pvId uint) (bool, error) {
	db := database.GetDbConn()
	var inv model.Inventory

	if err := db.Where("variant_id = ?", pvId).First(&inv).Error; err != nil {
		return false, errors.New("product variant does not exist with given id")
	}

	if int(inv.Quantity) < quantity {
		return false, nil
	}
	return true, nil
}

func AddItemCart(addCartItemDto common.AddCartItemDto) (*model.CartItem, error) {
	db := database.GetDbConn()

	var cart model.Cart
	if err := db.Raw("select * from carts where user_id = ?", &addCartItemDto.UserID).Scan(&cart); err != nil {
		cart = model.Cart{UserID: addCartItemDto.UserID}
		db.Create(&cart)
	}

	// Create new cart item
	cartItem := model.CartItem{
		CartID:           cart.ID,
		ProductVariantID: addCartItemDto.ProductVariantID,
		Quantity:         addCartItemDto.Quantity,
	}

	if err := db.Create(&cartItem).Error; err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func FetchCartItem(userId uint) ([]common.CartResponse, error) {
	db := database.GetDbConn()
	var carts []model.Cart

	db.Preload("CartItems.ProductVariant").Where("user_id = ?", userId).Find(&carts)

	if len(carts) == 0 {
		return nil, errors.New("cart_item not found for user")
	}

	var cartResponses []common.CartResponse

	for _, cart := range carts {
		cartResp := common.CartResponse{
			ID:     cart.ID,
			UserID: cart.UserID,
		}

		for _, item := range cart.CartItems {
			cartItem := common.CartItemInfo{
				CartItemID: item.ID,
				VariantID:  item.ProductVariantID,
				Quantity:   item.Quantity,
			}
			cartResp.CartItems = append(cartResp.CartItems, cartItem)
		}
		cartResponses = append(cartResponses, cartResp)
	}

	return cartResponses, nil
}

func RemoveCartItem(cartItemId uint) error {
	db := database.GetDbConn()
	var cartItem model.CartItem

	db.First(&cartItem, cartItemId)

	if cartItem.ID == 0 {
		return errors.New("cartItem not found")
	}

	if err := db.Where("id = ?", cartItemId).Delete(&model.CartItem{}).Error; err != nil {
		return errors.New("failed to delete cart items")
	}
	return nil
}
