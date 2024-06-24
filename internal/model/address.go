package model

type Address struct {
	BaseEntity
	UserID       uint   `json:"user_id"`
	User         User   `gorm:"constraint:OnDelete:CASCADE"`
	AddressLine1 string `json:"address_line1" gorm:"type:varchar(255); not null"`
	AddressLine2 string `json:"address_line2" gorm:"type:varchar(255)"`
	City         string `json:"city" gorm:"type:varchar(100); not null"`
	State        string `json:"state" gorm:"type:varchar(100); not null"`
	Zipcode      string `json:"zipcode" gorm:"type:varchar(20); not null"`
	Country      string `json:"country" gorm:"type:varchar(100); not null"`
}
