package cart

import (
	"encoding/json"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/service"
	"log/slog"
	"net/http"
	"strconv"
)

func AddItemCartHandler(w http.ResponseWriter, r *http.Request) {
	dto := common.AddCartItemDto{}

	//fetch body
	if ok := common.ReadReqPayload(w, r, &dto); !ok {
		return
	}

	//check stock available
	if isStock, err := service.CheckStockAvailable(dto.Quantity, dto.ProductVariantID); err != nil {
		common.HandleErrorRes(w, map[string]string{"message": "product variant does not exist with given id"})
		return
	} else if !isStock {
		common.HandleErrorRes(w, map[string]string{"message": "item out-of-stock"})
		return
	}

	//validation check
	if errMap := validationCartItem(dto); len(errMap) > 0 {
		common.HandleErrorRes(w, errMap)
		return
	}

	//add item to cart
	cartItem, err := service.AddItemCart(dto)

	if err != nil {
		common.HandleErrorRes(w, map[string]string{"message": "failed to add item in cart"})
		return
	}

	common.SendOkRes(w, map[string]string{"id": strconv.Itoa(int(cartItem.ID)), "cart_id": strconv.Itoa(int(cartItem.CartID)),
		"variant_id": strconv.Itoa(int(cartItem.ProductVariantID)), "quantity": strconv.Itoa(int(dto.Quantity))})
}

func validationCartItem(dto common.AddCartItemDto) map[string]string {
	errors := make(map[string]string)

	//check userId exists
	if !service.CheckUserExists(dto.UserID) {
		errors["userId"] = "user does not exist with given id"
	}

	// Check productVariantId exists
	if !service.CheckVariantExists(dto.ProductVariantID) {
		errors["variantId"] = "product variant does not exist with given id"
	}

	//check quantity
	if dto.Quantity <= 0 {
		errors["quantity"] = "quantity cannot be zero or negative"
	}

	return errors
}

func FetchCartDetailsHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("user_id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Debug("cannot read user_id from " + r.PathValue("user_id"))
		return
	}

	// check user exists with userId
	if !service.CheckUserExists(uint(userId)) {
		common.HandleErrorRes(w, map[string]string{"message": "user does not exists with given id"})
		return
	}

	// fetch user's cart items
	cartItem, err := service.FetchCartItem(uint(userId))
	if err != nil {
		common.HandleErrorRes(w, map[string]string{"message": err.Error()})
		return
	}

	err = json.NewEncoder(w).Encode(cartItem)

	if err != nil {
		slog.Error(err.Error())
		common.HandleErrorRes(w, map[string]string{"message": "failed to fetch user's cart item"})
	}

}

func RemoveCartItemsHandler(w http.ResponseWriter, r *http.Request) {
	cartItemId, err := strconv.Atoi(r.PathValue("cart_item_id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Debug("cannot read cart item id from " + r.PathValue("cart_item_id"))
		return
	}

	//delete cartItem by cartItemId
	err = service.RemoveCartItem(uint(cartItemId))

	if err != nil {
		common.HandleErrorRes(w, map[string]string{"message": err.Error()})
		slog.Error("Failed to delete cart item:", err)
		return
	}

	common.SendOkRes(w, map[string]string{"message": "cart item deleted successfully"})
}
