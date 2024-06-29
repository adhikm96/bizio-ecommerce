package order

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"net/http"
)

func UpdateOrderStatusHandler(writer http.ResponseWriter, request *http.Request) {

}

func FetchOrderHandler(writer http.ResponseWriter, request *http.Request) {

}

func CreateOrderHandler(writer http.ResponseWriter, request *http.Request) {

	createOrderDto := common.OrderCreateDto{}

	//

	if !common.ReadReqPayload(writer, request, &createOrderDto) {
		return
	}

	// create order
}
