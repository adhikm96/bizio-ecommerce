package test

import (
	"encoding/json"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	testutil "github.com/Digital-AIR/bizio-ecommerce/test/util"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
)

func TestCreateCartItem(t *testing.T) {
	db := database.GetDbConn()

	user, err := testutil.GetUser()
	assert.Nil(t, err)

	pv, err := testutil.GetVariant()
	assert.NoError(t, err)

	cartItem, err := testutil.GetCartItem()
	assert.Nil(t, err)

	//create inv
	inv := model.Inventory{VariantID: pv.ID, Quantity: 2000, ReorderLevel: 100}
	db.Create(&inv)

	//clear cart item table
	db.Delete(&cartItem)

	//count cartItem
	var count int64
	db.Model(&model.CartItem{}).Where("id = ?", cartItem.ID).Count(&count)

	dto := common.AddCartItemDto{
		UserID:           user.ID,
		ProductVariantID: pv.ID,
		Quantity:         cartItem.Quantity,
	}

	payload, err := json.Marshal(dto)
	assert.Nil(t, err)

	//add cartItem
	resPayload, resp, err := testutil.MakeReq("POST", "/cart", payload, nil)
	if err != nil {
		slog.Error(err.Error())
		t.Fail()
		return
	}

	var updatedCount int64
	db.Model(&model.CartItem{}).Count(&updatedCount)

	//count updated
	assert.Equal(t, count+1, updatedCount)

	db.Raw("select * from cart_items where product_variant_id = ?", pv.ID).Scan(&pv)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEqual(t, uint(0), cartItem.ID)
	assert.Equal(t, dto.Quantity, cartItem.Quantity)

	slog.Info(string(resPayload))

	//non-existing userId
	dto = common.AddCartItemDto{
		UserID:           uint(rand.Int()),
		ProductVariantID: pv.ID,
		Quantity:         cartItem.Quantity,
	}

	payload2, err := json.Marshal(dto)
	assert.Nil(t, err)

	resPayload, resp, err = testutil.MakeReq("POST", "/cart", payload2, nil)
	if resp.StatusCode != 400 {
		slog.Error(err.Error())
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	slog.Info(string(resPayload))

	//check count updated
	db.Model(&model.CartItem{}).Count(&updatedCount)
	assert.Equal(t, count+1, updatedCount)

	//get cart item api
	resPayload, resp, err = testutil.MakeReq("GET", "/cart/"+strconv.Itoa(int(user.ID)), nil, nil)
	if resp.StatusCode != 200 {
		t.Fail()
		return
	}

	var cartResponse common.CartResponse
	err = json.Unmarshal(resPayload, &cartResponse)
	if err != nil {
		return
	}

	db.Raw("select * from cart_items where id = ?", cartItem.ID).Scan(&cartItem)

	assert.NotNil(t, cartItem)
	assert.Equal(t, 1, len(cartResponse.CartItems))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	slog.Info(string(resPayload))

	// get cart item non-existent userId
	resPayload, resp, _ = testutil.MakeReq("GET", "/cart/100", nil, nil)
	if resp.StatusCode != 400 {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, string(resPayload), "user does not exists with given id")
	slog.Info(string(resPayload))

}

func TestCartItemOutOfStock(t *testing.T) {
	db := database.GetDbConn()

	user, err := testutil.GetUser()
	assert.Nil(t, err)

	pv, err := testutil.GetVariant()
	assert.Nil(t, err)

	//create inv
	inv := model.Inventory{VariantID: pv.ID, Quantity: 2000, ReorderLevel: 100}
	db.Create(&inv)

	//check quantity out-of-stock
	dto := common.AddCartItemDto{
		UserID:           user.ID,
		ProductVariantID: pv.ID,
		Quantity:         50000,
	}

	payload, err := json.Marshal(dto)

	resPayload, resp, err := testutil.MakeReq("POST", "/cart", payload, nil)
	if resp.StatusCode != 400 {
		slog.Error(err.Error())
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, string(resPayload), "item out-of-stock")
	slog.Info(string(resPayload))
}

func TestCartItemQuantity(t *testing.T) {
	db := database.GetDbConn()

	user, err := testutil.GetUser()
	assert.Nil(t, err)

	pv, err := testutil.GetVariant()
	assert.Nil(t, err)

	//create inv
	inv := model.Inventory{VariantID: pv.ID, Quantity: 2000, ReorderLevel: 100}
	db.Create(&inv)

	//quantity negative
	dto := common.AddCartItemDto{
		UserID:           user.ID,
		ProductVariantID: pv.ID,
		Quantity:         -1,
	}

	payload, err := json.Marshal(dto)

	resPayload, resp, err := testutil.MakeReq("POST", "/cart", payload, nil)
	if resp.StatusCode != 400 {
		slog.Error(err.Error())
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, string(resPayload), "quantity cannot be zero or negative")
	slog.Info(string(resPayload))
}

func TestCartItemInvalidVariantID(t *testing.T) {
	db := database.GetDbConn()

	user, err := testutil.GetUser()
	assert.Nil(t, err)

	//create inv
	inv := model.Inventory{VariantID: uint(rand.Int()), Quantity: 2000, ReorderLevel: 100}
	db.Create(&inv)

	//invalid pvId
	dto := common.AddCartItemDto{
		UserID:           user.ID,
		ProductVariantID: inv.VariantID,
		Quantity:         10,
	}

	payload, err := json.Marshal(dto)

	resPayload, resp, err := testutil.MakeReq("POST", "/cart", payload, nil)
	if resp.StatusCode != 400 {
		slog.Error(string(resPayload))
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, string(resPayload), "product variant does not exist with given id")
	slog.Info(string(resPayload))
}

func TestRemoveCartItem(t *testing.T) {
	db := database.GetDbConn()

	cart, err := testutil.GetCart()
	assert.Nil(t, err)

	pv, err := testutil.GetVariant()
	assert.Nil(t, err)

	cartItem := model.CartItem{CartID: cart.ID, ProductVariantID: pv.ID, Quantity: 10}
	db.Create(&cartItem)

	var count int64
	db.Model(&model.CartItem{}).Count(&count)

	// delete cart item non-existent cartItemId
	resPayload, resp, err := testutil.MakeReq("DELETE", "/cart/"+strconv.Itoa(rand.Int()), nil, nil)
	if resp.StatusCode != 400 {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, string(resPayload), "cartItem not found")

	var updatedCount int64
	db.Model(&model.CartItem{}).Count(&updatedCount)

	assert.Equal(t, count, updatedCount)

	// delete cart item exists cartItemId
	resPayload, resp, err = testutil.MakeReq("DELETE", "/cart/"+strconv.Itoa(int(cartItem.ID)), nil, nil)
	if resp.StatusCode != 200 {
		t.Fail()
		return
	}

	//after delete cartItem
	db.Model(&model.CartItem{}).Count(&updatedCount)
	assert.Equal(t, count-1, updatedCount)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(resPayload), "cart item deleted successfully")

}
