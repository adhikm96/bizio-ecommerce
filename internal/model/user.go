package model

type User struct {
	BaseEntity
	Username     string `json:"username" gorm:"type:varchar(100);size:100;uniqueIndex;not null"`
	Email        string `json:"email" gorm:"type:varchar(100);size:100;uniqueIndex;not null"`
	PasswordHash string `json:"password_hash" gorm:"type:varchar(255);size:255;not null"`
}
