package service

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
)

func CheckSameUserIdAndProductIdExists(userId uint, productId uint) bool {
	db := database.GetDbConn()
	var exists bool

	row := db.Raw("SELECT EXISTS(SELECT 1 FROM reviews WHERE user_id =? AND product_id =?)", userId, productId).Row()
	row.Scan(&exists)

	return exists
}

func CreateReview(reviewCreateDto common.ReviewCreateDto, productId uint) (*model.Review, error) {

	review := model.Review{
		UserID:    reviewCreateDto.UserID,
		Rating:    reviewCreateDto.Rating,
		Comment:   reviewCreateDto.Comment,
		ProductID: productId,
	}

	db := database.GetDbConn()
	return &review, db.Create(&review).Error
}

func GetProductReview(productId uint) []*common.ReviewListDto {
	var reviews []*common.ReviewListDto
	db := database.GetDbConn()
	db.Table("reviews").Select("id, user_id, product_id, rating, comment").Where("product_id = ?", productId).Scan(&reviews)
	return reviews
}
