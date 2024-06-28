package common

import "github.com/Digital-AIR/bizio-ecommerce/internal/model"

type NotificationCreateDto struct {
	UserId           uint   `json:"user_id"`
	NotificationType string `json:"notification_type"`
	Message          string `json:"message"`
}

type NotificationListDto struct {
	UserID           uint                     `json:"user_id"`
	ID               uint                     `json:"id"`
	NotificationType string                   `json:"notification_type"`
	Message          string                   `json:"message"`
	Status           model.NotificationStatus `json:"status" `
}

type ReviewCreateDto struct {
	UserID  uint   `json:"user_id"`
	Rating  uint   `json:"rating"`
	Comment string `json:"comment"`
}

type ReviewListDto struct {
	UserID    uint   `json:"user_id"`
	ProductID uint   `json:"product_id"`
	Rating    uint   `json:"rating"`
	Comment   string `json:"comment"`
}

type AddCartItemDto struct {
	UserID           uint `json:"user_id"`
	ProductVariantID uint `json:"variant_id"`
	Quantity         int  `json:"quantity"`
}

type CartResponse struct {
	ID        uint           `json:"id"`
	UserID    uint           `json:"user_id"`
	CartItems []CartItemInfo `json:"cart_items"`
}

type CartItemInfo struct {
	CartID    uint `json:"cart_id"`
	VariantID uint `json:"variant_id"`
	Quantity  int  `json:"quantity"`
}
