package test

import (
	"encoding/json"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	"github.com/Digital-AIR/bizio-ecommerce/internal/server"
	testutil "github.com/Digital-AIR/bizio-ecommerce/test/util"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestCreateCartItem(t *testing.T) {
	startServer()
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

	//get cart item api
	resPayload, resp, err = testutil.MakeReq("GET", "/cart/"+strconv.Itoa(int(user.ID)), nil, nil)
	if resp.StatusCode != 200 {
		t.Fail()
		return
	}

	db.Raw("select * from cart_items where id = ?", cartItem.ID).Scan(&cartItem)

	assert.NotNil(t, cartItem)

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
	startServer()
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
	startServer()
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
	assert.Contains(t, string(resPayload), "quantity cannot be negative")
	slog.Info(string(resPayload))
}

func TestCartItemInvalidVariantID(t *testing.T) {
	startServer()
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
	startServer()

	cart, err := testutil.GetCart()
	assert.Nil(t, err)

	// delete cart item non-existent cartId
	resPayload, resp, err := testutil.MakeReq("DELETE", "/cart/100", nil, nil)
	if resp.StatusCode != 400 {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, string(resPayload), "cart not found")

	// delete cart item exists cartId
	resPayload, resp, err = testutil.MakeReq("DELETE", "/cart/"+strconv.Itoa(int(cart.ID)), nil, nil)
	if resp.StatusCode != 200 {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(resPayload), "cart item deleted successfully")

}
func startServer() {
	go server.InitServer()
	time.Sleep(time.Second * 1)
}
