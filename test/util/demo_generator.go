package test_util

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	"github.com/shopspring/decimal"
	"time"
)

func GetBrand() (*model.Brand, error) {
	brand := model.Brand{
		Name:        RandomString(10),
		Description: "",
	}
	res := database.GetDbConn().Create(&brand)
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
	res := database.GetDbConn().Create(&category)
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
		Price:     decimal.New(1, 2),
		ProductID: product.ID,
	}
	res := database.GetDbConn().Create(&variant)

	if res.Error != nil {
		return nil, res.Error
	}
	return &variant, nil
}

func GetAttribute() (*model.Attribute, error) {
	attr := model.Attribute{
		Name: "attribute1",
	}

	res := database.GetDbConn().Create(&attr)

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
	res := database.GetDbConn().Create(&product)
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

	res := database.GetDbConn().Create(&inventory)

	if res.Error != nil {
		return nil, res.Error
	}

	return &inventory, nil
}

func GetUser() (*model.User, error) {
	user := model.User{
		Username:     RandomString(10),
		Email:        RandomString(10) + "@example.com",
		PasswordHash: "password",
	}
	return &user, database.GetDbConn().Create(&user).Error
}

func GetAddress() (*model.Address, error) {

	user, err := GetUser()
	if err != nil {
		return nil, err
	}

	address := model.Address{
		UserID:       user.ID,
		AddressLine1: "123 Main Street",
		AddressLine2: "line 2",
		City:         "Springfield",
		State:        "IL",
		Zipcode:      "62701",
		Country:      "USA",
	}

	return &address, database.GetDbConn().Create(&address).Error
}

func GetOutOfStockCart(pv *model.ProductVariant) (*model.Cart, error) {

	user, err := GetUser()
	if err != nil {
		return nil, err
	}

	cart := model.Cart{
		UserID: user.ID,
	}

	db := database.GetDbConn()

	if db.Create(&cart).Error != nil {
		return nil, err
	}

	if db.Create(&model.Inventory{
		Quantity:     10,
		ReorderLevel: 5,
		VariantID:    pv.ID,
	}).Error != nil {
		return nil, err
	}

	if db.Create(&model.CartItem{
		Quantity:         20,
		CartID:           cart.ID,
		ProductVariantID: pv.ID,
	}).Error != nil {
		return nil, err
	}

	return &cart, nil
}

func GetDiscount(percentage float64) (*model.Discount, error) {
	discount := model.Discount{
		BaseEntity:         model.BaseEntity{},
		Code:               RandomString(5),
		Description:        "test discount",
		DiscountPercentage: percentage,
		MaxUses:            1,
		CurrentUses:        0,
		StartDate:          time.Now(),
		EndDate:            time.Now().Add(time.Minute * 5),
	}

	return &discount, database.GetDbConn().Create(&discount).Error
}
