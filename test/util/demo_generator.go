package test_util

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"

	"github.com/shopspring/decimal"
)

func GetBrand() (*model.Brand, error) {
	brand := model.Brand{
		Name:        RandomString(10),
		Description: "",
	}
	res := database.NewDatabaseConnection().Create(&brand)
	if res.Error != nil {
		return nil, res.Error
	}
	return &brand, nil
}

func GetCategory() (*model.Category, error) {
	category := model.Category{
		Name:        RandomString(10),
		Description: "",
	}
	res := database.NewDatabaseConnection().Create(&category)
	if res.Error != nil {
		return nil, res.Error
	}
	return &category, nil
}

func GetVariant() (*model.ProductVariant, error) {

	product, err := GetProduct()

	if err != nil {
		return nil, err
	}

	variant := model.ProductVariant{
		Sku:       RandomString(10),
		Price:     decimal.Decimal{},
		ProductID: product.ID,
	}
	res := database.NewDatabaseConnection().Create(&variant)

	if res.Error != nil {
		return nil, res.Error
	}
	return &variant, nil
}

func GetAttribute() (*model.Attribute, error) {
	attr := model.Attribute{
		Name: "attribute1",
	}

	res := database.NewDatabaseConnection().Create(&attr)

	if res.Error != nil {
		return nil, res.Error
	}

	return &attr, nil
}

func GetProduct() (*model.Product, error) {

	cat, err := GetCategory()

	if err != nil {
		return nil, err
	}

	brand, err := GetBrand()
	if err != nil {
		return nil, err
	}

	product := model.Product{
		Name:        RandomString(10),
		Description: "",
		CategoryID:  cat.ID,
		BrandID:     brand.ID,
	}
	res := database.NewDatabaseConnection().Create(&product)
	if res.Error != nil {
		return nil, res.Error
	}
	return &product, nil
}

func GetInventory() (*model.Inventory, error) {
	variant, err := GetVariant()

	if err != nil {
		return nil, err
	}

	inventory := model.Inventory{
		Quantity:     10,
		ReorderLevel: 5,
		VariantID:    variant.ID,
	}

	res := database.NewDatabaseConnection().Create(&inventory)

	if res.Error != nil {
		return nil, res.Error
	}

	return &inventory, nil
}

func GetUser() (*model.User, error) {
	user := model.User{
		Username:     RandomString(5),
		Email:        RandomString(5) + "@example.com",
		PasswordHash: "password",
	}

	res := database.NewDatabaseConnection().Create(&user)

	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func GetCart() (*model.Cart, error) {
	user, err := GetUser()

	if err != nil {
		return nil, err
	}

	cart := model.Cart{
		UserID: user.ID,
	}

	res := database.NewDatabaseConnection().Create(&cart)

	if res.Error != nil {
		return nil, res.Error
	}

	return &cart, nil
}

func GetCartItem() (*model.CartItem, error) {
	cart, err := GetCart()

	if err != nil {
		return nil, err
	}

	pv, err := GetVariant()
	if err != nil {
		return nil, err
	}

	cartItem := model.CartItem{
		CartID:           cart.ID,
		ProductVariantID: pv.ID,
		Quantity:         10,
	}

	res := database.NewDatabaseConnection().Create(&cartItem)

	if res.Error != nil {
		return nil, res.Error
	}

	return &cartItem, nil
}
