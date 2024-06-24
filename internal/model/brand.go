package model

type Brand struct {
	BaseEntity
	Name        string    `json:"name" gorm:"type:varchar(100); not null"`
	Description string    `json:"description" gorm:"text"`
	Products    []Product `json:"products" gorm:"foreignKey:BrandID"`
}
