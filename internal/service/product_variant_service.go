package service

import "github.com/Digital-AIR/bizio-ecommerce/internal/database"

func CheckVariantExists(productVariantId uint) bool {
	db := database.NewDatabaseConnection()

	var productVariantExists bool
	db.Raw("select exists(select 1 from product_variants where id = ?)", productVariantId).Scan(&productVariantExists)

	return productVariantExists
}
