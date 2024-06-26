package review

import (
	"encoding/json"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	"github.com/Digital-AIR/bizio-ecommerce/internal/server"
	testutil "github.com/Digital-AIR/bizio-ecommerce/internal/test/util"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCreateValidReview(t *testing.T) {

	startServer()
	db := database.NewDatabaseConnection()

	category := model.Category{Name: "testcategory", Description: "test category description"}
	db.Create(&category)

	brand := model.Brand{Name: "testbrand", Description: "test brand description"}
	db.Create(&brand)

	user := model.User{Username: testutil.RandomString(5), Email: testutil.RandomString(5) + "@example.com", PasswordHash: "password"}
	db.Create(&user)

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

func startServer() {
	go server.InitServer()
	time.Sleep(time.Second * 1)
}
