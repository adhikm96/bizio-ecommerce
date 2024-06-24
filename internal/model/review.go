package model

type Review struct {
	BaseEntity
	ProductID uint    `json:"product_id" gorm:"index:idx_reviews_product_id"`
	Product   Product `gorm:"constraint:OnDelete:CASCADE"`
	UserID    uint    `json:"user_id" gorm:"index:idx_reviews_user_id"`
	User      User    `gorm:"constraint:OnDelete:CASCADE"`
	Rating    uint    `json:"rating"`
	Comment   string  `json:"comment" gorm:"type:text"`
}
