package notification

import (
	"encoding/json"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	"github.com/Digital-AIR/bizio-ecommerce/internal/server"
	testutil "github.com/Digital-AIR/bizio-ecommerce/test/util"
	"log/slog"
	"strconv"
	"testing"
	"time"
)

func TestNotificationApi(t *testing.T) {

	startServer()

	user1 := model.User{
		Username:     testutil.RandomString(5),
		Email:        testutil.RandomString(5) + "@example.com",
		PasswordHash: "password",
	}

	// create notification test

	db := database.NewDatabaseConnection()

	res := db.Create(&user1)

	if res.Error != nil {
		t.Fail()
		return
	}

	notiCreateDto := common.NotificationCreateDto{
		NotificationType: "EMAIL",
		Message:          "Sample Email",
	}

	payload, err := json.Marshal(notiCreateDto)

	_, response, _ := testutil.MakeReq("POST", "/notifications/"+strconv.Itoa(int(user1.ID)), []byte(""), nil)

	if response.StatusCode != 400 {
		t.Fail()
		return
	}

	resPayload, _, err := testutil.MakeReq("POST", "/notifications/"+strconv.Itoa(int(user1.ID)), payload, nil)

	if err != nil {
		slog.Error(err.Error())
		t.Fail()
		return
	}

	slog.Info(string(resPayload))

	resBody := make(map[string]string)

	err = json.Unmarshal(resPayload, &resBody)
	if err != nil {
		slog.Error(err.Error())
		t.Fail()
		return
	}

	nId, _ := strconv.Atoi(resBody["notification_id"])

	// get notification api

	_, response, _ = testutil.MakeReq("GET", "/notifications/0/read", []byte(""), nil)

	if response.StatusCode != 400 {
		t.Fail()
		return
	}

	resPayload, _, err = testutil.MakeReq("GET", "/notifications/"+strconv.Itoa(nId)+"/read", []byte(""), nil)

	if err != nil {
		slog.Error(err.Error())
		t.Fail()
		return
	}
	slog.Info(string(resPayload))

	resBody2 := common.NotificationListDto{}

	err = json.Unmarshal(resPayload, &resBody2)
	if err != nil {
		slog.Error(err.Error())
		t.Fail()
		return
	}

	if resBody2.NotificationType != notiCreateDto.NotificationType {
		t.Fail()
		return
	}

	if resBody2.Message != notiCreateDto.Message {
		t.Fail()
		return
	}

	if resBody2.UserID != user1.ID {
		t.Fail()
		return
	}

	// get user's notification api
	_, response, _ = testutil.MakeReq("GET", "/notifications/0", []byte(""), nil)

	if response.StatusCode != 400 {
		t.Fail()
		return
	}

	resPayload, _, err = testutil.MakeReq("GET", "/notifications/"+strconv.Itoa(int(user1.ID)), []byte(""), nil)

	if err != nil {
		slog.Error(err.Error())
		t.Fail()
		return
	}

	var resBody3 []common.NotificationListDto

	slog.Info(string(resPayload))

	err = json.Unmarshal(resPayload, &resBody3)
	if err != nil {
		slog.Error(err.Error())
		t.Fail()
		return
	}

	if len(resBody3) != 1 {
		t.Fail()
		return
	}

	if resBody3[0].NotificationType != notiCreateDto.NotificationType {
		t.Fail()
		return
	}

	if resBody3[0].Message != notiCreateDto.Message {
		t.Fail()
		return
	}

	if resBody3[0].UserID != user1.ID {
		t.Fail()
		return
	}
}

func startServer() {
	go server.InitServer()
	time.Sleep(time.Second * 1)
}
