package inventory

import (
	"encoding/json"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/service"
	"log/slog"
	"net/http"
)

func FetchInventoryHandler(writer http.ResponseWriter, request *http.Request) {
	invID, err := common.FetchPathVariable(writer, request, "id")

	if err != nil {
		return
	}

	// fetch inventory details
	inventory, err := service.FetchInventory(uint(invID))

	err = json.NewEncoder(writer).Encode(inventory)
	if err != nil {
		slog.Error(err.Error())
		common.HandleErrorRes(writer, map[string]string{"message": "failed to fetch inventory"})
	}
}

func UpdateInventoryHandler(writer http.ResponseWriter, request *http.Request) {
	invID, err := common.FetchPathVariable(writer, request, "id")

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

	err = service.UpdateInventory(uint(invID), invUpdateDto)

	if err != nil {
		common.HandleErrorRes(writer, map[string]string{"message": err.Error()})
	}
}

func ValidReqPayload(dto common.InventoryUpdateDto) map[string]string {
	errMap := make(map[string]string)

	if dto.Quantity <= 0 {
		errMap["quantity"] = "quantity should be greater than zero"
	}

	if dto.ReorderLevel <= 0 {
		errMap["reorder_level"] = "reorder_level should be greater than zero"
	}

	return errMap
}
