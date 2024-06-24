package model

type Inventory struct {
	BaseEntity
	Quantity     uint `json:"quantity" gorm:"not null"`
	ReorderLevel uint `json:"reorder_level" gorm:"not null"`
	VariantID    uint `gorm:"uniqueIndex:unique_inventories_variant_id"`
}
