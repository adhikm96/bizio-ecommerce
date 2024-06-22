package model

type User struct {
	BaseEntity
	Username     string `json:"username" gorm:"size:100, not null"`
	Email        string `json:"email" gorm:"size:100, not null"`
	PasswordHash string `json:"password_hash" gorm:"size:255, not null"`
}
