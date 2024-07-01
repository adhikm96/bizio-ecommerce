package test

import (
	"encoding/json"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	test_util "github.com/Digital-AIR/bizio-ecommerce/test/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrderCreateWithOutOfStock(t *testing.T) {
	pv, err := test_util.GetVariant()
	assert.Nil(t, err)

	address, err := test_util.GetAddress()
	assert.Nil(t, err)

	cart, err := test_util.GetOutOfStockCart(pv)
	assert.Nil(t, err)

	// creating order from out of stock cart
	orderCreateDto := common.OrderCreateDto{
		CartId:       cart.ID,
		DiscountCode: "",
		UserId:       address.UserID,
		AddressId:    address.ID,
	}
	data, response, err := test_util.MakeReqWithBody("POST", "/orders", orderCreateDto, nil)
	assert.NotNil(t, err)

	assert.Equal(t, response.StatusCode, 400)
	assert.Contains(t, string(data), "some cart items are out of stock")
}

func TestOrderCreateWithStock(t *testing.T) {
	pv, err := test_util.GetVariant()
	assert.Nil(t, err)

	pv2, err := test_util.GetVariant()
	assert.Nil(t, err)

	address, err := test_util.GetAddress()
	assert.Nil(t, err)

	user, err := test_util.GetUser()
	assert.Nil(t, err)

	cart := model.Cart{
		UserID: user.ID,
	}

	db := database.GetDbConn()

	assert.Nil(t, db.Create(&cart).Error)

	assert.Nil(t, db.Create(&model.Inventory{
		Quantity:     10,
		ReorderLevel: 5,
		VariantID:    pv.ID,
	}).Error)

	assert.Nil(t, db.Create(&model.Inventory{
		Quantity:     5,
		ReorderLevel: 1,
		VariantID:    pv2.ID,
	}).Error)

	assert.Nil(t, db.Create(&model.CartItem{
		Quantity:         10,
		CartID:           cart.ID,
		ProductVariantID: pv.ID,
	}).Error)

	assert.Nil(t, db.Create(&model.CartItem{
		Quantity:         5,
		CartID:           cart.ID,
		ProductVariantID: pv2.ID,
	}).Error)

	// creating order from stock cart
	orderCreateDto := common.OrderCreateDto{
		CartId:       cart.ID,
		DiscountCode: "",
		UserId:       address.UserID,
		AddressId:    address.ID,
	}

	data, response, err := test_util.MakeReqWithBody("POST", "/orders", orderCreateDto, nil)
	assert.Nil(t, err)

	var orderResp common.OrderResp

	assert.Nil(t, json.Unmarshal(data, &orderResp))

	assert.Equal(t, response.StatusCode, 200)
	assert.Equal(t, model.OrderStatus("created"), orderResp.Status)
	assert.Equal(t, float64(1500), orderResp.FinalAmount)
	assert.Equal(t, float64(0), orderResp.DiscountAmount)
	assert.Equal(t, float64(1500), orderResp.TotalAmount)
	assert.Equal(t, "", orderResp.DiscountCode)

	// check items added in order

	var cnt int
	err = db.Raw("select count(order_items.id) from order_items where order_id = ?", orderResp.Id).Scan(&cnt).Error
	assert.Nil(t, err)
	assert.Equal(t, cnt, 2)
}

