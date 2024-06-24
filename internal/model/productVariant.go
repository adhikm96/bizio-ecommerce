package model

import "github.com/shopspring/decimal"

type ProductVariant struct {
	BaseEntity
	Sku         string          `json:"sku" gorm:"type:varchar(100); uniqueIndex:unique_product_variants_sku; not null"`
	Price       decimal.Decimal `json:"price" gorm:"type:decimal(10,2);"`
	ProductID   uint            `gorm:"index:idx_product_variants_product_id"`
	Inventories []Inventory     `json:"inventories" gorm:"foreignKey:VariantID"`
}
