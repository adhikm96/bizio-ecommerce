package model

type Attribute struct {
	BaseEntity
	Name string `json:"name" gorm:"type:varchar(100)"`
}

type AttributeValue struct {
	BaseEntity
	AttributeID uint      `json:"attribute_id" gorm:"not null"`
	Attribute   Attribute `gorm:"constraint:OnDelete:CASCADE"`
	Value       string    `json:"value" gorm:"type:varchar(100); not null"`
}
