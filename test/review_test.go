package test

import (
	"encoding/json"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	testutil "github.com/Digital-AIR/bizio-ecommerce/test/util"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestCreateValidReview(t *testing.T) {
	db := database.NewDatabaseConnection()

	category := model.Category{Name: "testcategory", Description: "test category description"}
	db.Create(&category)

	brand := model.Brand{Name: "testbrand", Description: "test brand description"}
	db.Create(&brand)

	user := model.User{Username: testutil.RandomString(5), Email: testutil.RandomString(5) + "@example.com", PasswordHash: "password"}
	user1 := model.User{Username: testutil.RandomString(5), Email: testutil.RandomString(5) + "@example1.com", PasswordHash: "password1"}
	db.Create(&user)
	db.Create(&user1)

	product := model.Product{Name: "testproduct", Description: "test description", CategoryID: category.ID, BrandID: brand.ID}
	db.Create(&product)

	reviewCreateDto := common.ReviewCreateDto{
		UserID:  user.ID,
		Rating:  4,
		Comment: "Good product",
	}

	payload, err := json.Marshal(reviewCreateDto)

	//valid product review
	resPayload, resp, err := testutil.MakeReq("POST", "/products/"+strconv.Itoa(int(product.ID))+"/reviews", payload, nil)
	if err != nil {
		slog.Error(err.Error())
		t.Fail()
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	slog.Info(string(resPayload))

	//non-existent productId
	resPayload, resp, err = testutil.MakeReq("POST", "/products/"+strconv.Itoa(100)+"/reviews", payload, nil)
	if resp.StatusCode != 400 {
		slog.Error(err.Error())
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, strings.ToLower(string(resPayload)), "productid not found")
	slog.Info(string(resPayload))

	//duplicate userId and productId
	resPayload, resp, err = testutil.MakeReq("POST", "/products/"+strconv.Itoa(int(product.ID))+"/reviews", payload, nil)
	if resp.StatusCode != 400 {
		slog.Error(err.Error())
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, strings.ToLower(string(resPayload)), "user has already reviewed product")
	slog.Info(string(resPayload))

	reviewCreateDto = common.ReviewCreateDto{
		UserID:  user.ID,
		Rating:  8,
		Comment: "Good product",
	}

	payload1, err := json.Marshal(reviewCreateDto)

	//invalid rating value
	resPayload, resp, err = testutil.MakeReq("POST", "/products/"+strconv.Itoa(int(product.ID))+"/reviews", payload1, nil)
	if resp.StatusCode != 400 {
		slog.Error(err.Error())
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, strings.ToLower(string(resPayload)), "rating must be between 1 and 5")
	slog.Info(string(resPayload))

	// Review without comment
	reviewCreateDto = common.ReviewCreateDto{
		UserID:  user1.ID,
		Rating:  3,
		Comment: "",
	}
	payload2, err := json.Marshal(reviewCreateDto)

	resPayload, resp, err = testutil.MakeReq("POST", "/products/"+strconv.Itoa(int(product.ID))+"/reviews", payload2, nil)
	if resp.StatusCode != 400 {
		slog.Error(err.Error())
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, string(resPayload), "comment is required")
	slog.Info(string(resPayload))

	//get product reviews api
	resPayload, resp, err = testutil.MakeReq("GET", "/products/"+strconv.Itoa(int(product.ID))+"/reviews", nil, nil)
	if resp.StatusCode != 200 {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	slog.Info(string(resPayload))

	// get reviews non-existent productId
	resPayload, resp, err = testutil.MakeReq("GET", "/products/100/reviews", nil, nil)
	if resp.StatusCode != 400 {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, string(resPayload), "productId not found")

	//get invalid productId format
	resPayload, resp, err = testutil.MakeReq("GET", "/products/test/reviews", nil, nil)
	if resp.StatusCode != 400 {
		t.Fail()
		return
	}
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, string(resPayload), "")

}
