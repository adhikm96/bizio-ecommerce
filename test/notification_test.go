package test

import (
	"encoding/json"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	testutil "github.com/Digital-AIR/bizio-ecommerce/test/util"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"strconv"
	"testing"
)

func TestNotificationsHandler(t *testing.T) {
	user1 := model.User{
		Username:     testutil.RandomString(5),
		Email:        testutil.RandomString(5) + "@example.com",
		PasswordHash: "password",
	}

	// create notification test

	db := database.GetDbConn()

	res := db.Create(&user1)

	assert.Nil(t, res.Error)

	notiCreateDto := common.NotificationCreateDto{
		NotificationType: "EMAIL",
		Message:          "Sample Email",
		UserId:           user1.ID,
	}

	payload, err := json.Marshal(notiCreateDto)

	_, response, _ := testutil.MakeReq("POST", "/notifications", []byte(""), nil)

	assert.Equal(t, response.StatusCode, 400)

	resPayload, _, err := testutil.MakeReq("POST", "/notifications", payload, nil)

	assert.Nil(t, err)

	slog.Info(string(resPayload))

	resBody := make(map[string]string)

	err = json.Unmarshal(resPayload, &resBody)
	assert.Nil(t, err)

	nId, _ := strconv.Atoi(resBody["notification_id"])

	// check status of notification
	notification := model.Notification{}
	db.Find(&notification, "id = ?", nId)

	assert.Equal(t, notification.Status, model.UNREAD_NOTIFICATION)

	// get notification api
	_, response, _ = testutil.MakeReq("PUT", "/notifications/0/read", []byte(""), nil)

	assert.Equal(t, response.StatusCode, 400)

	_, _, err = testutil.MakeReq("PUT", "/notifications/"+strconv.Itoa(nId)+"/read", []byte(""), nil)

	assert.Nil(t, err)

	// check status of notification
	notification = model.Notification{}
	db.Find(&notification, "id = ?", nId)

	assert.Equal(t, notification.Status, model.READ_NOTIFICATION)

	// get user's notification api
	_, response, _ = testutil.MakeReq("GET", "/notifications/0", []byte(""), nil)

	assert.Equal(t, response.StatusCode, 400)

	resPayload, _, err = testutil.MakeReq("GET", "/notifications/"+strconv.Itoa(int(user1.ID)), []byte(""), nil)

	assert.Nil(t, err)

	var resBody3 []common.NotificationListDto

	slog.Info(string(resPayload))

	err = json.Unmarshal(resPayload, &resBody3)
	assert.Nil(t, err)

	assert.Equal(t, len(resBody3), 1)
	assert.Equal(t, resBody3[0].NotificationType, notification.NotificationType)
	assert.Equal(t, resBody3[0].Message, notification.Message)
	assert.Equal(t, resBody3[0].UserID, notification.UserID)
	assert.Equal(t, resBody3[0].Status, notification.Status)

}
