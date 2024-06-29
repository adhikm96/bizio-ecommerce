package inventory

import (
	"encoding/json"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/service"
	"log/slog"
	"net/http"
	"strconv"
)

const HighStockQuantityThreshold = 99999

func FetchInventoryHandler(writer http.ResponseWriter, request *http.Request) {
	variantId, err := common.FetchPathVariable(writer, request, "variantId")

	if err != nil {
		return
	}

	// fetch inventory details
	inventory, err := service.FetchInventory(uint(variantId))

	if err != nil {
		common.HandleErrorRes(writer, map[string]string{"message": err.Error()})
		return
	}

	err = json.NewEncoder(writer).Encode(inventory)
	if err != nil {
		slog.Error(err.Error())
		common.HandleErrorRes(writer, map[string]string{"message": "failed to fetch inventory"})
	}
}

func UpdateInventoryHandler(writer http.ResponseWriter, request *http.Request) {
	variantId, err := common.FetchPathVariable(writer, request, "variantId")

	if err != nil {
		return
	}

	invUpdateDto := common.InventoryUpdateDto{}

	if !common.ReadReqPayload(writer, request, &invUpdateDto) {
		return
	}

	// validate payload

	if errMap := ValidReqPayload(invUpdateDto); len(errMap) > 0 {
		common.HandleErrorRes(writer, errMap)
		return
	}

	// TODO : admin user check pending

	err = service.UpdateInventory(uint(variantId), invUpdateDto)

	if err != nil {
		common.HandleErrorRes(writer, map[string]string{"message": err.Error()})
	}
}

func ValidReqPayload(dto common.InventoryUpdateDto) map[string]string {
	errMap := make(map[string]string)

	if dto.Quantity <= 0 {
		errMap["quantity"] = "quantity should be greater than zero"
	}

	if dto.Quantity > HighStockQuantityThreshold {
		errMap["quantity"] = "quantity should be less or equals " + strconv.Itoa(HighStockQuantityThreshold)
	}

	if dto.ReorderLevel <= 0 {
		errMap["reorder_level"] = "reorder_level should be greater than zero"
	}

	return errMap
}
