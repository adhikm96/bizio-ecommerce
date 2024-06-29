package order

import (
	"encoding/json"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/service"
	"log/slog"
	"net/http"
)

func UpdateOrderStatusHandler(writer http.ResponseWriter, request *http.Request) {

}

func FetchOrderHandler(writer http.ResponseWriter, request *http.Request) {

}

func CreateOrderHandler(writer http.ResponseWriter, request *http.Request) {

	createOrderDto := common.OrderCreateDto{}

	// fetch body
	if !common.ReadReqPayload(writer, request, &createOrderDto) {
		return
	}

	// validate the req
	if errMap := validateOrderCreate(createOrderDto); len(errMap) > 0 {
		common.HandleErrorRes(writer, errMap)
		return
	}

	// check cartItem is available
	if outOfStockCartItems := service.FindOutOfStockItems(createOrderDto.CartId); len(outOfStockCartItems) > 0 {
		common.HandleErrorRes(writer, map[string]string{"message": "some cart items are out of stock"})
		return
	}

	// create order
	orderResp, err := service.CreateOrder(&createOrderDto)

	if err != nil {
		common.HandleErrorRes(writer, map[string]string{"message": "failed to create order"})
	}

	err = json.NewEncoder(writer).Encode(orderResp)

	if err != nil {
		slog.Error(err.Error())
		common.HandleErrorRes(writer, map[string]string{"message": "some error occurred while creating order"})
	}
}

func validateOrderCreate(dto common.OrderCreateDto) map[string]string {
	// nothing to validate as of now
	return map[string]string{}
}
