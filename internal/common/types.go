package common

type InventoryDetail struct {
	Id           uint `json:"id"`
	Quantity     uint `json:"quantity"`
	ReorderLevel uint `json:"reorder_level"`
	VariantID    uint `json:"variant_id"`
}

type InventoryUpdateDto struct {
	Quantity     uint `json:"quantity"`
	ReorderLevel uint `json:"reorder_level"`
}
