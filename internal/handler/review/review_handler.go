package review

import (
	"encoding/json"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/service"
	"log/slog"
	"net/http"
	"strconv"
)

func CreateReviewHanlder(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.Atoi(r.PathValue("product_id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Debug("invalid product id " + r.PathValue("product_id"))
		return
	}

	//check productId exists
	if !service.CheckProductExists(uint(productId)) {
		common.HandleErrorRes(w, map[string]string{"message": "productId not found"})
		return
	}

	reviewCreateDto := common.ReviewCreateDto{}

	if ok := common.ReadReqPayload(w, r, &reviewCreateDto); !ok {
		return
	}

	//check userId exists
	if !service.CheckUserExists(reviewCreateDto.UserID) {
		common.HandleErrorRes(w, map[string]string{"message": "userId not found"})
		return
	}

	if reviewCreateDto.Rating < 1 || reviewCreateDto.Rating > 5 {
		common.HandleErrorRes(w, map[string]string{"message": "rating must be between 1 and 5"})
		return
	}

	if service.CheckSameUserIdAndProductIdExists(reviewCreateDto.UserID, uint(productId)) {
		common.HandleErrorRes(w, map[string]string{"message": "user has already reviewed product"})
		return
	}

	review, err := service.CreateReview(reviewCreateDto, uint(productId))

	if err != nil {
		common.HandleErrorRes(w, map[string]string{"message": "failed to create review"})
		return
	}

	common.SendOkRes(w, map[string]string{"id": strconv.Itoa(int(review.ID)), "rating": strconv.Itoa(int(review.Rating)), "comment": review.Comment})
}

func FetchReviewHandler(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.Atoi(r.PathValue("product_id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Debug("invalid productId " + r.PathValue("product_id"))
		return
	}

	//check productId exists
	if !service.CheckProductExists(uint(productId)) {
		common.HandleErrorRes(w, map[string]string{"message": "productId not found"})
		return
	}

	//get product review
	review := service.GetProductReview(uint(productId))

	err = json.NewEncoder(w).Encode(review)

	if err != nil {
		common.HandleErrorRes(w, map[string]string{"message": "failed to get product review"})
		return
	}
}
