package model

import "time"

type BaseEntity struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;type:serial"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
