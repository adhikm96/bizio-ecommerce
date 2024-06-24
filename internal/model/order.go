package model

type Order struct {
	BaseEntity
	UserID         uint    `json:"user_id" gorm:"index:idx_orders_user_id"`
	AddressID      uint    `json:"address_id"`
	TotalAmount    float64 `json:"total_amount" gorm:"type:decimal(10, 2); not null"`
	DiscountAmount float64 `json:"discount_amount" gorm:"type:decimal(10, 2); default:0.0"`
	FinalAmount    float64 `json:"final_amount" gorm:"type:decimal(10, 2); not null"`
	DiscountCode   string  `json:"discount_code" gorm:"type:varchar(50)"`
	Status         string  `json:"status" gorm:"type:varchar(50); index:idx_orders_status; not null"`
}
