package model

type Product struct {
	BaseEntity
	Name            string           `json:"name" gorm:"type:varchar(255); not null"`
	Description     string           `json:"description" gorm:"text"`
	CategoryID      uint             `gorm:"index:idx_products_category_id"`
	BrandID         uint             `gorm:"index:idx_products_brand_id"`
	ProductVariants []ProductVariant `json:"product_variants" gorm:"foreignKey:ProductID"`
}
