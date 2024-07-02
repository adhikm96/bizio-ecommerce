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

type OrderCreateDto struct {
	CartId       uint   `json:"cart_id"`
	DiscountCode string `json:"discount_code"`
	UserId       uint   `json:"user_id"`
	AddressId    uint   `json:"address_id"`
}

type OrderUpdateDto struct {
	Status model.OrderStatus `json:"status"`
}

type OrderResp struct {
	Id             uint              `json:"id"`
	TotalAmount    float64           `json:"total_amount"`
	DiscountAmount float64           `json:"discount_amount"`
	FinalAmount    float64           `json:"final_amount"`
	DiscountCode   string            `json:"discount_code"`
	Status         model.OrderStatus `json:"status"`
}

type InventoryDetail struct {
	Id           uint `json:"id"`
	Quantity     uint `json:"quantity"`
	ReorderLevel uint `json:"reorder_level"`
	VariantID    uint `json:"variant_id"`
}

type InventoryUpdateDto struct {
	Quantity     uint `json:"quantity"`
	ReorderLevel uint `json:"reorder_level"`
}

type OrderDetails struct {
	ID             uint               `json:"id"`
	UserID         uint               `json:"user_id"`
	AddressID      uint               `json:"address_id"`
	TotalAmount    float64            `json:"total_amount"`
	DiscountAmount float64            `json:"discount_amount"`
	FinalAmount    float64            `json:"final_amount"`
	DiscountCode   string             `json:"discount_code"`
	Status         model.OrderStatus  `json:"status"`
	Items          []*OrderItemDetail `json:"items"`
}

type OrderItemCreate struct {
	Quantity uint    `json:"quantity"`
	Price    float64 `json:"price"`
	PvId     uint    `json:"pv_id"`
}

type OrderItemDetail struct {
	Id       uint    `json:"id"`
	Quantity uint    `json:"quantity"`
	Price    float64 `json:"price"`
	PvId     uint    `json:"pv_id"`
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
	CartItemID uint `json:"cart_item_id"`
	VariantID  uint `json:"variant_id"`
	Quantity   int  `json:"quantity"`
}