func TestOrderCreateWithStockWithDiscount(t *testing.T) {

	discount, err := test_util.GetDiscount(50)
	assert.Nil(t, err)

	pv, err := test_util.GetVariant()
	assert.Nil(t, err)

	pv2, err := test_util.GetVariant()
	assert.Nil(t, err)

	address, err := test_util.GetAddress()
	assert.Nil(t, err)

	user, err := test_util.GetUser()
	assert.Nil(t, err)

	cart := model.Cart{
		UserID: user.ID,
	}

	db := database.GetDbConn()

	assert.Nil(t, db.Create(&cart).Error)

	assert.Nil(t, db.Create(&model.Inventory{
		Quantity:     10,
		ReorderLevel: 5,
		VariantID:    pv.ID,
	}).Error)

	assert.Nil(t, db.Create(&model.Inventory{
		Quantity:     5,
		ReorderLevel: 1,
		VariantID:    pv2.ID,
	}).Error)

	assert.Nil(t, db.Create(&model.CartItem{
		Quantity:         10,
		CartID:           cart.ID,
		ProductVariantID: pv.ID,
	}).Error)

	assert.Nil(t, db.Create(&model.CartItem{
		Quantity:         5,
		CartID:           cart.ID,
		ProductVariantID: pv2.ID,
	}).Error)

	// creating order from stock cart
	orderCreateDto := common.OrderCreateDto{
		CartId:       cart.ID,
		DiscountCode: discount.Code,
		UserId:       address.UserID,
		AddressId:    address.ID,
	}

	data, response, err := test_util.MakeReqWithBody("POST", "/orders", orderCreateDto, nil)
	assert.Nil(t, err)

	var orderResp common.OrderResp

	assert.Nil(t, json.Unmarshal(data, &orderResp))

	assert.Equal(t, response.StatusCode, 200)
	assert.Equal(t, model.OrderStatus("created"), orderResp.Status)
	assert.Equal(t, float64(750), orderResp.FinalAmount)
	assert.Equal(t, float64(750), orderResp.DiscountAmount)
	assert.Equal(t, float64(1500), orderResp.TotalAmount)
	assert.Equal(t, discount.Code, orderResp.DiscountCode)

	// check items added in order

	var cnt int
	err = db.Raw("select count(order_items.id) from order_items where order_id = ?", orderResp.Id).Scan(&cnt).Error
	assert.Nil(t, err)
	assert.Equal(t, cnt, 2)
}

func TestOrderCreateWithStockWith100PercentDiscount(t *testing.T) {

	discount, err := test_util.GetDiscount(100)
	assert.Nil(t, err)

	pv, err := test_util.GetVariant()
	assert.Nil(t, err)

	pv2, err := test_util.GetVariant()
	assert.Nil(t, err)

	address, err := test_util.GetAddress()
	assert.Nil(t, err)

	user, err := test_util.GetUser()
	assert.Nil(t, err)

	cart := model.Cart{
		UserID: user.ID,
	}

	db := database.GetDbConn()

	assert.Nil(t, db.Create(&cart).Error)

	assert.Nil(t, db.Create(&model.Inventory{
		Quantity:     10,
		ReorderLevel: 5,
		VariantID:    pv.ID,
	}).Error)

	assert.Nil(t, db.Create(&model.Inventory{
		Quantity:     5,
		ReorderLevel: 1,
		VariantID:    pv2.ID,
	}).Error)

	assert.Nil(t, db.Create(&model.CartItem{
		Quantity:         10,
		CartID:           cart.ID,
		ProductVariantID: pv.ID,
	}).Error)

	assert.Nil(t, db.Create(&model.CartItem{
		Quantity:         5,
		CartID:           cart.ID,
		ProductVariantID: pv2.ID,
	}).Error)

	// creating order from stock cart
	orderCreateDto := common.OrderCreateDto{
		CartId:       cart.ID,
		DiscountCode: discount.Code,
		UserId:       address.UserID,
		AddressId:    address.ID,
	}

	data, response, err := test_util.MakeReqWithBody("POST", "/orders", orderCreateDto, nil)
	assert.Nil(t, err)

	var orderResp common.OrderResp

	assert.Nil(t, json.Unmarshal(data, &orderResp))

	assert.Equal(t, response.StatusCode, 200)
	assert.Equal(t, model.OrderStatus("created"), orderResp.Status)
	assert.Equal(t, float64(0), orderResp.FinalAmount)
	assert.Equal(t, float64(1500), orderResp.DiscountAmount)
	assert.Equal(t, float64(1500), orderResp.TotalAmount)
	assert.Equal(t, discount.Code, orderResp.DiscountCode)

	// check items added in order

	var cnt int
	err = db.Raw("select count(order_items.id) from order_items where order_id = ?", orderResp.Id).Scan(&cnt).Error
	assert.Nil(t, err)
	assert.Equal(t, cnt, 2)
}
