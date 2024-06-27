package test

import (
	"encoding/json"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	test_util "github.com/Digital-AIR/bizio-ecommerce/test/util"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"strconv"
	"testing"
)

func TestNewInventoryCreate(t *testing.T) {

	db := database.NewDatabaseConnection()

	variant, err := test_util.GetVariant()
	assert.Nil(t, err)

	inv := model.Inventory{}

	db.Raw("select * from inventories where variant_id = ?", variant.ID).Scan(&inv)

	assert.Equal(t, uint(0), inv.ID)

	invUpdateDto := common.InventoryUpdateDto{
		Quantity:     1000,
		ReorderLevel: 500,
	}

	data, err := json.Marshal(invUpdateDto)
	assert.Nil(t, err)

	_, _, err = test_util.MakeReq("PUT", "/admin/inventory/"+strconv.Itoa(int(variant.ID)), data, nil)
	assert.Nil(t, err)

	db.Raw("select * from inventories where variant_id = ?", variant.ID).Scan(&inv)

	assert.NotEqual(t, uint(0), inv.ID)
	assert.Equal(t, invUpdateDto.Quantity, inv.Quantity)
	assert.Equal(t, invUpdateDto.ReorderLevel, inv.ReorderLevel)
}

func TestInventoryFlow(t *testing.T) {
	inventory, err := test_util.GetInventory()

	assert.Nil(t, err)

	_, response, _ := test_util.MakeReq("GET", "/inventory/0", []byte(""), nil)

	assert.Equal(t, response.StatusCode, 400)

	// test get api
	resPayload, _, err := test_util.MakeReq("GET", "/inventory/"+strconv.Itoa(int(inventory.VariantID)), []byte(""), nil)

	assert.Nil(t, err)

	slog.Info(string(resPayload))

	inventoryDetail := common.InventoryDetail{}

	err = json.Unmarshal(resPayload, &inventoryDetail)
	assert.Nil(t, err)

	assert.Equal(t, inventory.Quantity, inventoryDetail.Quantity)
	assert.Equal(t, inventory.ReorderLevel, inventoryDetail.ReorderLevel)
	assert.Equal(t, inventory.VariantID, inventoryDetail.VariantID)
	assert.Equal(t, inventory.ID, inventoryDetail.Id)

	// update test
	invUpdateDto := common.InventoryUpdateDto{
		Quantity:     0,
		ReorderLevel: 0,
	}

	data, err := json.Marshal(invUpdateDto)
	assert.Nil(t, err)

	_, response, err = test_util.MakeReq("PUT", "/admin/inventory/"+strconv.Itoa(int(inventory.VariantID)), data, nil)
	assert.Equal(t, response.StatusCode, 400)

	invUpdateDto = common.InventoryUpdateDto{
		Quantity:     1000,
		ReorderLevel: 500,
	}

	data, err = json.Marshal(invUpdateDto)
	assert.Nil(t, err)

	_, _, err = test_util.MakeReq("PUT", "/admin/inventory/"+strconv.Itoa(int(inventory.VariantID)), data, nil)
	assert.Nil(t, err)

	database.NewDatabaseConnection().First(&inventory, "id = ?", inventory.ID)

	assert.Equal(t, inventory.Quantity, invUpdateDto.Quantity)
	assert.Equal(t, inventory.ReorderLevel, invUpdateDto.ReorderLevel)

}
