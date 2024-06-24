package model

type Inventory struct {
	BaseEntity
	Quantity     uint `json:"quantity"`
	ReorderLevel uint `json:"reorder_level"`
	VariantID    uint `gorm:"uniqueIndex:unique_inventories_variant_id"`
}
