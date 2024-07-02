package order

import (
	"encoding/json"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	"github.com/Digital-AIR/bizio-ecommerce/internal/service"
	"log/slog"
	"net/http"
	"slices"
	"strings"
)

func UpdateOrderStatusHandler(writer http.ResponseWriter, request *http.Request) {
	orderId, err := common.FetchPathVariable(writer, request, "order_id")

	if err != nil {
		return
	}

	orderUpdateDto := common.OrderUpdateDto{}

	if !common.ReadReqPayload(writer, request, &orderUpdateDto) {
		return
	}

	// validate
	if errMap := validateOrderUpdate(orderUpdateDto); len(errMap) > 0 {
		common.HandleErrorRes(writer, errMap)
		return
	}

	err = service.UpdateOrderStatus(uint(orderId), orderUpdateDto.Status)

	if err != nil {
		common.HandleErrorRes(writer, map[string]string{"message": err.Error()})
	}
}

func validateOrderUpdate(dto common.OrderUpdateDto) map[string]string {
	errMap := make(map[string]string)

	if dto.Status == "" {
		errMap["status"] = "status required"
	}
	if dto.Status != "" && !slices.Contains(model.ValidOrdersStatus, strings.ToLower(string(dto.Status))) {
		errMap["status"] = "invalid status"
	}

	return errMap
}

func FetchOrderHandler(writer http.ResponseWriter, request *http.Request) {
	orderId, err := common.FetchPathVariable(writer, request, "order_id")

	if err != nil {
		return
	}

	// fetch order details

	orderDetails, err := service.FetchOrderDetails(uint(orderId))

	if err != nil {
		common.HandleErrorRes(writer, map[string]string{"message": err.Error()})
		return
	}

	err = json.NewEncoder(writer).Encode(orderDetails)
	if err != nil {
		common.HandleErrorRes(writer, map[string]string{"message": "failed to fetch order details"})
	}
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
	if outOfStockCartItems, err := service.FindOutOfStockItems(createOrderDto.CartId); len(outOfStockCartItems) > 0 {
		if err != nil {
			slog.Error(err.Error())
			common.HandleErrorRes(writer, map[string]string{"message": "failed to create order"})
			return
		}
		common.HandleErrorRes(writer, map[string]string{"message": "some cart items are out of stock"})
		return
	}

	// create order
	orderResp, err := service.CreateOrder(&createOrderDto)

	if err != nil {
		common.HandleErrorRes(writer, map[string]string{"message": "failed to create order"})
	}

	err = json.NewEncoder(writer).Encode(&orderResp)

	if err != nil {
		slog.Error(err.Error())
		common.HandleErrorRes(writer, map[string]string{"message": "some error occurred while creating order"})
	}
}

func validateOrderCreate(dto common.OrderCreateDto) map[string]string {
	// nothing to validate as of now
	return map[string]string{}
}
