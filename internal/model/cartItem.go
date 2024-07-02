package model

type CartItem struct {
	BaseEntity
	CartID           uint           `json:"cart_id"`
	Cart             Cart           `gorm:"constraint:OnDelete:CASCADE"`
	ProductVariantID uint           `json:"variant_id" gorm:"index:idx_cart_items_variant_id"`
	ProductVariant   ProductVariant `gorm:"foreignKey:ProductVariantID"`
	Quantity         int            `json:"quantity" gorm:"not null"`
}
