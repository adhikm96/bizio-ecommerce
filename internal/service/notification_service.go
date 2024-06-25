package service

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
)

func FetchUserNotifications(userId uint) []*common.NotificationListDto {
	db := database.NewDatabaseConnection()
	var notifications []*common.NotificationListDto
	db.Raw("SELECT id, user_id, notification_type, message, status FROM notifications WHERE user_id = ?", userId).Scan(&notifications)
	return notifications
}

func FetchNotification(nId uint) *common.NotificationListDto {
	db := database.NewDatabaseConnection()
	var notification common.NotificationListDto
	db.Raw("SELECT id, user_id, notification_type, message, status FROM notifications WHERE id = ?", nId).Scan(&notification)
	return &notification
}
