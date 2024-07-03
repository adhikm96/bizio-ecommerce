package test

import (
	"encoding/json"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	"github.com/Digital-AIR/bizio-ecommerce/internal/service"
	test_util "github.com/Digital-AIR/bizio-ecommerce/test/util"
	"github.com/stretchr/testify/assert"
	"strconv"
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
	assert.Equal(t, model.OrderStatus(model.DeliveryState), orderResp.Status)
	assert.Equal(t, float64(1500), orderResp.FinalAmount)
	assert.Equal(t, float64(0), orderResp.DiscountAmount)
	assert.Equal(t, float64(1500), orderResp.TotalAmount)
	assert.Equal(t, "", orderResp.DiscountCode)

	// check items added in order

	var cnt int
	err = db.Raw("select count(order_items.id) from order_items where order_id = ?", orderResp.Id).Scan(&cnt).Error
	assert.Nil(t, err)
	assert.Equal(t, cnt, 2)

	// fetch order test

	data, _, err = test_util.MakeReqWithBody("GET", "/orders/"+strconv.Itoa(int(orderResp.Id)), new(string), nil)
	assert.Nil(t, err)

	order, err := service.FetchOrder(orderResp.Id)
	assert.Nil(t, err)

	var orderDetail common.OrderDetails

	err = json.Unmarshal(data, &orderDetail)
	assert.Nil(t, err)

	assert.Equal(t, model.OrderStatus(model.DeliveryState), order.Status)
	assert.Equal(t, float64(1500), order.FinalAmount)
	assert.Equal(t, float64(0), order.DiscountAmount)
	assert.Equal(t, float64(1500), order.TotalAmount)
	assert.Equal(t, "", order.DiscountCode)

	orderItems, err := service.FetchOrderItem(order.ID)
	assert.Nil(t, err)

	for _, orderItem := range orderItems {
		if orderItem.PvId == pv2.ID {
			assert.Equal(t, orderItem.Quantity, uint(5))
			assert.Equal(t, orderItem.Price, float64(100))
		} else if orderItem.PvId == pv.ID {
			assert.Equal(t, orderItem.Quantity, uint(10))
			assert.Equal(t, orderItem.Price, float64(100))
		}
	}
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
	assert.Equal(t, model.OrderStatus(model.DeliveryState), orderResp.Status)
	assert.Equal(t, float64(750), orderResp.FinalAmount)
	assert.Equal(t, float64(750), orderResp.DiscountAmount)
	assert.Equal(t, float64(1500), orderResp.TotalAmount)
	assert.Equal(t, discount.Code, orderResp.DiscountCode)

	// check items added in order

	var cnt int
	err = db.Raw("select count(order_items.id) from order_items where order_id = ?", orderResp.Id).Scan(&cnt).Error
	assert.Nil(t, err)
	assert.Equal(t, cnt, 2)

	// fetch order test

	data, _, err = test_util.MakeReqWithBody("GET", "/orders/"+strconv.Itoa(int(orderResp.Id)), new(string), nil)
	assert.Nil(t, err)

	order, err := service.FetchOrder(orderResp.Id)
	assert.Nil(t, err)

	var orderDetail common.OrderDetails

	err = json.Unmarshal(data, &orderDetail)
	assert.Nil(t, err)

	assert.Equal(t, model.OrderStatus(model.DeliveryState), order.Status)
	assert.Equal(t, float64(750), order.FinalAmount)
	assert.Equal(t, float64(750), order.DiscountAmount)
	assert.Equal(t, float64(1500), order.TotalAmount)
	assert.Equal(t, discount.Code, order.DiscountCode)

	orderItems, err := service.FetchOrderItem(order.ID)
	assert.Nil(t, err)

	for _, orderItem := range orderItems {
		if orderItem.PvId == pv2.ID {
			assert.Equal(t, orderItem.Quantity, uint(5))
			assert.Equal(t, orderItem.Price, float64(100))
		} else if orderItem.PvId == pv.ID {
			assert.Equal(t, orderItem.Quantity, uint(10))
			assert.Equal(t, orderItem.Price, float64(100))
		}
	}

	// order status test

	orderStatusUpdate := common.OrderUpdateDto{Status: model.CompleteState}

	_, _, err = test_util.MakeReqWithBody("PUT", "/orders/"+strconv.Itoa(int(orderResp.Id)), orderStatusUpdate, nil)

	assert.Nil(t, err)

	order, err = service.FetchOrder(orderResp.Id)
	assert.Nil(t, err)

	assert.Equal(t, model.OrderStatus(model.CompleteState), order.Status)

	orderStatusUpdate = common.OrderUpdateDto{Status: model.DeliveryState}
	data, response, _ = test_util.MakeReqWithBody("PUT", "/orders/"+strconv.Itoa(int(orderResp.Id)), orderStatusUpdate, nil)
	assert.Equal(t, response.StatusCode, 400)

	assert.Contains(t, string(data), "order already in a final state")

	orderStatusUpdate = common.OrderUpdateDto{Status: "invalid"}
	data, response, _ = test_util.MakeReqWithBody("PUT", "/orders/"+strconv.Itoa(int(orderResp.Id)), orderStatusUpdate, nil)
	assert.Equal(t, response.StatusCode, 400)

	assert.Contains(t, string(data), "invalid status")

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
	assert.Equal(t, model.OrderStatus(model.DeliveryState), orderResp.Status)
	assert.Equal(t, float64(0), orderResp.FinalAmount)
	assert.Equal(t, float64(1500), orderResp.DiscountAmount)
	assert.Equal(t, float64(1500), orderResp.TotalAmount)
	assert.Equal(t, discount.Code, orderResp.DiscountCode)

	// check items added in order

	var cnt int
	err = db.Raw("select count(order_items.id) from order_items where order_id = ?", orderResp.Id).Scan(&cnt).Error
	assert.Nil(t, err)
	assert.Equal(t, cnt, 2)

	// fetch order test

	data, _, err = test_util.MakeReqWithBody("GET", "/orders/"+strconv.Itoa(int(orderResp.Id)), new(string), nil)
	assert.Nil(t, err)

	order, err := service.FetchOrder(orderResp.Id)
	assert.Nil(t, err)

	var orderDetail common.OrderDetails

	err = json.Unmarshal(data, &orderDetail)
	assert.Nil(t, err)

	assert.Equal(t, model.OrderStatus(model.DeliveryState), order.Status)
	assert.Equal(t, float64(0), order.FinalAmount)
	assert.Equal(t, float64(1500), order.DiscountAmount)
	assert.Equal(t, float64(1500), order.TotalAmount)
	assert.Equal(t, discount.Code, order.DiscountCode)

	orderItems, err := service.FetchOrderItem(order.ID)
	assert.Nil(t, err)

	for _, orderItem := range orderItems {
		if orderItem.PvId == pv2.ID {
			assert.Equal(t, orderItem.Quantity, uint(5))
			assert.Equal(t, orderItem.Price, float64(100))
		} else if orderItem.PvId == pv.ID {
			assert.Equal(t, orderItem.Quantity, uint(10))
			assert.Equal(t, orderItem.Price, float64(100))
		}
	}
}
