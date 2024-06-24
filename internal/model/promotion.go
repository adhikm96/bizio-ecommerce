package model

import "time"

type Promotion struct {
	BaseEntity
	Name               string    `json:"name" gorm:"type:varchar(255); not null"`
	Description        string    `json:"description" gorm:"type:text"`
	DiscountPercentage float64   `json:"discount_percentage" gorm:"type:decimal(5,2); not null"`
	StartDate          time.Time `json:"start_date" gorm:"index:idx_promotions_start_date; not null"`
	EndDate            time.Time `json:"end_date" gorm:"index:idx_promotions_end_date; not null"`
}
