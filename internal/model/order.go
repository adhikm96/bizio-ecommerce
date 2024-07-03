package model

type OrderStatus string

const DeliveryState = "delivery"
const PaymentState = "payment"
const ConfirmState = "confirm"
const CompleteState = "complete"

var ValidOrdersStatus = []string{
	DeliveryState,
	PaymentState,
	ConfirmState,
	CompleteState,
}

var FinalOrderState = []string{
	CompleteState,
}

type Order struct {
	BaseEntity
	UserID         uint        `json:"user_id" gorm:"index:idx_orders_user_id"`
	User           User        `gorm:"foreignKey:UserID"`
	AddressID      uint        `json:"address_id"`
	Address        Address     `gorm:"foreignKey:AddressID"`
	TotalAmount    float64     `json:"total_amount" gorm:"type:decimal(10, 2); not null"`
	DiscountAmount float64     `json:"discount_amount" gorm:"type:decimal(10, 2); default:0.0"`
	FinalAmount    float64     `json:"final_amount" gorm:"type:decimal(10, 2); not null"`
	DiscountCode   string      `json:"discount_code" gorm:"type:varchar(50)"`
	Status         OrderStatus `json:"status" gorm:"type:varchar(50); index:idx_orders_status; not null"`
}
