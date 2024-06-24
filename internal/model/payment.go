package model

type Payment struct {
	BaseEntity
	OrderID       uint    `json:"order_id" gorm:"index:idx_payments_order_id"`
	Order         Order   `gorm:"constraint:OnDelete:CASCADE"`
	Amount        float64 `json:"amount" gorm:"type:decimal(10,2); not null"`
	PaymentMethod string  `json:"payment_method" gorm:"type:varchar(50); not null"`
	Status        string  `json:"status" gorm:"type:varchar(50); index:idx_payments_status; not null"`
	TransactionID string  `json:"transaction_id" gorm:"type:varchar(100);"`
}
