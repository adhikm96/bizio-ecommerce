package model

type VariantAttribute struct {
	BaseEntity
	ProductVariantID uint           `json:"variant_id"`
	ProductVariant   ProductVariant `gorm:"constraint:OnDelete:CASCADE"`
	AttributeID      uint           `json:"attribute_id"`
	Attribute        Attribute      `gorm:"constraint:OnDelete:CASCADE"`
	AttributeValueID uint           `json:"attribute_value_id"`
	AttributeValue   AttributeValue `gorm:"constraint:OnDelete:CASCADE"`
}
