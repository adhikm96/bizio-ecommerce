package model

type ProductPromotion struct {
	BaseEntity
	ProductID   uint      `json:"product_id" gorm:"index:idx_product_promotions_product_id"`
	Product     Product   `gorm:"constraint:OnDelete:CASCADE"`
	PromotionID uint      `json:"promotion_id" gorm:"index:idx_product_promotions_promotion_id"`
	Promotion   Promotion `gorm:"constraint:OnDelete:CASCADE"`
}
