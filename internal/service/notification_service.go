package service

import (
	"errors"
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
)

func FetchUserNotifications(userId uint) []*common.NotificationListDto {
	db := database.NewDatabaseConnection()
	var notifications []*common.NotificationListDto
	db.Raw("SELECT id, user_id, notification_type, message, status FROM notifications WHERE user_id = ?", userId).Scan(&notifications)
	return notifications
}

func UpdateNotificationAsRead(nId uint) error {
	db := database.NewDatabaseConnection()
	var notification model.Notification

	db.Find(&notification, "id = ?", nId)

	if notification.ID == 0 {
		return errors.New("notification not found")
	}

	if notification.Status == model.READ_NOTIFICATION {
		return errors.New("already marked as read")
	}

	notification.Status = model.READ_NOTIFICATION
	db.Save(notification)
	return nil
}

func CreateNotification(notificationCreateDto common.NotificationCreateDto) (*model.Notification, error) {

	// create
	notification := model.Notification{
		UserID:           notificationCreateDto.UserId,
		NotificationType: notificationCreateDto.NotificationType,
		Message:          notificationCreateDto.Message,
		Status:           model.UNREAD_NOTIFICATION,
	}

	db := database.NewDatabaseConnection()
	return &notification, db.Create(&notification).Error
}
