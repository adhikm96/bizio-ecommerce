package notification

import (
	"encoding/json"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	"github.com/Digital-AIR/bizio-ecommerce/internal/service"
	"log/slog"
	"net/http"
	"strconv"
)

func GetHandler(writer http.ResponseWriter, request *http.Request) {
	nID, err := strconv.Atoi(request.PathValue("id"))

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		slog.Debug("cannot read id from " + request.PathValue("id"))
		return
	}

	notification := service.Fetch(uint(nID))

	if notification.ID == 0 {
		common.HandleErrorRes(writer, map[string]string{"message": "notification does not exists with given id"})
		return
	}

	err = json.NewEncoder(writer).Encode(notification)

	if err != nil {
		slog.Error(err.Error())
		common.HandleErrorRes(writer, map[string]string{"message": "failed to fetch user's notifications"})
	}
}

func UsersNotificationHandler(writer http.ResponseWriter, request *http.Request) {
	userId, err := strconv.Atoi(request.PathValue("user_id"))

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		slog.Debug("cannot read user_id from " + request.PathValue("user_id"))
		return
	}

	// check user exists with userId

	if !service.CheckUserExists(uint(userId)) {
		common.HandleErrorRes(writer, map[string]string{"message": "user does not exists with given id"})
		return
	}

	// fetch user's notifications
	notifications := service.FetchUserNotifications(uint(userId))

	err = json.NewEncoder(writer).Encode(notifications)

	if err != nil {
		slog.Error(err.Error())
		common.HandleErrorRes(writer, map[string]string{"message": "failed to fetch user's notifications"})
	}
}

func CreateHandler(writer http.ResponseWriter, request *http.Request) {

	notificationCreateDto := common.NotificationCreateDto{}
	userId, err := strconv.Atoi(request.PathValue("user_id"))

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		slog.Debug("cannot read user_id from " + request.PathValue("user_id"))
		return
	}

	// fetch body
	if ok := common.ReadReqPayload(writer, request, &notificationCreateDto); !ok {
		return
	}

	// validate
	if errMap := validateNotificationCreate(notificationCreateDto); len(errMap) > 0 {
		common.HandleErrorRes(writer, errMap)
		return
	}

	// create
	notification := model.Notification{
		UserID:           uint(userId),
		NotificationType: notificationCreateDto.NotificationType,
		Message:          notificationCreateDto.Message,
		Status:           model.UNREAD_NOTIFICATION,
	}

	db := database.NewDatabaseConnection()

	res := db.Create(&notification)

	if res.Error != nil {
		slog.Error(err.Error())
		common.HandleErrorRes(writer, map[string]string{"message": "failed to create notification"})
		return
	}

	// respond
	common.HandleRes(writer, map[string]string{"notification_id": strconv.Itoa(int(notification.ID))}, http.StatusOK)
}

func validateNotificationCreate(dto common.NotificationCreateDto) map[string]string {

	errMap := make(map[string]string)

	// check type
	if dto.NotificationType == "" {
		errMap["notification_type"] = "notification_type is required"
	}

	return errMap
}
