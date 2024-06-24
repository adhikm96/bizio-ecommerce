package model

type OrderItem struct {
	BaseEntity
	OrderID          uint           `json:"order_id" gorm:"index:idx_order_items_order_id"`
	Order            Order          `gorm:"constraint:OnDelete:CASCADE"`
	ProductVariantID uint           `json:"variant_id" gorm:"index:idx_order_items_variant_id"`
	ProductVariant   ProductVariant `gorm:"foreignKey:variant_id"`
	Quantity         uint           `json:"quantity" gorm:"not null"`
	Price            float64        `json:"price" gorm:"type:decimal(10, 2); not null"`
}
