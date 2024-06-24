package model

import "time"

type Discount struct {
	BaseEntity
	Code               string    `json:"code" gorm:"type:varchar(50); uniqueKey; not null; index:idx_discounts_code"`
	Description        string    `json:"description"`
	DiscountPercentage float64   `json:"discount_percentage" gorm:"type:decimal(5, 2); not null"`
	MaxUses            int       `json:"max_uses" gorm:"type:int; default:1"`
	CurrentUses        int       `json:"current_uses" gorm:"type:int; default:0"`
	StartDate          time.Time `json:"start_date" gorm:"index:idx_discounts_start_date; not null"`
	EndDate            time.Time `json:"end_date" gorm:"index:idx_discounts_end_date; not null"`
}
