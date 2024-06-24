package model

type Cart struct {
	BaseEntity
	UserID uint `json:"user_id" gorm:"index:idx_carts_user_id"`
	User   User `gorm:"constraint:OnDelete:CASCADE;"`
}
